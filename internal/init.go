package internal

import (
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/es"
)

func InitProjects() {
	conf.MustInitConf()
	es.InitElasticSearchClient()
}
