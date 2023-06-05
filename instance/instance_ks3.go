package instance

import "fmt"

type InstancesKS3BucketMeta struct {
	CreationDate string `json:"CreationDate"`
	Name         string `json:"Name"`
	Region       string `json:"Region"`
	Type         string `json:"Type"`
}

//InstanceKS3
type InstanceKS3 struct {
	InstanceBase
	meta *InstancesKS3BucketMeta
}

//GetMeta
func (i *InstanceKS3) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceKS3) GetInstanceName() string {
	return i.meta.Name
}

func (i *InstanceKS3) GetInstanceIP() string {
	return "KS3"
}

//GetFieldValueByName
func (i *InstanceKS3) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceKS3) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceKS3
func NewInstanceKS3(instanceId string, meta *InstancesKS3BucketMeta) (*InstanceKS3, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceKS3{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
