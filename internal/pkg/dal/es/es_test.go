package es

import "testing"

//
//import (
//	"context"
//	"fmt"
//	"github.com/elastic/go-elasticsearch/v7"
//	"github.com/elastic/go-elasticsearch/v7/esapi"
//	"github.com/olivere/elastic/v7"
//	"log"
//	"os"
//	"testing"
//	"time"
//)

func TestInitEsClient(t *testing.T) {
	//queryEs := elastic.NewBoolQuery()
	//queryEs = queryEs.Must(elastic.NewMatchQuery("report_id", "924"))
	//
	//Client, _ = elastic.NewClient(
	//	elastic.SetURL("http://172.17.101.191:9200"),
	//	elastic.SetSniff(false),
	//	elastic.SetBasicAuth("elastic", "ZSrfx4R6ICa3skGBpCdf"),
	//	elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
	//	elastic.SetHealthcheckInterval(30*time.Second),
	//)
	//
	//_, _, err := Client.Ping("http://172.17.101.191:9200").Do(context.Background())
	//if err != nil {
	//	fmt.Println("es连接失败", err)
	//	return
	//}

	// 删除index
	//arr := []int{1067}
	//for _, index := range arr {
	//	r, err := Client.DeleteIndex(fmt.Sprintf("%d", index)).Do(context.Background())
	//	if err != nil {
	//		fmt.Println("es删除索引", err)
	//	}
	//	fmt.Println(r)
	//}

	//查询es中的所有index
	//cfg := elasticsearch.Config{
	//	Addresses: []string{
	//		"http://es-cn-7pp2x5zd10001vn8v.elasticsearch.aliyuncs.com:9200",
	//	},
	//	Username: "elastic",
	//	Password: "dnRqw5rRicYjbkkE",
	//}
	//es, err := elasticsearch.NewClient(cfg)
	//if err != nil {
	//	panic(err)
	//}
	//res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), es)
	//if err != nil {
	//	return
	//}
	//defer res.Body.Close()
	//
	//fmt.Println(res.String())

	//indice := Client.IndexGetSettings("924")
	//str, _ := json.Marshal(indice)
	//fmt.Println(indice)
	//mapping := Client.GetMapping()
	//service := mapping.Index(reportId)
	//client.IndexGetSettings()
	//res, err := Client.Search("924").Query(queryEs).Sort("time_stamp", true).Size(100000).Pretty(true).Do(context.Background())
	//if err != nil {
	//	fmt.Println("123", err)
	//}
	//var sceneTestResultDataMsg kao.SceneTestResultDataMsg
	//str := []byte("{\n\t\"end\": false,\n\t\"team_id\": 0,\n\t\"report_id\": \"875\",\n\t\"report_name\": \"123\",\n\t\"plan_id\": 181,\n\t\"plan_name\": \"123\",\n\t\"scene_id\": 198,\n\t\"scene_name\": \"if控制器\",\n\t\"results\": {\n\t\t\"ea75c258-258a-4ddf-89d0-77121d221ca4\": {\n\t\t\t\"event_id\": \"ea75c258-258a-4ddf-89d0-77121d221ca4\",\n\t\t\t\"name\": \"新建接口 \",\n\t\t\t\"plan_id \": 181,\n\t\t\t\"plan_name \": \"123 \",\n\t\t\t\"scene_id \": 198,\n\t\t\t\"scene_name \": \"if控制器 \",\n\t\t\t\"total_request_num \": 3,\n\t\t\t\"total_request_time \": 3350,\n\t\t\t\"success_num \": 3,\n\t\t\t\"error_num \": 0,\n\t\t\t\"error_threshold \": 0,\n\t\t\t\"avg_request_time \": 1116,\n\t\t\t\"max_request_time \": 1124,\n\t\t\t\"min_request_time \": 1111,\n\t\t\t\"custom_request_time_line \": 0,\n\t\t\t\"ninety_request_time_line \": 90,\n\t\t\t\"ninety_five_request_time_line \": 95,\n\t\t\t\"ninety_nine_request_time_line \": 99,\n\t\t\t\"custom_request_time_line_value \": 0,\n\t\t\t\"ninety_request_time_line_value \": 1124,\n\t\t\t\"ninety_five_request_time_line_value \": 1124,\n\t\t\t\"ninety_nine_request_time_line_value \": 1124,\n\t\t\t\"send_bytes \": 105,\n\t\t\t\"received_bytes \": 177,\n\t\t\t\"qps \": 0.9\n\t\t}\n\t},\n\t\"machine \": null,\n\t\"time_stamp \": 1664334772\n}")
	//str := "{\n\t\"end\": false,\n\t\"team_id\": 0,\n\t\"report_id\": \"873\",\n\t\"report_name\": \"test\",\n\t\"plan_id\": 182,\n\t\"plan_name\": \"test\",\n\t\"scene_id\": 199,\n\t\"scene_name\": \"jaskldjklasjd\",\n\t\"results\": {\n\t\t\"56d4c8dd-e1aa-4ac3-a8f8-055d0bc74f67\": {\n\t\t\t\"event_id\": \"56d4c8dd-e1aa-4ac3-a8f8-055d0bc74f67\",\n\t\t\t\"name\": \"新建接口\",\n\t\t\t\"plan_id\": 182,\n\t\t\t\"plan_name\": \"test\",\n\t\t\t\"scene_id\": 199,\n\t\t\t\"scene_name\": \"jaskldjklasjd\",\n\t\t\t\"total_request_num\": 960,\n\t\t\t\"total_request_time\": 357,\n\t\t\t\"success_num\": 0,\n\t\t\t\"error_num\": 960,\n\t\t\t\"error_threshold\": 0,\n\t\t\t\"avg_request_time\": 0,\n\t\t\t\"max_request_time\": 4,\n\t\t\t\"min_request_time\": 0,\n\t\t\t\"custom_request_time_line\": 0,\n\t\t\t\"ninety_request_time_line,\"\n\t\t\tninety_five_request_time_line \":95,\"\n\t\t\tninety_nine_request_time_line \":99,\"\n\t\t\tcustom_request_time_line_value \":0,\"\n\t\t\tninety_request_time_line_value \":1,\"\n\t\t\tninety_five_request_time_line_value \":2,\"\n\t\t\tninety_nine_request_time_line_value \":3,\"\n\t\t\tsend_bytes \":33600,\"\n\t\t\treceived_bytes \":56640,\"\n\t\t\tqps \":2689.08},\"\n\t\t\t72 d6a8c2 - 4 adc - 4950 - 89 a7 - 169516953520 f376 \":{\"\n\t\t\tevent_id \":\"\n\t\t\t72 d6a8c2 - 4 adc - 4950 - 89 a7 - 16953520 f376 \",\"\n\t\t\tname \":\"\n\t\t\t新建接口 \",\"\n\t\t\tplan_id \":182,\"\n\t\t\tplan_name \":\"\n\t\t\ttest \",\"\n\t\t\tscene_id \":199,\"\n\t\t\tscene_name \":\"\n\t\t\tjaskldjklasjd \",\"\n\t\t\ttotal_request_num \":960,\"\n\t\t\ttotal_request_time \":900,\"\n\t\t\tsuccess_num \":0,\"\n\t\t\terror_num \":960,\"\n\t\t\terror_threshold \":0,\"\n\t\t\tavg_request_time \":0,\"\n\t\t\tmax_request_time \":4,\"\n\t\t\tminuest_time \":0,\"\n\t\t\tcustom_request_time_line \":0,\"\n\t\t\tninety_request_time_line \":90,\"\n\t\t\tninety_five_request_time_line \":95,\"\n\t\t\tninety_nine_request_time_line \":99,\"\n\t\t\tcustom_request_time_line_value \":0,\"\n\t\t\tninety_request_time_line_value \":2,\"\n\t\t\tninety_five_request_time_line_value \":3,\"\n\t\t\tninety_nine_request_time_line_value \":4,\"\n\t\t\tsend_bytes \":33:90,\"\n\t\t\tninety_five_request_time_line \":95,\"\n\t\t\tninety_nine_request_time_line \":99,\"\n\t\t\tcustom_request_time_line_value \":0,\"\n\t\t\tninety_request_time_line_value \":1,\"\n\t\t\tninety_five_request_time_line_value \":2,\"\n\t\t\tninety_nine_request_time_line_value \":3,\"\n\t\t\tsend_bytes \":34300,\"\n\t\t\treceived_bytes \":57820,\"\n\t\t\tqps \":2745.1},\"\n\t\t\t72 d6a8c2 - 4 adc - 4950 - 89 a7 - 16953520 f376 \":{\"\n\t\t\tevent_id \":\"\n\t\t\t72 d6a8c2 - 4 adc - 4950 - 89 a7 - 16953520 f376 \",\"\n\t\t\tname \":\"\n\t\t\t新建接口 \",\"\n\t\t\tplan_id \":182,\"\n\t\t\tplan_name \":\"\n\t\t\ttest \",\"\n\t\t\tscene_id \":199,\"\n\t\t\tscene_name \":\"\n\t\t\tjaskldjklasjd \",\"\n\t\t\ttotal_request_num \":980,\"\n\t\t\ttotal_request_time \":900,\"\n\t\t\tsuccess_num \":0,\"\n\t\t\terror_num \":980,\"\n\t\t\terror_threshold \":0,\"\n\t\t\tavg_request_time \":0,\"\n\t\t\tmax_request_time \":4,\"\n\t\t\tmin_rst_time \":0,\"\n\t\t\tcustom_request_time_line \":0,\"\n\t\t\tninety_request_time_line \":90,\"\n\t\t\tninety_five_request_time_line \":95,\"\n\t\t\tninety_nine_request_time_line \":99,\"\n\t\t\tcustom_request_time_line_value \":0,\"\n\t\t\tninety_request_time_line_value \":2,\"\n\t\t\tninety_five_request_time_line_value \":3,\"\n\t\t\tninety_nine_request_time_line_value \":4,\"\n\t\t\tsend_bytes \":34300,\"\n\t\t\treceived_bytes \":57820,\"\n\t\t\tqps \":1088.89}},\"\nmachine \":null,\"\ntime_stamp \":1664334555}"
	//_, err := Client.Index().Index("report").BodyString("123").Do(context.Background())
	//_, err := Client.Delete().Index("report").Do(context.Background())
	//for _, v := range res.Hits.Hits {
	//	str, _ := json.Marshal(v.Source)
	//	fmt.Println(string(str))
	//}

}

func TestPost(t *testing.T) {

	//host := "http://172.17.101.191:9200/"
	//
	//url := host + "925"
	//body := "{\"settings\": {\"index\": {\"max_result_window\": 20000000}}}"
	//err := Post(url, body)
	//if err != nil {
	//
	//}
}
