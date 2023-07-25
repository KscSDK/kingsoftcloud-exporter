package instance

import "fmt"

type InstanceListenerMeta struct {
	//监听器的ID
	ListenerId string `json:"ListenerId"`

	//监听器的名称
	ListenerName string `json:"ListenerName"`

	//监听器的协议
	ListenerProtocol string `json:"ListenerProtocol"`

	//监听器的协议端口
	ListenerPort int `json:"ListenerPort"`

	//监听器的状态 可取值: start | stop
	ListenerState string `json:"ListenerState"`

	//监听器创建时间
	CreateTime string `json:"CreateTime"`
}

//InstanceEIP
type InstanceListener struct {
	InstanceBase
	meta *InstanceListenerMeta
}

//GetMeta
func (i *InstanceListener) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceListener) GetInstanceName() string {
	return i.meta.ListenerName
}

func (i *InstanceListener) GetInstanceIP() string {
	return `LISTENER`
}

//GetFieldValueByName
func (i *InstanceListener) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceListener) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceListener
func NewInstanceListener(instanceId string, meta *InstanceListenerMeta) (*InstanceListener, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceListener{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
