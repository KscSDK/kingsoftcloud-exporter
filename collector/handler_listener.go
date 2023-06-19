package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_LISTENER = "LISTENER"
)

func init() {
	registerHandler(Namespace_LISTENER, defaultHandlerEnabled, NewListenerHandler)
}

type listenerHandler struct {
	baseProductHandler
}

func (h *listenerHandler) GetNamespace() string {
	return Namespace_LISTENER
}

func NewListenerHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &listenerHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
