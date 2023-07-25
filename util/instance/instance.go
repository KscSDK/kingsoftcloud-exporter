package instance

// KscInstance 每个产品的实例对象, 可用于配置导出指标的额外label填充, 根据字段名获取值
type KscInstance interface {
	// 获取实例的id
	GetInstanceID() string

	// 获取实例名称
	GetInstanceName() string

	// 获取实例IP
	GetInstanceIP() string

	// 用于查询云监控数据的主键字段, 一般是实例id
	GetMonitorQueryKey() string

	// 根据字段名称获取该字段的值, 由各个产品接口具体实现
	GetFieldValueByName(string) (string, error)

	// 根据字段名称获取该字段的值, 由各个产品接口具体实现
	GetFieldValuesByName(string) (map[string][]string, error)

	// 获取实例raw元数据, 每个实例类型不一样
	GetMeta() interface{}
}

//InstanceBase 基本资源信息
type InstanceBase struct {
	//实例名称
	InstanceName string

	//实例ID
	InstanceID string

	//实例IP
	InstanceIP string

	//实例所属机房
	Region string
}

//GetInstanceID
func (i *InstanceBase) GetInstanceID() string {
	return i.InstanceID
}

//GetInstanceID
func (i *InstanceBase) GetInstanceName() string {
	return i.InstanceName
}

//GetInstanceID
func (i *InstanceBase) GetInstanceIP() string {
	return i.InstanceIP
}

//GetMonitorQueryKey
func (i *InstanceBase) GetMonitorQueryKey() string {
	return i.InstanceID
}
