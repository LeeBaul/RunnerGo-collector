package internal

import (
	"RunnerGo-collector/internal/pkg/conf"
	"RunnerGo-collector/internal/pkg/dal/redis"
	log "RunnerGo-collector/internal/pkg/log"
)

func InitProjects() {
	conf.MustInitConf()
	log.InitLogger()
	//es.InitEsClient(conf.Conf.ES.Host, conf.Conf.ES.Username, conf.Conf.ES.Password)
	redis.InitRedisClient(conf.Conf.Redis1.Address, conf.Conf.Redis1.Password, conf.Conf.Redis1.DB, conf.Conf.Redis2.Address, conf.Conf.Redis2.Password, conf.Conf.Redis2.DB)
}
