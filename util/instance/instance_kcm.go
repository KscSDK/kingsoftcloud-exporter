package instance

import (
	"fmt"
	"strings"
)

type InstancesKCMMeta struct {
	InstanceId string `json:"instanceId"`
}

//InstanceKCM
type InstanceKCM struct {
	InstanceBase
	meta *InstancesKCMMeta
}

//GetMeta
func (i *InstanceKCM) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceKCM) GetInstanceName() string {
	return i.meta.InstanceId
}

func (i *InstanceKCM) GetInstanceIP() string {
	return strings.SplitN(i.meta.InstanceId, ":", 2)[0]
}

//GetFieldValueByName
func (i *InstanceKCM) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceKCM) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceKCM
func NewInstanceKCM(instanceId string, meta *InstancesKCMMeta) (*InstanceKCM, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceKCM{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
