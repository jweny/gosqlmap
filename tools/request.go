package tools

import (
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"strings"
	"time"
)
type ReqConf struct {
	url string
	method string
	data string
	cookies map[string]interface{}
}

type SinglePageBaseData struct {
	BaseBody       []byte
	BaseBodyLength int
	BaseStatusCode int
}

//
//type Result struct {
//	method string
//	paramMap url.Values
//	isConnect bool
//	isStability bool
//}

var UPPER_RATIO_BOUND = 0.98
var LOWER_RATIO_BOUND = 0.02
var DIFF_TOLERANCE = 0.05

var REQUEST_NUMBER = 0

var scriptUserAgent = "Mozilla/5.0 (Windows NT 6.2; rv:30.0) Gecko/20150101 Firefox/32.0"

var httpClient = fasthttp.Client{
	Name:                          scriptUserAgent,
	ReadTimeout:                   20 * time.Second,
	WriteTimeout:                  20 * time.Second,
	MaxResponseBodySize:           1024 * 1024 * 2,
	DisableHeaderNamesNormalizing: true,
	Dial: func(addr string) (net.Conn, error) {
		return fasthttp.DialTimeout(addr, 5*time.Second)
	},
	TLSConfig: &tls.Config{
		// 不验证 ssl 证书，可能存在风险
		InsecureSkipVerify: true,
	},
}

//var goHTTPClient = http.Client{
//	Timeout: 10 * time.Second,
//}

// GetOriginalBody 返回未压缩的 Body
func GetOriginalBody(response *fasthttp.Response) ([]byte, error) {
	contentEncoding := strings.ToLower(string(response.Header.Peek("Content-Encoding")))
	var body []byte
	var err error
	switch contentEncoding {
	case "", "none", "identity":
		body, err = response.Body(), nil
	case "gzip":
		body, err = response.BodyGunzip()
	case "deflate":
		body, err = response.BodyInflate()
	default:
		// TODO: support `br`
		body, err = []byte{}, fmt.Errorf("unsupported Content-Encoding: %v", contentEncoding)
	}
	return body, err
}

// httpDoTimeout 提供http 请求，返回 响应码 body 和 错误
func httpDoTimeout(conf *ReqConf) (int, []byte, error) {
	// count request
	REQUEST_NUMBER = REQUEST_NUMBER + 1

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	request.SetRequestURI(conf.url)
	request.Header.SetMethod(conf.method)
	//todo 加个随机UA
	request.Header.SetUserAgent("Mozilla/5.0 (Windows NT 6.2; rv:30.0) Gecko/20150101 Firefox/32.0")

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if conf.cookies != nil {
		for key, v := range conf.cookies {
			request.Header.SetCookie(key, v.(string))
		}
	}
	err := httpClient.DoTimeout(request, response, 10*time.Second)
	//TrafficOut.Println("HTTP request",REQUEST_NUMBER,"\n",request.Header.String())
	if err != nil {
		return 0, nil, err
	}
	body, err := GetOriginalBody(response)
	//TrafficOut.Println("HTTP response",REQUEST_NUMBER,"\n",response.Header.String())
	return response.StatusCode(), body, err
}