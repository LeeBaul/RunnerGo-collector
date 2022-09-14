package handler

import (
	"context"
	"encoding/json"

	"github.com/davecgh/go-spew/spew"

	"kp-collector/internal/pkg/dal/es"
	"kp-collector/internal/pkg/dal/kao"
)

func SaveEs(ctx context.Context, value []byte) {
	spew.Dump(value)

	var ret kao.ResultDataMsg
	if err := json.Unmarshal(value, &ret); err != nil {

	}

	//json.Marshal(ret)

	resp, err := es.Client().Index().
		Index("report").BodyJson(ret).Do(ctx)

	spew.Dump(resp, err)

}
