package command

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/controller-manager-library/pkg/config"
	"github.com/gardener/controller-manager-library/pkg/configmain"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/spf13/cobra"
)

type Main func(ctx context.Context, src config.OptionSource) error

func NewCommand(ctx context.Context, version, use, short, long string, main Main, cfgsrc config.OptionSource) *cobra.Command {
	ctx, cfg := configmain.WithConfig(ctx, nil)

	cfg.AddSource("lookup", cfgsrc)

	fileName := ""
	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Version: version,
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if fileName != "" {
			logger.Infof("reading config from file %q", fileName)
			if err := config.MergeConfigFile(fileName, cmd.Flags(), false); err != nil {
				return fmt.Errorf("invalid config file %q; %s", fileName, err)
			}
		}
		if err := cfg.Evaluate(); err != nil {
			return err
		}
		if err := main(ctx, cfgsrc); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		return nil
	}

	cfg.AddToCommand(cmd)
	cmd.Flags().StringVarP(&fileName, "config", "", "", "config file")
	return cmd
}
