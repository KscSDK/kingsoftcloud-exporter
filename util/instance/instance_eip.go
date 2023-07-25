package instance

import (
	"fmt"
)

//InstanceEIPMeta
type InstanceEIPMeta struct {
	//弹性IP的ID
	AllocationId string `json:"AllocationId"`

	//弹性IP的带宽
	BandWidth int64 `json:"BandWidth"`

	//弹性 IP 的计费方式
	ChargeType string `json:"ChargeType"`

	//IP 版本
	IpVersion string `json:"IpVersion"`

	//弹性 IP 的线路类型的 ID
	LineId string `json:"LineId"`

	Mode string `json:"Mode"`

	//项目的ID
	ProjectId string `json:"ProjectId"`

	//弹性IP
	PublicIp string `json:"PublicIp"`

	//弹性IP的状态，已绑定(associate)，未绑定(disassociate)
	State string `json:"State"`

	//弹性IP创建时间
	CreateTime string `json:"CreateTime"`
}

//InstanceEIP
type InstanceEIP struct {
	InstanceBase
	meta *InstanceEIPMeta
}

//GetMeta
func (i *InstanceEIP) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceEIP) GetInstanceName() string {
	return "eip"
}

func (i *InstanceEIP) GetInstanceIP() string {
	return i.meta.PublicIp
}

//GetFieldValueByName
func (i *InstanceEIP) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceEIP) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceEIP
func NewInstanceEIP(instanceId string, meta *InstanceEIPMeta) (*InstanceEIP, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceEIP{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
