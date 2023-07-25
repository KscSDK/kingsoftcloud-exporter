package instance

import "fmt"

type InstanceKCSMeta struct {
	//缓存服务ID
	CacheId string `json:"cacheId"`

	//缓存服务名称
	Name string `json:"name"`

	//缓存服务端口
	Port int `json:"port"`

	//实例版本
	Protocol string `json:"protocol"`

	//缓存服务VIP
	Vip string `json:"vip"`

	//公网IP
	Eip string `json:"eip"`

	//缓存服务模式， 取值范围：[1: 集群(Cluster) | 2: 单主从(Single) | 3	自定义集群]
	Mode int `json:"mode"`

	//机房信息
	Region string `json:"region"`

	//缓存服务创建时间
	CreateTime string `json:"createTime"`
}

//InstanceKCS
type InstanceKCS struct {
	InstanceBase
	meta *InstanceKCSMeta
}

//GetMeta
func (i *InstanceKCS) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceKCS) GetInstanceName() string {
	return i.meta.Name
}

//GetInstanceIP
func (i *InstanceKCS) GetInstanceIP() string {
	return i.meta.Vip
}

//GetFieldValueByName
func (i *InstanceKCS) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceKCS) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceKCS
func NewInstanceKCS(instanceId string, meta *InstanceKCSMeta) (*InstanceKCS, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceKCS{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
