package kmetal

import (
	"fmt"
	"reflect"

	"github.com/gardener/controller-manager-library/pkg/logger"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-go/api/models"
)

const Mega = 1024 * 1024

const NAME = "name"
const ID = "id"

const ORIGIN = "ORIGIN"
const FORWARDED_IN = "__X-Forwarded-For__"

const UUID = "uuid"
const MACS_IN = "__mac__"
const PARTITION_IN = "partition"
const PARTITION = "partition_info"

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

const CMDLINE = "commandLine"
const INITRD = "initrd"
const KERNEL = "kernel"

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
			logger.Infof("lookup mac %s found uuid %s", mac, s(resp.Machines[0].ID))
			return resp.Machines[0], nil
		}
		logger.Infof("lookup mac %s not found", mac)
	}
	return nil, nil
}

func FillMetadataByFields(logger logger.LogContext, m *models.V1MachineResponse, fields Fields, metadata map[string]interface{}) {
	for n, f := range fields {
		value, err := f.source.Get(m)
		if err == nil && value != nil {
			t := f.source.Type()
			v := reflect.ValueOf(value)
			for t.Kind() == reflect.Ptr {
				v = v.Elem()
				t = t.Elem()
			}

			// TODO: extend fieldpath to handle maps
			err := f.target.Set(metadata, v.Interface())
			if err != nil {
				logger.Warnf("cannot set field %q: %s", n, err)
			}
			/*
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
			*/
		}
	}
}

func FillMetadata(logger logger.LogContext, m *models.V1MachineResponse, metadata map[string]interface{}) {
	metadata[UUID] = s(m.ID)
	attributes := map[string]interface{}{}
	if m.Hardware != nil {
		macs := map[string]interface{}{}
		list := []string{}
		for _, n := range m.Hardware.Nics {
			if n.Mac != nil {
				if s(n.Name) == SPECIAL_NIC {
					if s(n.Mac) != "" {
						macs[SPECIAL_NIC] = []string{s(n.Mac)}
					}
				} else {
					if s(n.Mac) != "" {
						list = append(list, s(n.Mac))
					}
				}
			}
		}
		macs[MACS_REGULAR] = list
		attributes[MACS] = macs

		set(CORES, m.Hardware.CPUCores, attributes)
		set(MEMORY, m.Hardware.Memory, attributes, MEM_SCALE)
		disks := []interface{}{}
		for _, d := range m.Hardware.Disks {
			disk := map[string]interface{}{}
			set(DSK_NAME, d.Name, disk)
			set(DSK_SIZE, d.Size, disk, DSK_SCALE)
			set(DSK_PRIMARY, d.Primary, disk)
			disks = append(disks, disk)
		}
		attributes[DISKS] = disks
	}
	if m.State != nil {
		set(STATE, m.State.Value, attributes)
	}
	if m.Allocation != nil {
		set(DESC, m.Allocation.Description, attributes)
		if m.Allocation.BootInfo != nil {
			set(CMDLINE, m.Allocation.BootInfo.Cmdline, attributes)
			set(INITRD, m.Allocation.BootInfo.Initrd, attributes)
			set(KERNEL, m.Allocation.BootInfo.Kernel, attributes)
		}
		set(PROJECT, m.Allocation.Project, attributes)
	}
	if m.Bios != nil {
		set(VENDOR, m.Bios.Vendor, attributes)
		set(VERSION, m.Bios.Version, attributes)
	}
	metadata[ATTRIBUTES] = attributes
}

func set(name string, value interface{}, dst map[string]interface{}, scales ...int) {
	scale := 1
	for _, s := range scales {
		scale = scale * s
	}
	switch v := value.(type) {
	case *int64:
		if v != nil {
			dst[name] = *v / int64(scale)
		}
	case *int32:
		if v != nil {
			dst[name] = *v / int32(scale)
		}
	case *int:
		if v != nil {
			dst[name] = *v / (scale)
		}
	case *bool:
		if v != nil {
			dst[name] = *v
		}
	case *string:
		if v != nil && *v != "" {
			dst[name] = *v
		}
	case string:
		if v != "" {
			dst[name] = v
		}
	}
}
