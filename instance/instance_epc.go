package instance

import "fmt"

//安全组信息
type NetworkInterfaceAttribute struct {
	//网卡的类型，主网卡(primary)、从网卡(extension)
	NetworkInterfaceType string `json:"NetworkInterfaceType"`

	//服务器的网卡在VPC中的IP
	PrivateIpAddress string `json:"PrivateIpAddress"`
}

type InstanceEPCMeta struct {
	//可用区的名称
	AvailabilityZone string `json:"AvailabilityZone"`

	//裸金属服务器资源ID
	HostId string `json:"HostId"`

	//裸金属服务器名称
	HostName string `json:"HostName"`

	//裸金属服务器机型
	HostType string `json:"HostType"`

	//关联的网卡信息
	NetworkInterfaceAttributeSet []*NetworkInterfaceAttribute `json:"NetworkInterfaceAttributeSet"`

	//创建时间
	CreateTime string `json:"CreateTime"`
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
	return i.meta.HostName
}

//GetInstanceIP
func (i *InstanceEPC) GetInstanceIP() string {
	if len(i.meta.NetworkInterfaceAttributeSet) <= 0 {
		return "EPC"
	}
	return i.meta.NetworkInterfaceAttributeSet[0].PrivateIpAddress
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
