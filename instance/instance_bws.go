package instance

import "fmt"

type InstanceBWSMeta struct {
	AvailabilityZone string `json:"AvailabilityZone"`
	HostName         string `json:"HostName"`
	InstanceId       string `json:"InstanceId"`
	InstanceName     string `json:"InstanceName"`
	PrivateIpAddress string `json:"PrivateIpAddress"`
}

//InstanceBWS
type InstanceBWS struct {
	InstanceBase
	meta *InstanceBWSMeta
}

//GetMeta
func (i *InstanceBWS) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceBWS) GetInstanceName() string {
	return i.meta.InstanceName
}

//GetInstanceIP
func (i *InstanceBWS) GetInstanceIP() string {
	return i.meta.PrivateIpAddress
}

//GetFieldValueByName
func (i *InstanceBWS) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceBWS) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceBWS
func NewInstanceBWS(instanceId string, meta *InstanceBWSMeta) (*InstanceBWS, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceBWS{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
