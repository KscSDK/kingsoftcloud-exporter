package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_SLB = "SLB"
)

func init() {
	registerHandler(Namespace_SLB, defaultHandlerEnabled, NewSLBHandler)
}

type slbHandler struct {
	baseProductHandler
}

func (h *slbHandler) GetNamespace() string {
	return Namespace_SLB
}

func NewSLBHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &slbHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
