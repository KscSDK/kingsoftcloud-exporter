package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_PGS = "PGS"
)

func init() {
	registerHandler(Namespace_PGS, defaultHandlerEnabled, NewPGSHandler)
}

type pgsHandler struct {
	baseProductHandler
}

func (h *pgsHandler) GetNamespace() string {
	return Namespace_PGS
}

func NewPGSHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &pgsHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
