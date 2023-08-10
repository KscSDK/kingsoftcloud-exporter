package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_MONGO = "MONGO"
)

func init() {
	registerHandler(Namespace_MONGO, defaultHandlerEnabled, NewMongoDBHandler)
}

type mongoDBHandler struct {
	baseProductHandler
}

func (h *mongoDBHandler) GetNamespace() string {
	return Namespace_MONGO
}

func NewMongoDBHandler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &mongoDBHandler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
