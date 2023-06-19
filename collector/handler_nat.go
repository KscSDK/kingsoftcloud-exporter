package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_NAT = "NAT"
)

func init() {
	registerHandler(Namespace_NAT, defaultHandlerEnabled, NewNATHandler)
}

type natHandler struct {
	baseProductHandler
}

func (h *natHandler) GetNamespace() string {
	return Namespace_NAT
}

func NewNATHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &natHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
