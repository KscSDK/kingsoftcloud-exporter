package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_KEC = "KEC"
)

func init() {
	registerHandler(Namespace_KEC, defaultHandlerEnabled, NewKECHandler)
}

type kecHandler struct {
	baseProductHandler
}

func (h *kecHandler) GetNamespace() string {
	return Namespace_KEC
}

func NewKECHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &kecHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
