package tools

import (
	"regexp"
	"strings"
)


var dbmsErrorKeyword map[string]string

// 启发式检测sql注入:发包尝试让Web应用报错,目的为探测该参数点是否是动态的、是否为可能的注入点
func heuristicCheckSqlInjection(conf *ReqConf) (string, string, error){
	// error based 生成字典
	dbmsErrorKeyword = genDbmsErrorFromXml()

	dynamicParams, err := checkParamIsDynamic(conf)
	if err != nil{
		return "", "", err
	}
	if len(dynamicParams) == 0 {
		return "", "", err
	}
	// 遍历每一个动态参数
	for key, value := range dynamicParams {
		// 生成payload 用于报错
		payload := genHeuristicCheckPayload()
		currentUrl := strings.Replace(conf.url, key+"="+value[0], key+"="+value[0]+payload, 1)
		currentConf := ReqConf{url: currentUrl}
		_, currentBody, err := httpDoTimeout(&currentConf)
		if err != nil {
			Error.Println(err)
			return "", "",nil
		}
		flag := checkIsSamePage(currentBody)
		if flag == false {
			Payload.Println(currentUrl)
			dbms := getDBMSBasedOnErrors(currentBody)
			if dbms != ""{
				Info.Println("heuristic (basic) test shows that GET parameter", key, "might be injectable,  possible DBMS:", dbms)
				return key, dbms, nil
			}
		}
	}
	return "", "", nil
}

// 正则匹配出数据库类型
func getDBMSBasedOnErrors(currentBody []byte) string {
	for key , value := range dbmsErrorKeyword{
		match, _ := regexp.MatchString(key, string(currentBody))
		if match == true{
			Debug.Println(key, value)
			return value
		}
	}
	return ""
}

