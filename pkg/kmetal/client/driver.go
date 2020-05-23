package client

import (
	metalgo "github.com/metal-stack/metal-go"
)

func NewDriver(cfg *DriverConfig) (*metalgo.Driver, error) {
	return metalgo.NewDriver(cfg.DriverURL, cfg.Token, cfg.HMAC)
}
