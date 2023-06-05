package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_EIP = "EIP"
)

func init() {
	registerHandler(Namespace_EIP, defaultHandlerEnabled, NewEIPHandler)
}

type eipHandler struct {
	baseProductHandler
}

func (h *eipHandler) GetNamespace() string {
	return Namespace_EIP
}

func NewEIPHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &eipHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
