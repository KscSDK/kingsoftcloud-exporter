package instance

import "fmt"

//InstancePEERMeta
type InstancePEERMeta struct {
	//发起端VPC的对等连接ID
	VpcPeeringConnectionId string `json:"VpcPeeringConnectionId"`

	//VPC的对等连接类型, 有效值：local| remote
	VpcPeeringConnectionType string `json:"VpcPeeringConnectionType"`

	//peering的名称
	PeeringName string `json:"PeeringName"`

	//peering的状态
	State string `json:"State"`

	//VPC创建时间
	CreateTime string `json:"CreateTime"`
}

//InstancePEER
type InstancePEER struct {
	InstanceBase
	meta *InstancePEERMeta
}

//GetMeta
func (i *InstancePEER) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstancePEER) GetInstanceName() string {
	return i.meta.PeeringName
}

//GetInstanceIP
func (i *InstancePEER) GetInstanceIP() string {
	return `PEER`
}

//GetFieldValueByName
func (i *InstancePEER) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstancePEER) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstancePEER
func NewInstancePEER(instanceId string, meta *InstancePEERMeta) (*InstancePEER, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstancePEER{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
