package instance

import (
	"sync"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// 可用于产品的实例的缓存, InstanceRepository
type InstanceCache struct {
	Raw            InstanceRepository
	cache          map[string]KscInstance
	lastReloadTime time.Time
	logger         log.Logger
	mu             sync.Mutex
	reloadInterval time.Duration
}

//GetInstanceKey
func (c *InstanceCache) GetInstanceKey() string {
	return c.Raw.GetInstanceKey()
}

//Get
func (c *InstanceCache) Get(id string) (KscInstance, error) {
	i, exists := c.cache[id]
	if exists {
		return i, nil
	}

	i, err := c.Raw.Get(id)
	if err != nil {
		return nil, err
	}

	c.cache[i.GetInstanceID()] = i
	return i, nil
}

//ListByIds
func (c *InstanceCache) ListByIds(ids []string) (instances []KscInstance, err error) {
	err = c.checkNeedReload()
	if err != nil {
		return nil, err
	}

	var notExists []string
	for _, id := range ids {
		i, ok := c.cache[id]
		if ok {
			instances = append(instances, i)
		} else {
			notExists = append(notExists, id)
		}
	}
	return
}

func (c *InstanceCache) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {
	err = c.checkNeedReload()
	if err != nil {
		return
	}

	for _, ins := range c.cache {
		for k, v := range filters {
			tv, e := ins.GetFieldValueByName(k)
			if e != nil {
				break
			}
			if v != tv {
				break
			}
		}
		instances = append(instances, ins)
	}

	return
}

func (c *InstanceCache) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	err = c.checkNeedReload()
	if err != nil {
		return
	}

	for _, ins := range c.cache {
		for k, v := range filters {
			tv, e := ins.GetFieldValueByName(k)
			if e != nil {
				break
			}
			if v != tv {
				break
			}
		}
		instances = append(instances, ins)
	}

	return
}

//checkNeedReload
func (c *InstanceCache) checkNeedReload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.lastReloadTime.IsZero() && time.Now().Sub(c.lastReloadTime) < c.reloadInterval {
		return nil
	}

	var instances []KscInstance
	var err error
	if config.ExporterRunningMode == config.ExporterMode_Mock {
		instances, err = c.Raw.ListByMonitors(map[string]interface{}{})
	} else {
		instances, err = c.Raw.ListByFilters(map[string]interface{}{})
	}

	if err != nil {
		return err
	}
	numChanged := 0
	if len(instances) > 0 {
		newCache := map[string]KscInstance{}
		for _, instance := range instances {
			newCache[instance.GetInstanceID()] = instance
		}
		numChanged = len(newCache) - len(c.cache)
		c.cache = newCache
	}
	c.lastReloadTime = time.Now()

	level.Info(c.logger).Log("msg", "Reload instance cache", "num", len(c.cache), "changed", numChanged)
	return nil
}

func NewInstanceCache(
	repo InstanceRepository,
	reloadInterval time.Duration,
	logger log.Logger,
) InstanceRepository {
	cache := &InstanceCache{
		Raw:            repo,
		cache:          map[string]KscInstance{},
		reloadInterval: reloadInterval,
		logger:         logger,
	}
	return cache
}
