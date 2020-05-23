package client

import (
	"os"
	"path/filepath"

	"github.com/gardener/controller-manager-library/pkg/config"
)

type DriverConfig struct {
	DriverURL string
	HMAC      string
	Token     string
}

func (this *DriverConfig) AddOptionsToSet(set config.OptionSet) {
	set.AddStringOption(&this.DriverURL, "driver", "", "", "Driver URL")
	set.AddStringOption(&this.Token, "token", "", "", "Token")
	set.AddStringOption(&this.HMAC, "hmac", "", "", "HMAC")
}

func (this *DriverConfig) Evaluate() error {
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetDriverConfig(path string, d *DriverConfig) (*DriverConfig, error) {
	cfg := DriverConfig{}
	if path == "" {
		path = os.Getenv("METALCONFIG")
	}
	if path == "" {
		home := os.Getenv("HOME")
		if home != "" {
			p := filepath.Join(home, ".metalctl", "config.yaml")
			if fileExists(p) {
				path = p
			}
		}
	}

	if path != "" {
		c, err := ReadConfig(path)
		if err != nil {
			return nil, err
		}
		cur, err := c.GetCurrentConfig()
		if err != nil {
			return nil, err
		}
		cfg = *cur
	}

	if d != nil {
		if d.DriverURL != "" {
			cfg.DriverURL = d.DriverURL
		}
		if d.HMAC != "" {
			cfg.HMAC = d.HMAC
		}
		if d.Token != "" {
			cfg.Token = d.Token
		}
	}
	return &cfg, nil
}
