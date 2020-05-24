package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/gardener/controller-manager-library/pkg/config"
	"github.com/gardener/controller-manager-library/pkg/configmain"
	"github.com/gardener/controller-manager-library/pkg/ctxutil"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/gardener/controller-manager-library/pkg/server"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings"

	"github.com/mandelsoft/kmetal/pkg/command"
	"github.com/mandelsoft/kmetal/pkg/kmetal"
	"github.com/mandelsoft/kmetal/pkg/kmetal/client"
)

var Version = "dev-version"

type Config struct {
	MetalConfig string
	client.DriverConfig
	Port int
}

func (this *Config) AddOptionsToSet(set config.OptionSet) {
	set.AddIntOption(&this.Port, "port", "", 8080, "server port")
	set.AddStringOption(&this.MetalConfig, "metalconfig", "", "", "config file for metal-api")
	this.DriverConfig.AddOptionsToSet(set)
}

func (this *Config) Evaluate() error {
	return this.DriverConfig.Evaluate()
}

////////////////////////////////////////////////////////////////////////////////

const Mega = 1024 * 1024

func main() {
	command.Start(Server)
}

func Server(ctx context.Context) *cobra.Command {
	return command.NewCommand(ctx, Version, "<options>", "machine lookup server", "lookup machine objects", doit, &Config{})
}

func doit(ctx context.Context, src config.OptionSource) error {
	logger := logger.New()
	cfg := src.(*Config)

	scfg := configmain.Get(ctx).GetSource(server.OPTION_SOURCE).(*server.Config)
	if scfg.ServerPortHTTP == 0 {
		scfg.ServerPortHTTP = 8080
	}
	access, err := client.GetDriverConfig(cfg.MetalConfig, &cfg.DriverConfig)

	if err != nil {
		return err
	}
	logger.Infof("Driver: %s", access.DriverURL)
	logger.Infof("HMAC  : %s...", strings.ShortenString(access.HMAC, 3))

	driver, err := client.NewDriver(access)
	if err != nil {
		return err
	}

	ctx = ctxutil.WaitGroupContext(ctx)
	s := server.NewHTTPServer(ctx, logger, "metal-lookup")
	s.Register("/lookup", NewLookupHandler(ctx, driver).Lookup)
	s.Register("/healthz", Healthz)

	s.Start(nil, scfg.BindAddress, scfg.ServerPortHTTP)
	ctxutil.WaitGroupWait(ctx, 0)
	return nil
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type LookupHandler struct {
	ctx    context.Context
	logger logger.LogContext
	driver *metalgo.Driver
}

func NewLookupHandler(ctx context.Context, driver *metalgo.Driver) *LookupHandler {
	return &LookupHandler{
		ctx:    ctx,
		logger: logger.New(),
		driver: driver,
	}
}

func (this *LookupHandler) Lookup(w http.ResponseWriter, r *http.Request) {
	err := this.lookup(w, r)
	if err != nil {
		this.logger.Infof("failed: %s", err)
	}
}

func fail(w http.ResponseWriter, status int, msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	w.WriteHeader(status)
	w.Write([]byte(err.Error() + "\n"))
	return err
}

func (this *LookupHandler) lookup(w http.ResponseWriter, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" && ct != restful.MIME_JSON {
		return fail(w, http.StatusUnsupportedMediaType, "required %s\n", restful.MIME_JSON)
	}
	metadata := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&metadata)

	if err != nil {
		return fail(w, http.StatusBadRequest, "%s", err)
	}

	macs := []string{}
	e := metadata[kmetal.MACS_IN]
	if e != nil {
		if m, ok := e.([]interface{}); ok {
			for _, v := range m {
				macs = append(macs, v.(string))
			}
		} else {
			return fail(w, http.StatusBadRequest, "invalid mac list %T", e)
		}
	}
	uuid := ""
	e = metadata[kmetal.UUID]
	if e != nil {
		if m, ok := e.(string); ok {
			uuid = m
		} else {
			return fail(w, http.StatusBadRequest, "invalid uuid")
		}
	}
	logger.Infof("lookup uuud: %s, macs: %s", uuid, macs)
	machine, err := kmetal.Lookup(this.logger, this.driver, uuid, macs)
	if err != nil {
		return fail(w, http.StatusBadRequest, "cannot lookup: %s", err)
	}

	if machine != nil {
		kmetal.FillMetadata(machine, metadata)
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fail(w, http.StatusInternalServerError, "%s", err)
	}
	w.Header().Set("Content-Type", restful.MIME_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return nil
}
