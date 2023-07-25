package collector

import (
	"github.com/go-kit/log"
)

const (
	Namespace_KS3 = "KS3"
)

func init() {
	registerHandler(Namespace_KS3, defaultHandlerEnabled, NewKS3Handler)
}

type ks3Handler struct {
	baseProductHandler
}

func (h *ks3Handler) GetNamespace() string {
	return Namespace_KS3
}

func NewKS3Handler(c *KscProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &ks3Handler{
		baseProductHandler{
			collector: c,
			logger:    logger,
		},
	}
	return
}
