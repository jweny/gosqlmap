package tools

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

var baseData SinglePageBaseData
var source = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {

}

// 连接性
func checkConnect(conf *ReqConf) (bool, error) {
	Info.Println("testing connection to the target URL")
	statusCode, body, err := httpDoTimeout(conf)
	if err != nil {
		Error.Println("error while connect to:", conf.url, "error:", err)
		return false, err
	}
	if statusCode == 200 {
		baseData = SinglePageBaseData{
			BaseStatusCode: statusCode,
			BaseBodyLength: len(body),
			BaseBody:       body,
		}
		Info.Println(conf.url, "can connect")
		return true, nil
	}
	Warning.Println(conf.url, "can not connect")
	return false, nil
}

// 稳定性
func checkStability(conf *ReqConf) (bool, error) {
	Info.Println("testing if the target URL content is stable")
	secondStatusCode, secondBody, err := httpDoTimeout(conf)
	isStability := false
	if err != nil {
		Error.Println("error while check stability to:", conf.url, "error:", err)
		return false, err
	}
	if secondStatusCode == baseData.BaseStatusCode && checkIsSamePage(secondBody) {
		isStability = true
		Info.Println(conf.url, "content is stable")
	} else {
		// 再尝试一次
		thirdStatusCode, thirdBody, err := httpDoTimeout(conf)
		if err != nil {
			Error.Println("error while check stability to:", conf.url, "error:", err)
			return false, err
		}
		if thirdStatusCode == baseData.BaseStatusCode && checkIsSamePage(thirdBody) {
			isStability = true
			Info.Println(conf.url, "content is stable")
		} else {
			Warning.Println(conf.url, "content is not stable")
		}
	}
	return isStability, nil
}

func getAllGetParams(conf *ReqConf) (url.Values, error) {
	targetUrl, err := url.Parse(conf.url)
	if err != nil {
		return nil, err
	}
	paramMap, err := url.ParseQuery(targetUrl.RawQuery)
	if err != nil {
		return nil, err
	}
	return paramMap, nil
}

func checkParamIsDynamic(conf *ReqConf) (url.Values, error) {
	paramMap, err := getAllGetParams(conf)
	Info.Println("testing parameters is dynamic")
	if err != nil {
		return nil, err
	}
	dynamicList := make(url.Values)
	// 遍历参数 替换payload尝试...
	if len(paramMap) == 0 {
		Warning.Println("not found any param")
	} else {
		// 遍历每一个参数
		for key, value := range paramMap {
			currentUrl := strings.Replace(conf.url, key+"="+value[0], key+"="+genRandom4Num(source), 1)
			currentConf := ReqConf{url: currentUrl}
			_, currentBody, err := httpDoTimeout(&currentConf)
			if err != nil {
				Error.Println(err)
			}
			flag := checkIsSamePage(currentBody)
			if flag == false {
				Info.Println("GET parameter", key, "appears to be dynamic")
				dynamicList[key] = value
			}
		}
	}
	return dynamicList, nil
}

//payload 替换掉第一个参数  waf条件： samePage为假 && 页面包含关键词
func checkWaf(conf *ReqConf) (bool, error) {
	Info.Println("checking if the target is protected by some kind of WAF/IPS")
	paramMap, err := getAllGetParams(conf)
	Info.Println("testing parameters is dynamic")
	if err != nil {
		return false, err
	}
	for key, value := range paramMap {
		currentUrl := strings.Replace(conf.url, key+"="+value[0], key+"="+IPS_WAF_CHECK_PAYLOAD, 1)
		currentConf := ReqConf{url: currentUrl}
		_, currentBody, err := httpDoTimeout(&currentConf)
		for _, wafKeyword := range WAF_CHECK_KEYWORD {
			if strings.Contains(string(currentBody), wafKeyword) {
				Warning.Println("\"heuristics detected that the target is protected by some kind of WAF/IPS")
				Warning.Println("stop to continue with further target testing")
				return true, nil
			}
		}
		if err != nil {
			Error.Println(err)
		}
		break
	}
	return false, nil
}
