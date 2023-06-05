package instance

import "fmt"

//InstanceSLBMeta
type InstanceSLBMeta struct {

	//负载均衡的ID
	LoadBalancerId string `json:"LoadBalancerId"`

	//负载均衡的名称
	LoadBalancerName string `json:"LoadBalancerName"`

	//负载均衡的状态，开启和关闭状态 可取值:[start|stop]
	LoadBalancerState string `json:"LoadBalancerState"`

	PublicIp string `json:"PublicIp"`

	//负载均衡创建时间
	CreateTime string `json:"CreateTime"`
}

//InstancesSLB
type InstanceSLB struct {
	InstanceBase
	meta *InstanceSLBMeta
}

//GetMeta
func (i *InstanceSLB) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceSLB) GetInstanceName() string {
	return i.meta.LoadBalancerName
}

//GetInstanceIP
func (i *InstanceSLB) GetInstanceIP() string {
	return "SLB"
}

//GetFieldValueByName
func (i *InstanceSLB) GetFieldValueByName(string) (string, error) {
	return "", nil
}

//GetFieldValuesByName
func (i *InstanceSLB) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceSLB
func NewInstanceSLB(instanceId string, meta *InstanceSLBMeta) (*InstanceSLB, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceSLB{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
