package instance

import (
	"fmt"
)

//InstanceDCGWMeta
type InstanceDCGWMeta struct {
	DirectConnectGatewayId   string `json:"DirectConnectGatewayId"`
	DirectConnectGatewayName string `json:"DirectConnectGatewayName"`
	CreateTime               string `json:"CreateTime"`
}

//InstanceDCGW
type InstanceDCGW struct {
	InstanceBase
	meta *InstanceDCGWMeta
}

//GetMeta
func (i *InstanceDCGW) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceDCGW) GetInstanceName() string {
	return "eip"
}

func (i *InstanceDCGW) GetInstanceIP() string {
	return `DC_GW`
}

//GetFieldValueByName
func (i *InstanceDCGW) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceDCGW) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceDCGW
func NewInstanceDCGW(instanceId string, meta *InstanceDCGWMeta) (*InstanceDCGW, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceDCGW{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
