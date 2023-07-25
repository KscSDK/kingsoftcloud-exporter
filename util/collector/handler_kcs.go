package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_KCS = "KCS"
)

func init() {
	registerHandler(Namespace_KCS, defaultHandlerEnabled, NewKCSHandler)
}

type kcsHandler struct {
	baseProductHandler
}

func (h *kcsHandler) GetNamespace() string {
	return Namespace_KCS
}

func NewKCSHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &kcsHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
