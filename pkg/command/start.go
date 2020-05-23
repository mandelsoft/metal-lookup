/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package command

import (
	"context"
	"github.com/gardener/controller-manager-library/pkg/ctxutil"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const DeletionActivity = "DeletionActivity"

var GracePeriod time.Duration

type CommandFactory func(ctx context.Context) *cobra.Command

func Start(fac CommandFactory) {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	var (
		cctx = ctxutil.CancelContext(ctxutil.WaitGroupContext(context.Background(), "main"))
		ctx  = ctxutil.TickContext(cctx, DeletionActivity)
		c    = make(chan os.Signal, 2)
		t    = make(chan os.Signal, 2)
	)

	signal.Notify(t, syscall.SIGTERM, syscall.SIGQUIT)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		cnt := 0
	loop:
		for {
			select {
			case <-c:
				cnt++
				if cnt == 2 {
					break loop
				}
				logger.Infof("process is being terminated")
				ctxutil.Cancel(ctx)
			case <-t:
				cnt++
				if cnt == 2 {
					break loop
				}
				grace := GracePeriod
				if grace > 0 {
					logger.Infof("process is being terminated with grace period for cleanup")
					go ctxutil.CancelAfterInactivity(ctx, DeletionActivity, grace)
				} else {
					logger.Infof("process is being terminated without grace period")
					ctxutil.Cancel(ctx)
				}
			}
		}
		logger.Infof("process is aborted immediately")
		os.Exit(0)
	}()

	//	if err := plugins.HandleCommandLine("--plugin-file", os.Args); err != nil {
	//		panic(err)
	//	}

	cmd := fac(ctx)
	cmd.Flags().DurationVarP(&GracePeriod, "grace-period", "", 120*time.Second, "Grace period for shutdown")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	var gracePeriod = 120 * time.Second
	logger.Infof("waiting for everything to shutdown (max. %d seconds)", gracePeriod/time.Second)
	ctxutil.WaitGroupWait(ctx, gracePeriod, "main")
	logger.Infof("program exists")
}
