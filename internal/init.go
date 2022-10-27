package internal

import (
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/redis"
	log "kp-collector/internal/pkg/log"
)

func InitProjects() {
	conf.MustInitConf()
	log.InitLogger()
	//es.InitEsClient(conf.Conf.ES.Host, conf.Conf.ES.Username, conf.Conf.ES.Password)
	redis.InitRedisClient(conf.Conf.Redis.Address, conf.Conf.Redis.Password, conf.Conf.Redis.DB)
}
