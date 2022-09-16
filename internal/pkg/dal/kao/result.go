package kao

// ResultDataMsg 请求结果数据结构
type ResultDataMsg struct {
	MachineIp             string `json:"machineIp" bson:"machineIp"`
	MachineName           string `json:"machineName" bson:"machineName"`
	ReportId              string `json:"reportId" bson:"reportId"`
	ReportName            string
	EventId               string
	PlanId                int64  // 任务ID
	PlanName              string //
	SceneId               int64  // 场景
	SceneName             string
	ApiId                 string  // 接口ID
	ApiName               string  // 接口名称
	RequestTime           uint64  // 请求响应时间
	CustomRequestTimeLine int64   // 自定义响应时间线
	ErrorThreshold        float32 // 自定义错误率
	ErrorType             int64   // 错误类型：1. 请求错误；2. 断言错误
	IsSucceed             bool    // 请求是否有错：true / false   为了计数
	ErrorMsg              string  // 错误信息
	SendBytes             uint64  // 发送字节数
	ReceivedBytes         uint64  // 接收字节数
	Timestamp             int64
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
type ApiTestResultDataMsg struct {
	ReportId                   string `json:"report_id" bson:"report_id"`
	ReportName                 string `json:"report_name" bson:"report_name"`
	PlanId                     int64  `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName                   string `json:"plan_name" bson:"plan_name"` //
	SceneId                    int64  `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName                  string `json:"scene_name" bson:"scene_name"`
	TargetId                   int64  `json:"target_id" bson:"target_id"`                   // 接口ID
	Name                       string `json:"name" bson:"name"`                             // 接口名称
	TotalRequestNum            uint64 `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime           uint64 `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                 uint64 `json:"success_num" bson:"success_num"`
	ErrorNum                   uint64 `json:"error_num" bson:"error_num"`               // 错误数
	AvgRequestTime             uint64 `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime             uint64 `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime             uint64 `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine      uint64 `json:"custom_request_time_line" bson:"custom_request_time_line"`
	CustomRequestTimeLineValue int64  `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLine      uint64 `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine  uint64 `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine  uint64 `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	SendBytes                  uint64 `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes              uint64 `json:"received_bytes" bson:"received_bytes"` // 接收字节数
}

// SceneTestResultDataMsg 场景的测试结果

type SceneTestResultDataMsg struct {
	ReportId   string                                `json:"report_id" bson:"report_id"`
	ReportName string                                `json:"report_name" bson:"report_name"`
	PlanId     int64                                 `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName   string                                `json:"plan_name" bson:"plan_name"` //
	SceneId    int64                                 `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName  string                                `json:"scene_name" bson:"scene_name"`
	Results    map[interface{}]*ApiTestResultDataMsg `json:"results" bson:"results"`
}

type RequestTimeList []uint64

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
func TimeLineCalculate(line int64, requestTimeList RequestTimeList) (requestTime uint64) {
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
