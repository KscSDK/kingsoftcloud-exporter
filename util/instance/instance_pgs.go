package instance

import "fmt"

//InstancePGSMeta
type InstancePGSMeta struct {
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

//InstancePGS
type InstancePGS struct {
	InstanceBase
	meta *InstancePGSMeta
}

//GetMeta
func (i *InstancePGS) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstancePGS) GetInstanceName() string {
	return i.meta.NatName
}

//GetInstanceIP
func (i *InstancePGS) GetInstanceIP() string {
	return `PGS`
}

//GetFieldValueByName
func (i *InstancePGS) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstancePGS) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstancePGS
func NewInstancePGS(instanceId string, meta *InstancePGSMeta) (*InstancePGS, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstancePGS{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
