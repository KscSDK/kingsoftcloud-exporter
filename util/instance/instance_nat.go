package instance

import "fmt"

//InstanceNATMeta
type InstanceNATMeta struct {
	//VPC的ID
	VpcId string `json:"VpcId"`

	//NAT的ID
	NatId string `json:"NatId"`

	//NAT的名称
	NatName string `json:"NatName"`

	//NAT的作用范围，VPC是指NAT对整个VPC有效，subnet是指NAT对关联的子网有效
	NatMode string `json:"NatMode"`

	//NAT的类型
	NatType string `json:"NatType"`

	//NAT的IP数量
	NatIpNumber int64 `json:"NatIpNumber"`

	//NAT的带宽
	BandWidth int64 `json:"BandWidth"`

	//VPC创建时间
	CreateTime string `json:"CreateTime"`
}

//InstanceNAT
type InstanceNAT struct {
	InstanceBase
	meta *InstanceNATMeta
}

//GetMeta
func (i *InstanceNAT) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceNAT) GetInstanceName() string {
	return i.meta.NatName
}

//GetInstanceIP
func (i *InstanceNAT) GetInstanceIP() string {
	return `NAT`
}

//GetFieldValueByName
func (i *InstanceNAT) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceNAT) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceNAT
func NewInstanceNAT(instanceId string, meta *InstanceNATMeta) (*InstanceNAT, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceNAT{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
