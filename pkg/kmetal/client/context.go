package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// MetalConfig contains all configuration contexts of metalctl
type MetalConfig struct {
	source          string
	CurrentContext  string `yaml:"current"`
	PreviousContext string `yaml:"previous"`
	Contexts        map[string]Context
}

// Context configure metalctl behaviour
type Context struct {
	ApiURL       string `yaml:"url"`
	IssuerURL    string `yaml:"issuer_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	HMAC         string `yaml:"hmac"`
}

func ReadConfig(path string) (*MetalConfig, error) {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(os.Getenv("HOME"), path[2:])
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read metal config %q: %s", path, err)
	}
	cfg := &MetalConfig{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal metal config %q: %s", path, err)
	}
	cfg.source = path
	return cfg, err

}

func (this *MetalConfig) GetCurrentConfig() (*DriverConfig, error) {
	if this.CurrentContext == "" {
		return nil, fmt.Errorf("no current context in metal config %q", this.source)
	}
	if cur, ok := this.Contexts[this.CurrentContext]; ok {
		return &DriverConfig{
			DriverURL: cur.ApiURL,
			HMAC:      cur.HMAC,
			Token:     "",
		}, nil
	}
	return nil, fmt.Errorf("current context %q not found in metal config %q", this.CurrentContext, this.source)
}
