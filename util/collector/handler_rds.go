package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_RDS = "KRDS"
)

func init() {
	registerHandler(Namespace_RDS, defaultHandlerEnabled, NewRDSHandler)
}

type rdsHandler struct {
	baseProductHandler
}

func (h *rdsHandler) GetNamespace() string {
	return Namespace_RDS
}

func NewRDSHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &rdsHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
