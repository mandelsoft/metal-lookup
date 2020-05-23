package main

import (
	"context"
	"fmt"

	"github.com/gardener/controller-manager-library/pkg/config"
	"github.com/spf13/cobra"

	"github.com/mandelsoft/kmetal/pkg/command"
	"github.com/mandelsoft/kmetal/pkg/kmetal/client"

	metalgo "github.com/metal-stack/metal-go"
)

var Version = "dev-version"

type Config struct {
	MetalConfig string
	client.DriverConfig

	Mac  string
	UUID string
}

func (this *Config) AddOptionsToSet(set config.OptionSet) {
	set.AddStringOption(&this.MetalConfig, "metalconfig", "", "", "config file for metal-api")
	this.DriverConfig.AddOptionsToSet(set)

	set.AddStringOption(&this.Mac, "mac", "", "", "mac address")
	set.AddStringOption(&this.UUID, "uuid", "", "", "UUID")
}

func (this *Config) Evaluate() error {
	if this.Mac == "" && this.UUID == "" {
		return fmt.Errorf("mac or uuid must be given§")
	}
	return this.DriverConfig.Evaluate()
}

////////////////////////////////////////////////////////////////////////////////

const Mega = 1024 * 1024

func main() {
	command.Start(Lookup)
}

func Lookup(ctx context.Context) *cobra.Command {
	return command.NewCommand(ctx, Version, "<options>", "machine lookup", "lookup machine objects", doit, &Config{})
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
	mfr := &metalgo.MachineFindRequest{}
	if cfg.Mac != "" {
		mfr.NicsMacAddresses = []string{cfg.Mac}
	}
	if cfg.UUID != "" {
		mfr.ID = &cfg.UUID
	}
	resp, err := driver.MachineFind(mfr)
	if err != nil {
		return err
	}
	if len(resp.Machines) == 0 {
		return fmt.Errorf("no machine found\n")
	}
	for _, m := range resp.Machines {
		fmt.Printf("UUID:        %s\n", p(m.ID))
		fmt.Printf("Liveliness:  %s\n", p(m.Liveliness))
		a := []string{}
		if m.Hardware != nil {
			for _, n := range m.Hardware.Nics {
				if n.Mac != nil {
					a = append(a, p(n.Mac))
				}
			}
			fmt.Printf("Macs:        %s\n", a)
			if m.Hardware.CPUCores != nil {
				fmt.Printf("Cores:       %d\n", *m.Hardware.CPUCores)
			}
			if m.Hardware.Memory != nil {
				fmt.Printf("Memory:      %dM\n", *m.Hardware.Memory/Mega)
			}
			for _, d := range m.Hardware.Disks {
				fmt.Printf("  Disk:      %s\n", p(d.Name))
				if d.Size != nil {
					fmt.Printf("  Size:      %dM\n", *d.Size/Mega)
				}
				if d.Primary != nil {
					fmt.Printf("  Primary:   %t\n", *d.Primary)
				}
			}
		}
		if m.State != nil {
			fmt.Printf("State:     %s\n", p(m.State.Value))
		}
		if m.Allocation != nil {
			fmt.Printf("Desc:      %s\n", m.Allocation.Description)
		}
		if m.Bios != nil {
			fmt.Printf("Vendor:      %s\n", p(m.Bios.Vendor))
			fmt.Printf("Version:     %s\n", p(m.Bios.Version))
		}
	}
	return nil
}

func p(s *string) string {
	if s == nil {
		return "<none>"
	}
	return *s
}
