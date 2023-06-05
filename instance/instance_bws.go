package instance

import "fmt"

type InstanceBWSMeta struct {
	//共享带宽的ID
	BandWidthShareId string `json:"BandWidthShareId"`

	//共享带宽的名称
	BandWidthShareName string `json:"BandWidthShareName"`

	//共享带宽创建时间
	CreateTime string `json:"CreateTime"`
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
	return i.meta.BandWidthShareName
}

//GetInstanceIP
func (i *InstanceBWS) GetInstanceIP() string {
	return `BWS`
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
