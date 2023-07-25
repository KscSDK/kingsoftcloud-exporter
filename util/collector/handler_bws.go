package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_BWS = "BWS"
)

func init() {
	registerHandler(Namespace_BWS, defaultHandlerEnabled, NewBWSHandler)
}

type bwsHandler struct {
	baseProductHandler
}

func (h *bwsHandler) GetNamespace() string {
	return Namespace_BWS
}

func NewBWSHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &bwsHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
