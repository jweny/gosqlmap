package tools

func start(conf *ReqConf) {
	// 连接性检查，初始化result
	isConnect, err := checkConnect(conf)
	if isConnect == false || err != nil {
		return
	}
	// 稳定性检查，更新result
	isStability, err := checkStability(conf)
	if isStability == false || err != nil {
		return
	}
	// 启发式sql注入检测
	injParam, dbms, err := heuristicCheckSqlInjection(conf)
	if err != nil{
		Error.Println(err)
	}
	Info.Println(injParam, dbms)
}


