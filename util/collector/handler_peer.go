package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_PEER = "PEER"
)

func init() {
	registerHandler(Namespace_PEER, defaultHandlerEnabled, NewPEERHandler)
}

type peerHandler struct {
	baseProductHandler
}

func (h *peerHandler) GetNamespace() string {
	return Namespace_PEER
}

func NewPEERHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &peerHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
