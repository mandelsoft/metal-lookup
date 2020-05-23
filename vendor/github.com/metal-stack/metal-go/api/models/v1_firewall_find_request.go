// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// V1FirewallFindRequest v1 firewall find request
// swagger:model v1.FirewallFindRequest
type V1FirewallFindRequest struct {

	// allocation hostname
	// Required: true
	AllocationHostname *string `json:"allocation_hostname"`

	// allocation image id
	// Required: true
	AllocationImageID *string `json:"allocation_image_id"`

	// allocation name
	// Required: true
	AllocationName *string `json:"allocation_name"`

	// allocation project
	// Required: true
	AllocationProject *string `json:"allocation_project"`

	// allocation succeeded
	// Required: true
	AllocationSucceeded *bool `json:"allocation_succeeded"`

	// disk names
	// Required: true
	DiskNames []string `json:"disk_names"`

	// disk sizes
	// Required: true
	DiskSizes []int64 `json:"disk_sizes"`

	// fru board mfg
	// Required: true
	FruBoardMfg *string `json:"fru_board_mfg"`

	// fru board mfg serial
	// Required: true
	FruBoardMfgSerial *string `json:"fru_board_mfg_serial"`

	// fru board part number
	// Required: true
	FruBoardPartNumber *string `json:"fru_board_part_number"`

	// fru chassis part number
	// Required: true
	FruChassisPartNumber *string `json:"fru_chassis_part_number"`

	// fru chassis part serial
	// Required: true
	FruChassisPartSerial *string `json:"fru_chassis_part_serial"`

	// fru product manufacturer
	// Required: true
	FruProductManufacturer *string `json:"fru_product_manufacturer"`

	// fru product part number
	// Required: true
	FruProductPartNumber *string `json:"fru_product_part_number"`

	// fru product serial
	// Required: true
	FruProductSerial *string `json:"fru_product_serial"`

	// hardware cpu cores
	// Required: true
	HardwareCPUCores *int64 `json:"hardware_cpu_cores"`

	// hardware memory
	// Required: true
	HardwareMemory *int64 `json:"hardware_memory"`

	// id
	// Required: true
	ID *string `json:"id"`

	// ipmi address
	// Required: true
	IPMIAddress *string `json:"ipmi_address"`

	// ipmi interface
	// Required: true
	IPMIInterface *string `json:"ipmi_interface"`

	// ipmi mac address
	// Required: true
	IPMIMacAddress *string `json:"ipmi_mac_address"`

	// ipmi user
	// Required: true
	IPMIUser *string `json:"ipmi_user"`

	// liveliness
	// Required: true
	Liveliness *string `json:"liveliness"`

	// name
	// Required: true
	Name *string `json:"name"`

	// network asns
	// Required: true
	NetworkAsns []int64 `json:"network_asns"`

	// network destination prefixes
	// Required: true
	NetworkDestinationPrefixes []string `json:"network_destination_prefixes"`

	// network ids
	// Required: true
	NetworkIds []string `json:"network_ids"`

	// network ips
	// Required: true
	NetworkIps []string `json:"network_ips"`

	// network nat
	// Required: true
	NetworkNat *bool `json:"network_nat"`

	// network prefixes
	// Required: true
	NetworkPrefixes []string `json:"network_prefixes"`

	// network private
	// Required: true
	NetworkPrivate *bool `json:"network_private"`

	// network underlay
	// Required: true
	NetworkUnderlay *bool `json:"network_underlay"`

	// network vrfs
	// Required: true
	NetworkVrfs []int64 `json:"network_vrfs"`

	// nics mac addresses
	// Required: true
	NicsMacAddresses []string `json:"nics_mac_addresses"`

	// nics names
	// Required: true
	NicsNames []string `json:"nics_names"`

	// nics neighbor mac addresses
	// Required: true
	NicsNeighborMacAddresses []string `json:"nics_neighbor_mac_addresses"`

	// nics neighbor names
	// Required: true
	NicsNeighborNames []string `json:"nics_neighbor_names"`

	// nics neighbor vrfs
	// Required: true
	NicsNeighborVrfs []string `json:"nics_neighbor_vrfs"`

	// nics vrfs
	// Required: true
	NicsVrfs []string `json:"nics_vrfs"`

	// partition id
	// Required: true
	PartitionID *string `json:"partition_id"`

	// rackid
	// Required: true
	Rackid *string `json:"rackid"`

	// sizeid
	// Required: true
	Sizeid *string `json:"sizeid"`

	// state value
	// Required: true
	StateValue *string `json:"state_value"`

	// tags
	// Required: true
	Tags []string `json:"tags"`
}

