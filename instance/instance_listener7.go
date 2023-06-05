package instance

import "fmt"

//InstanceEIP
type InstanceListener7 struct {
	InstanceBase
	meta *InstanceListenerMeta
}

//GetMeta
func (i *InstanceListener7) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceListener7) GetInstanceName() string {
	return i.meta.ListenerName
}

func (i *InstanceListener7) GetInstanceIP() string {
	return `LISTENER7`
}

//GetFieldValueByName
func (i *InstanceListener7) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceListener7) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceListener7
func NewInstanceListener7(instanceId string, meta *InstanceListenerMeta) (*InstanceListener7, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceListener7{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
