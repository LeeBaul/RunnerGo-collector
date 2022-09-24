package kao

// ResultDataMsg 请求结果数据结构
type ResultDataMsg struct {
	End                   bool    `json:"end" bson:"end"` // 结束标记
	Name                  string  `json:"name" bson:"name"`
	TeamId                int64   `json:"team_id" bson:"team_id"`
	ReportId              string  `json:"report_id" bson:"report_id"`
	ReportName            string  `json:"report_name" bson:"report_name"`
	MachineNum            int64   `json:"machine_num" bson:"machine_num"` // 使用的压力机数量
	MachineIp             string  `json:"machine_ip" bson:"machine_ip"`
	EventId               string  `json:"event_id" bson:"event_id"`
	PlanId                int64   `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName              string  `json:"plan_name" bson:"plan_name"` //
	SceneId               int64   `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName             string  `json:"sceneName" bson:"scene_name"`
	RequestTime           int64   `json:"request_time" bson:"request_time"`                         // 请求响应时间
	CustomRequestTimeLine int64   `json:"custom_request_time_line" bson:"custom_request_time_line"` // 自定义响应时间线
	ErrorThreshold        float64 `json:"error_threshold" bson:"error_threshold"`                   // 自定义错误率
	ErrorType             int64   `json:"error_type" bson:"error_type"`                             // 错误类型：1. 请求错误；2. 断言错误
	IsSucceed             bool    `json:"is_succeed" bson:"is_succeed"`                             // 请求是否有错：true / false   为了计数
	ErrorMsg              string  `json:"error_msg" bson:"error_msg"`                               // 错误信息
	SendBytes             uint64  `json:"send_bytes" bson:"send_bytes"`                             // 发送字节数
	ReceivedBytes         uint64  `json:"received_bytes" bson:"received_bytes"`                     // 接收字节数
	Timestamp             int64   `json:"timestamp" bson:"timestamp"`
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
type ApiTestResultDataMsg struct {
	EventId                        string  `json:"event_id" bson:"event_id"`
	Name                           string  `json:"name" bson:"name"`
	PlanId                         int64   `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName                       string  `json:"plan_name" bson:"plan_name"` //
	SceneId                        int64   `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName                      string  `json:"scene_name" bson:"scene_name"`
	TotalRequestNum                int64   `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime               int64   `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                     int64   `json:"success_num" bson:"success_num"`
	ErrorNum                       int64   `json:"error_num" bson:"error_num"`               // 错误数
	ErrorThreshold                 float64 `json:"error_threshold" bson:"error_threshold"`   // 自定义错误率
	AvgRequestTime                 int64   `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime                 int64   `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime                 int64   `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine          int64   `json:"custom_request_time_line" bson:"custom_request_time_line"`
	NinetyRequestTimeLine          int64   `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine      int64   `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine      int64   `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	CustomRequestTimeLineValue     int64   `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLineValue     int64   `json:"ninety_request_time_line_value" bson:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue int64   `json:"ninety_five_request_time_line_value" bson:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue int64   `json:"ninety_nine_request_time_line_value" bson:"ninety_nine_request_time_line_value"`
	SendBytes                      uint64  `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes                  uint64  `json:"received_bytes" bson:"received_bytes"` // 接收字节数
	Qps                            float64 `json:"qps" bson:"qps"`
}

// SceneTestResultDataMsg 场景的测试结果

type SceneTestResultDataMsg struct {
	End        bool                             `json:"end" bson:"end"`
	TeamId     int64                            `json:"team_id" bson:"team_id"`
	ReportId   string                           `json:"report_id" bson:"report_id"`
	ReportName string                           `json:"report_name" bson:"report_name"`
	PlanId     int64                            `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName   string                           `json:"plan_name" bson:"plan_name"` //
	SceneId    int64                            `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName  string                           `json:"scene_name" bson:"scene_name"`
	Results    map[string]*ApiTestResultDataMsg `json:"results" bson:"results"`
	Machine    map[string]int64                 `json:"machine" bson:"machine"`
	TimeStamp  int64                            `json:"time_stamp" bson:"time_stamp"`
}
type RequestTimeList []int64

func (rt RequestTimeList) Len() int {
	return len(rt)
}

func (rt RequestTimeList) Less(i int, j int) bool {
	return rt[i] < rt[j]
}
func (rt RequestTimeList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}

// TimeLineCalculate 根据响应时间线，计算该线的值
func TimeLineCalculate(line int64, requestTimeList RequestTimeList) (requestTime int64) {
	if line > 0 && line < 100 {
		proportion := float64(line) / 100
		value := proportion * float64(len(requestTimeList))
		requestTime = requestTimeList[int(value)]
	}
	return

}

type ResultDataMsgList []ResultDataMsg

func (rt ResultDataMsgList) Len() int {
	return len(rt)
}

func (rt ResultDataMsgList) Less(i int, j int) bool {
	return rt[i].Timestamp < rt[j].Timestamp
}
func (rt ResultDataMsgList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}
