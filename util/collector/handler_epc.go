package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_EPC = "EPC"
)

func init() {
	registerHandler(Namespace_EPC, defaultHandlerEnabled, NewEPCHandler)
}

type epcHandler struct {
	baseProductHandler
}

func (h *epcHandler) GetNamespace() string {
	return Namespace_EPC
}

func NewEPCHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &epcHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
