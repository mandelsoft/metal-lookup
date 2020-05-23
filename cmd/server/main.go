package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/gardener/controller-manager-library/pkg/config"
	"github.com/gardener/controller-manager-library/pkg/ctxutil"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/gardener/controller-manager-library/pkg/server"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/spf13/cobra"

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
	cfg := src.(*Config)

	access, err := client.GetDriverConfig(cfg.MetalConfig, &cfg.DriverConfig)

	if err != nil {
		return err
	}
	fmt.Printf("Hallo\n")
	fmt.Printf("Driver: %s\n", access.DriverURL)
	fmt.Printf("HMAC  : %s\n", access.HMAC)

	driver, err := client.NewDriver(access)
	if err != nil {
		return err
	}

	ctx = ctxutil.WaitGroupContext(ctx)
	s := server.NewHTTPServer(ctx, logger.New(), "metal-lookup")
	s.Register("/lookup", NewLookupHandler(ctx, driver).Lookup)
	s.Register("/healthz", Healthz)

	s.Start(nil, "localhost", cfg.Port)
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
	ct := r.Header.Get("Content-Type")
	if ct != "" && ct != restful.MIME_JSON {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("required %s\n", restful.MIME_JSON)))
		return
	}
	metadata := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&metadata)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	macs := []string{}
	e := metadata["__mac__"]
	if e != nil {
		if m, ok := e.([]interface{}); ok {
			for _, v := range m {
				macs = append(macs, v.(string))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("invalid mac list %T\n", e)))
			return
		}
	}
	uuid := ""
	e = metadata["uuid"]
	if e != nil {
		if m, ok := e.(string); ok {
			uuid = m
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid uuid\n"))
			return
		}
	}
	logger.Infof("lookup uuud: %s, macs: %s", uuid, macs)
	machine, err := kmetal.Lookup(this.logger, this.driver, uuid, macs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		this.logger.Infof("macs: %s", e)
		w.Write([]byte(fmt.Sprintf("cannot lookup: %s\n", err)))
		return
	}

	if machine != nil {
		kmetal.FillMetadata(machine, metadata)
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error() + "\n"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", restful.MIME_JSON)
	w.Write(data)
}
