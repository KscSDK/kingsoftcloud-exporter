package instance

import "fmt"

var (
	MongoDBType_Sharding   = "Cluster"
	MongoDBType_ReplicaSet = "HighIO"
)

type InstanceMongoMeta struct {
	AvailabilityZone string `json:"Az"`
	InstanceId       string `json:"InstanceId"`
	InstanceName     string `json:"Name"`
	IP               string `json:"IP"`
	InstanceType     string `json:"InstanceType"`
}

//InstanceMongoNodeMeta
type InstanceMongoNodeMeta struct {
	InstanceId   string `json:"NodeId"`
	InstanceName string `json:"Name"`
	IP           string `json:"IP"`
	EIP          string `json:"Eip"`
	Role         string `json:"Role"`
	Status       string `json:"Status"`
}

//InstanceMongo
type InstanceMongo struct {
	InstanceBase
	meta *InstanceMongoNodeMeta
}

//GetMeta
func (i *InstanceMongo) GetMeta() interface{} {
	return i.meta
}

//GetInstanceID
func (i *InstanceMongo) GetInstanceID() string {
	return i.meta.InstanceId
}

//GetInstanceName
func (i *InstanceMongo) GetInstanceName() string {
	return i.meta.InstanceName
}

//GetInstanceIP
func (i *InstanceMongo) GetInstanceIP() string {
	return i.meta.IP
}

//GetFieldValueByName
func (i *InstanceMongo) GetFieldValueByName(string) (string, error) {
	return "", nil
}

func (i *InstanceMongo) GetFieldValuesByName(string) (map[string][]string, error) {
	return nil, nil
}

//NewInstanceEIP
func NewInstanceMongo(instanceId string, meta *InstanceMongoNodeMeta) (*InstanceMongo, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId is empty ")
	}

	if meta == nil {
		return nil, fmt.Errorf("meta is empty ")
	}

	ins := &InstanceMongo{
		InstanceBase: InstanceBase{
			InstanceID: instanceId,
		},
		meta: meta,
	}

	return ins, nil
}
