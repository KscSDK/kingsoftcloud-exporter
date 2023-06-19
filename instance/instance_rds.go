package instance

import "fmt"

//InstanceRDSMeta
type InstanceRDSMeta struct {
	//实例ID
	DBInstanceIdentifier string `json:"DBInstanceIdentifier"`

	//实例名称
	DBInstanceName string `json:"DBInstanceName"`

	//实例类型。取值范围：HRDS（高可用）,RR（只读实例）,TRDS（临时实例）,SINGLERDS（单实例）
	DBInstanceType string `json:"DBInstanceType"`

	//实例状态。
	DBInstanceStatus string `json:"DBInstanceStatus"`

	//实例Region信息
	Region string `json:"Region"`

	//数据库类型
	Engine string `json:"Engine"`

	//数据库版本号
	EngineVersion string `json:"EngineVersion"`

	//实例虚IP
	Vip string `json:"Vip"`

	//实例创建时间
	CreateTime string `json:"InstanceCreateTime"`
}

//InstancesRDS
type InstanceRDS struct {
	InstanceBase
	meta *InstanceRDSMeta
}

//GetMeta
func (i *InstanceRDS) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceRDS) GetInstanceName() string {
	return i.meta.DBInstanceName
}

//GetInstanceIP
func (i *InstanceRDS) GetInstanceIP() string {
	return i.meta.Vip
}

//GetFieldValueByName
func (i *InstanceRDS) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceRDS) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceRDS
func NewInstanceRDS(instanceId string, meta *InstanceRDSMeta) (*InstanceRDS, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceRDS{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
