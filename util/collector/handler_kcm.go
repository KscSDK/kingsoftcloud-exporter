package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_KCM = "KCM"
)

func init() {
	registerHandler(Namespace_KCM, defaultHandlerEnabled, NewKCMHandler)
}

type kcmHandler struct {
	baseProductHandler
}

func (h *kcmHandler) GetNamespace() string {
	return Namespace_KCM
}

func NewKCMHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &kcmHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
