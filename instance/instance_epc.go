package instance

import "fmt"

type InstanceEPCMeta struct {
	AvailabilityZone string `json:"AvailabilityZone"`
	HostName         string `json:"HostName"`
	InstanceId       string `json:"InstanceId"`
	InstanceName     string `json:"InstanceName"`
	PrivateIpAddress string `json:"PrivateIpAddress"`
}

//InstanceEPC
type InstanceEPC struct {
	InstanceBase
	meta *InstanceEPCMeta
}

//GetMeta
func (i *InstanceEPC) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceEPC) GetInstanceName() string {
	return i.meta.InstanceName
}

//GetInstanceIP
func (i *InstanceEPC) GetInstanceIP() string {
	return i.meta.PrivateIpAddress
}

//GetFieldValueByName
func (i *InstanceEPC) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceEPC) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceEPC
func NewInstanceEPC(instanceId string, meta *InstanceEPCMeta) (*InstanceEPC, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceEPC{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
