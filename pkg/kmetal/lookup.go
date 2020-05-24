package kmetal

import (
	"fmt"

	"github.com/gardener/controller-manager-library/pkg/logger"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-go/api/models"
)

const Mega = 1024 * 1024

func Lookup(logger logger.LogContext, driver *metalgo.Driver, uuid string, macs []string) (*models.V1MachineResponse, error) {
	mfr := &metalgo.MachineFindRequest{}
	if uuid != "" {
		mfr.ID = &uuid
		resp, err := driver.MachineFind(mfr)
		if err != nil {
			logger.Infof("lookup uuid %s failed: %s", uuid, err)
			return nil, err
		}
		if len(resp.Machines) > 0 {
			logger.Infof("lookup uuid %s found", uuid)
			return resp.Machines[0], nil
		}
		logger.Infof("lookup uuid %s not found", uuid)
	}

	for _, mac := range macs {
		mfr.ID = nil
		mfr.NicsMacAddresses = []string{mac}
		resp, err := driver.MachineFind(mfr)
		if err != nil {
			return nil, err
		}
		if len(resp.Machines) > 0 {
			if uuid != "" {
				mismatch := ""
				for _, m := range resp.Machines {
					if m.ID != nil {
						if *m.ID == "uuid" {
							return m, nil
						}
						mismatch = *m.ID
					}
				}
				if mismatch != "" {
					return nil, fmt.Errorf("uuid mismatch for mac %s: %s != %s", mac, uuid, mismatch)
				}
			}
			logger.Infof("lookup mac %s found uuid %s", mac, p(resp.Machines[0].ID))
			return resp.Machines[0], nil
		}
		logger.Infof("lookup mac %s not found", mac)
	}
	return nil, nil
}

func p(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func FillMetadata(m *models.V1MachineResponse, metadata map[string]interface{}) {
	attributes := map[string]interface{}{}
	if m.Hardware != nil {
		macs := map[string]interface{}{}
		list := []string{}
		for _, n := range m.Hardware.Nics {
			if n.Mac != nil {
				if p(n.Name) == "lo" {
					if p(n.Mac) != "" {
						macs["lo"] = []string{p(n.Mac)}
					}
				} else {
					if p(n.Mac) != "" {
						list = append(list, p(n.Mac))
					}
				}
			}
		}
		macs["regular"] = list
		metadata["macs"] = macs

		if m.Hardware.CPUCores != nil {
			attributes["cores"] = *m.Hardware.CPUCores
		}
		if m.Hardware.Memory != nil {
			attributes["memory"] = *m.Hardware.Memory / Mega
		}
		disks := []interface{}{}
		for _, d := range m.Hardware.Disks {
			disk := map[string]interface{}{}
			if d.Name != nil {
				disk["name"] = *d.Name
			}
			if d.Size != nil {
				disk["size"] = *d.Size / Mega
			}
			if d.Primary != nil {
				disk["primary"] = *d.Primary
			} else {
				disk["primary"] = false
			}
			disks = append(disks, disk)
		}
		attributes["disks"] = disks
	}
	if m.State != nil {
		attributes["state"] = *m.State.Value
	}
	if m.Allocation != nil {
		attributes["Description"] = m.Allocation.Description
		if m.Allocation.Image != nil {
			if m.Allocation.Image.URL != "" {
				attributes["imageURL"] = m.Allocation.Image.URL
			}
			if m.Allocation.Image.Name != "" {
				attributes["imageName"] = m.Allocation.Image.Name
			}
			if m.Allocation.Project != nil {
				attributes["project"] = *m.Allocation.Project
			}
		}
	}
	if m.Bios != nil {
		attributes["vendor"] = p(m.Bios.Vendor)
		attributes["bios-version"] = p(m.Bios.Version)
	}
	metadata["attributes"] = attributes
}
