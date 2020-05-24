package kmetal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gardener/controller-manager-library/pkg/logger"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-go/api/models"
)

const Mega = 1024 * 1024

const UUID = "uuid"
const MACS_IN = "__mac__"

const ATTRIBUTES = "attributes"

const MACS = "macs"
const SPECIAL_NIC = "lo"
const MACS_REGULAR = "regular"

const DISKS = "disks"
const DSK_NAME = "name"
const DSK_SIZE = "size"
const DSK_PRIMARY = "primary"

const STATE = "state"
const CORES = "cores"
const MEMORY = "memory"
const DESC = "description"
const PROJECT = "project"
const VENDOR = "vendor"
const VERSION = "bios"
const IMAGE_NAME = "imageName"
const IMAGE_URL = "imageURL"

const MEM_SCALE = Mega
const DSK_SCALE = Mega

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
						if *m.ID == uuid {
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

func FillMetadataByFields(m *models.V1MachineResponse, fields Fields, metadata map[string]interface{}) {
	for n, f := range fields {
		value, err := f.Get(m)
		if err == nil && value != nil {
			t := f.Type()
			v := reflect.ValueOf(value)
			for t.Kind() == reflect.Ptr {
				v = v.Elem()
				t = t.Elem()
			}

			// TODO: extend fieldpath to handle maps
			path := strings.Split(n, ".")
			m := map[string]interface{}(metadata)
			for _, c := range path[:len(path)-1] {
				if c == "" {
					continue
				}
				next := m[c]
				if next == nil {
					next = map[string]interface{}{}
					m[c] = next
				}
				if f, ok := next.(map[string]interface{}); ok {
					m = f
				} else {
					return
				}
			}
			m[path[len(path)-1]] = v.Interface()
		}
	}
}

func FillMetadata(m *models.V1MachineResponse, metadata map[string]interface{}) {
	metadata[UUID] = p(m.ID)
	attributes := map[string]interface{}{}
	if m.Hardware != nil {
		macs := map[string]interface{}{}
		list := []string{}
		for _, n := range m.Hardware.Nics {
			if n.Mac != nil {
				if p(n.Name) == SPECIAL_NIC {
					if p(n.Mac) != "" {
						macs[SPECIAL_NIC] = []string{p(n.Mac)}
					}
				} else {
					if p(n.Mac) != "" {
						list = append(list, p(n.Mac))
					}
				}
			}
		}
		macs[MACS_REGULAR] = list
		attributes[MACS] = macs

		if m.Hardware.CPUCores != nil {
			attributes[CORES] = *m.Hardware.CPUCores
		}
		if m.Hardware.Memory != nil {
			attributes[MEMORY] = *m.Hardware.Memory / MEM_SCALE
		}
		disks := []interface{}{}
		for _, d := range m.Hardware.Disks {
			disk := map[string]interface{}{}
			if d.Name != nil {
				disk[DSK_NAME] = *d.Name
			}
			if d.Size != nil {
				disk[DSK_SIZE] = *d.Size / DSK_SCALE
			}
			if d.Primary != nil {
				disk[DSK_PRIMARY] = *d.Primary
			} else {
				disk[DSK_PRIMARY] = false
			}
			disks = append(disks, disk)
		}
		attributes[DISKS] = disks
	}
	if m.State != nil {
		attributes[STATE] = *m.State.Value
	}
	if m.Allocation != nil {
		attributes[DESC] = m.Allocation.Description
		if m.Allocation.Image != nil {
			if m.Allocation.Image.URL != "" {
				attributes[IMAGE_URL] = m.Allocation.Image.URL
			}
			if m.Allocation.Image.Name != "" {
				attributes[IMAGE_NAME] = m.Allocation.Image.Name
			}
			if m.Allocation.Project != nil {
				attributes[PROJECT] = *m.Allocation.Project
			}
		}
	}
	if m.Bios != nil {
		attributes[VENDOR] = p(m.Bios.Vendor)
		attributes[VERSION] = p(m.Bios.Version)
	}
	metadata[ATTRIBUTES] = attributes
}
