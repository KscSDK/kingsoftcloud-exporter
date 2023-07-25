package instance

import "fmt"

type InstancesKECMeta struct {
	AvailabilityZone string `json:"AvailabilityZone"`
	HostName         string `json:"HostName"`
	InstanceId       string `json:"InstanceId"`
	InstanceName     string `json:"InstanceName"`
	PrivateIpAddress string `json:"PrivateIpAddress"`
}

//InstanceEIP
type InstanceKEC struct {
	InstanceBase
	meta *InstancesKECMeta
}

//GetMeta
func (i *InstanceKEC) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceKEC) GetInstanceName() string {
	return i.meta.InstanceName
}

func (i *InstanceKEC) GetInstanceIP() string {
	return i.meta.PrivateIpAddress
}

//GetFieldValueByName
func (i *InstanceKEC) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceKEC) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceEIP
func NewInstanceKEC(instanceId string, meta *InstancesKECMeta) (*InstanceKEC, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceKEC{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
