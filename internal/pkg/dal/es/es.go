package es

import (
	"fmt"

	"github.com/olivere/elastic/v7"

	"kp-collector/internal/pkg/conf"
)

var (
	esc *elastic.Client
)

func InitElasticSearchClient() {
	var err error
	esc, err = elastic.NewClient(
		elastic.SetURL(conf.Conf.ES.Host),
		elastic.SetBasicAuth(conf.Conf.ES.Username, conf.Conf.ES.Password))

	if err != nil {
		panic(err)
	}

	fmt.Println("es initialized")
}

func Client() *elastic.Client {
	return esc
}
