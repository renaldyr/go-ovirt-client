package ovirtclient

import (
	"fmt"
	"net"

	ovirtsdk "github.com/renaldyr/go-ovirt"
)

func (o *oVirtClient) UpdateNIC(
	vmid VMID,
	nicID NICID,
	params UpdateNICParameters,
	retries ...RetryStrategy,
) (result NIC, err error) {
	req := o.conn.SystemService().VmsService().VmService(string(vmid)).NicsService().NicService(string(nicID)).Update()

	nicBuilder := ovirtsdk.NewNicBuilder().Id(string(nicID))
	if name := params.Name(); name != nil {
		nicBuilder.Name(*name)
	}
	if vnicProfileID := params.VNICProfileID(); vnicProfileID != nil {
		nicBuilder.VnicProfile(ovirtsdk.NewVnicProfileBuilder().Id(string(*vnicProfileID)).MustBuild())
	}
	if mac := params.Mac(); mac != nil {
		if _, err := net.ParseMAC(*mac); err != nil {
			return nil, newError(EUnidentified, "Failed to parse MacAddress: %s", *mac)
		}
		nicBuilder.Mac(ovirtsdk.NewMacBuilder().Address(*mac).MustBuild())
	}

	req.Nic(nicBuilder.MustBuild())

	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("updating NIC %s for VM %s", nicID, vmid),
		o.logger,
		retries,
		func() error {
			update, err := req.Send()
			if err != nil {
				return wrap(err, EUnidentified, "Failed to update NIC %s", nicID)
			}
			sdkNIC, ok := update.Nic()
			if !ok {
				return newFieldNotFound("NIC update response", "NIC")
			}
			nic, err := convertSDKNIC(sdkNIC, o)
			if err != nil {
				return err
			}
			result = nic
			return nil
		})
	return result, err
}

func (m *mockClient) UpdateNIC(vmid VMID, nicID NICID, params UpdateNICParameters, retries ...RetryStrategy) (
	NIC,
	error,
) {
	m.lock.Lock()
	defer m.lock.Unlock()
	nic, ok := m.nics[nicID]
	if !ok {
		return nil, newError(ENotFound, "NIC not found")
	}
	if nic.vmid != vmid {
		return nil, newError(ENotFound, "NIC not found")
	}
	if name := params.Name(); name != nil {
		nic = nic.withName(*name)
	}
	if vnicProfileID := params.VNICProfileID(); vnicProfileID != nil {
		if _, ok := m.vnicProfiles[*vnicProfileID]; !ok {
			return nil, newError(ENotFound, "VNIC profile %s not found", *vnicProfileID)
		}
		nic = nic.withVNICProfileID(*vnicProfileID)
	}
	if mac := params.Mac(); mac != nil {
		if _, err := net.ParseMAC(*mac); err != nil {
			return nil, newError(EUnidentified, "Failed to parse MacAddress: %s", *mac)
		}
		nic = nic.withMac(*mac)
	}
	m.nics[nicID] = nic

	return nic, nil
}