// Validate validates this v1 firewall find request
func (m *V1FirewallFindRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAllocationHostname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAllocationImageID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAllocationName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAllocationProject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAllocationSucceeded(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDiskNames(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDiskSizes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruBoardMfg(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruBoardMfgSerial(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruBoardPartNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruChassisPartNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruChassisPartSerial(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruProductManufacturer(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruProductPartNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFruProductSerial(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHardwareCPUCores(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHardwareMemory(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIPMIAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIPMIInterface(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIPMIMacAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIPMIUser(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLiveliness(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkAsns(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkDestinationPrefixes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkIds(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkIps(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkNat(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkPrefixes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkPrivate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkUnderlay(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkVrfs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsMacAddresses(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsNames(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsNeighborMacAddresses(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsNeighborNames(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsNeighborVrfs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNicsVrfs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePartitionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRackid(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSizeid(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStateValue(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1FirewallFindRequest) validateAllocationHostname(formats strfmt.Registry) error {

	if err := validate.Required("allocation_hostname", "body", m.AllocationHostname); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateAllocationImageID(formats strfmt.Registry) error {

	if err := validate.Required("allocation_image_id", "body", m.AllocationImageID); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateAllocationName(formats strfmt.Registry) error {

	if err := validate.Required("allocation_name", "body", m.AllocationName); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateAllocationProject(formats strfmt.Registry) error {

	if err := validate.Required("allocation_project", "body", m.AllocationProject); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateAllocationSucceeded(formats strfmt.Registry) error {

	if err := validate.Required("allocation_succeeded", "body", m.AllocationSucceeded); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateDiskNames(formats strfmt.Registry) error {

	if err := validate.Required("disk_names", "body", m.DiskNames); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateDiskSizes(formats strfmt.Registry) error {

	if err := validate.Required("disk_sizes", "body", m.DiskSizes); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruBoardMfg(formats strfmt.Registry) error {

	if err := validate.Required("fru_board_mfg", "body", m.FruBoardMfg); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruBoardMfgSerial(formats strfmt.Registry) error {

	if err := validate.Required("fru_board_mfg_serial", "body", m.FruBoardMfgSerial); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruBoardPartNumber(formats strfmt.Registry) error {

	if err := validate.Required("fru_board_part_number", "body", m.FruBoardPartNumber); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruChassisPartNumber(formats strfmt.Registry) error {

	if err := validate.Required("fru_chassis_part_number", "body", m.FruChassisPartNumber); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruChassisPartSerial(formats strfmt.Registry) error {

	if err := validate.Required("fru_chassis_part_serial", "body", m.FruChassisPartSerial); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruProductManufacturer(formats strfmt.Registry) error {

	if err := validate.Required("fru_product_manufacturer", "body", m.FruProductManufacturer); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruProductPartNumber(formats strfmt.Registry) error {

	if err := validate.Required("fru_product_part_number", "body", m.FruProductPartNumber); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateFruProductSerial(formats strfmt.Registry) error {

	if err := validate.Required("fru_product_serial", "body", m.FruProductSerial); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateHardwareCPUCores(formats strfmt.Registry) error {

	if err := validate.Required("hardware_cpu_cores", "body", m.HardwareCPUCores); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateHardwareMemory(formats strfmt.Registry) error {

	if err := validate.Required("hardware_memory", "body", m.HardwareMemory); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateIPMIAddress(formats strfmt.Registry) error {

	if err := validate.Required("ipmi_address", "body", m.IPMIAddress); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateIPMIInterface(formats strfmt.Registry) error {

	if err := validate.Required("ipmi_interface", "body", m.IPMIInterface); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateIPMIMacAddress(formats strfmt.Registry) error {

	if err := validate.Required("ipmi_mac_address", "body", m.IPMIMacAddress); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateIPMIUser(formats strfmt.Registry) error {

	if err := validate.Required("ipmi_user", "body", m.IPMIUser); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateLiveliness(formats strfmt.Registry) error {

	if err := validate.Required("liveliness", "body", m.Liveliness); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkAsns(formats strfmt.Registry) error {

	if err := validate.Required("network_asns", "body", m.NetworkAsns); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkDestinationPrefixes(formats strfmt.Registry) error {

	if err := validate.Required("network_destination_prefixes", "body", m.NetworkDestinationPrefixes); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkIds(formats strfmt.Registry) error {

	if err := validate.Required("network_ids", "body", m.NetworkIds); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkIps(formats strfmt.Registry) error {

	if err := validate.Required("network_ips", "body", m.NetworkIps); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkNat(formats strfmt.Registry) error {

	if err := validate.Required("network_nat", "body", m.NetworkNat); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkPrefixes(formats strfmt.Registry) error {

	if err := validate.Required("network_prefixes", "body", m.NetworkPrefixes); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkPrivate(formats strfmt.Registry) error {

	if err := validate.Required("network_private", "body", m.NetworkPrivate); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkUnderlay(formats strfmt.Registry) error {

	if err := validate.Required("network_underlay", "body", m.NetworkUnderlay); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNetworkVrfs(formats strfmt.Registry) error {

	if err := validate.Required("network_vrfs", "body", m.NetworkVrfs); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsMacAddresses(formats strfmt.Registry) error {

	if err := validate.Required("nics_mac_addresses", "body", m.NicsMacAddresses); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsNames(formats strfmt.Registry) error {

	if err := validate.Required("nics_names", "body", m.NicsNames); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsNeighborMacAddresses(formats strfmt.Registry) error {

	if err := validate.Required("nics_neighbor_mac_addresses", "body", m.NicsNeighborMacAddresses); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsNeighborNames(formats strfmt.Registry) error {

	if err := validate.Required("nics_neighbor_names", "body", m.NicsNeighborNames); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsNeighborVrfs(formats strfmt.Registry) error {

	if err := validate.Required("nics_neighbor_vrfs", "body", m.NicsNeighborVrfs); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateNicsVrfs(formats strfmt.Registry) error {

	if err := validate.Required("nics_vrfs", "body", m.NicsVrfs); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validatePartitionID(formats strfmt.Registry) error {

	if err := validate.Required("partition_id", "body", m.PartitionID); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateRackid(formats strfmt.Registry) error {

	if err := validate.Required("rackid", "body", m.Rackid); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateSizeid(formats strfmt.Registry) error {

	if err := validate.Required("sizeid", "body", m.Sizeid); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateStateValue(formats strfmt.Registry) error {

	if err := validate.Required("state_value", "body", m.StateValue); err != nil {
		return err
	}

	return nil
}

func (m *V1FirewallFindRequest) validateTags(formats strfmt.Registry) error {

	if err := validate.Required("tags", "body", m.Tags); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1FirewallFindRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1FirewallFindRequest) UnmarshalBinary(b []byte) error {
	var res V1FirewallFindRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}