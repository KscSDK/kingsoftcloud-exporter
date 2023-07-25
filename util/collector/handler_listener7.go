package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_LISTENER7 = "LISTENER7"
)

func init() {
	registerHandler(Namespace_LISTENER7, defaultHandlerEnabled, NewListener7Handler)
}

type listener7Handler struct {
	baseProductHandler
}

func (h *listener7Handler) GetNamespace() string {
	return Namespace_LISTENER7
}

func NewListener7Handler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &listener7Handler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
