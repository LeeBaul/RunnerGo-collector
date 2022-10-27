package es

//var Client *elastic.Client
//var mapping = "{\"settings\": {\"index\": {\"max_result_window\": 20000000}}}"
//
//func InitEsClient(host, user, password string) {
//	Client, _ = elastic.NewClient(
//		elastic.SetURL(host),
//		elastic.SetSniff(false),
//		elastic.SetBasicAuth(user, password),
//		elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
//		elastic.SetHealthcheckInterval(30*time.Second),
//	)
//	_, _, err := Client.Ping(host).Do(context.Background())
//	if err != nil {
//		panic(fmt.Sprintf("es连接失败: %s", err))
//	}
//	return
//}

//func InsertTestData(sceneTestResultDataMsg kao.SceneTestResultDataMsg) (err error) {
//
//	//index := conf.Conf.ES.Index
//	index := sceneTestResultDataMsg.ReportId
//	exist, err := Client.IndexExists(index).Do(context.Background())
//	if err != nil {
//		log2.Logger.Error(fmt.Sprintf("es连接失败: %s", err))
//		SendStopMsg(conf.Conf.GRPC.Host, sceneTestResultDataMsg.ReportId)
//		return
//	}
//	if !exist {
//		_, clientErr := Client.CreateIndex(index).BodyString(mapping).Do(context.Background())
//		if clientErr != nil {
//			log2.Logger.Error("es创建索引", index, "失败", err)
//			SendStopMsg(conf.Conf.GRPC.Host, sceneTestResultDataMsg.ReportId)
//			return
//		}
//
//	}
//
//	_, err = Client.Index().Index(index).BodyJson(sceneTestResultDataMsg).Do(context.Background())
//	if err != nil {
//		log2.Logger.Error("es写入数据失败", err)
//		log2.Logger.Debug("sceneTestResultDataMsg.ReportId", sceneTestResultDataMsg)
//		SendStopMsg(conf.Conf.GRPC.Host, sceneTestResultDataMsg.ReportId)
//		return
//	}
//	if sceneTestResultDataMsg.End {
//		SendStopMsg(conf.Conf.GRPC.Host, sceneTestResultDataMsg.ReportId)
//	}
//	return
//}
