package kao

// ResultDataMsg 请求结果数据结构
type ResultDataMsg struct {
	MachineIp             string `json:"machineIp" bson:"machineIp"`
	MachineName           string `json:"machineName" bson:"machineName"`
	ReportId              string `json:"reportId" bson:"reportId"`
	ReportName            string
	EventId               string
	PlanId                string // 任务ID
	PlanName              string //
	SceneId               string // 场景
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
}
