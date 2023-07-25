package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_DCGW = "DCGW"
)

func init() {
	registerHandler(Namespace_DCGW, defaultHandlerEnabled, NewDCGWHandler)
}

type dcgwHandler struct {
	baseProductHandler
}

func (h *dcgwHandler) GetNamespace() string {
	return Namespace_DCGW
}

func NewDCGWHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &dcgwHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
