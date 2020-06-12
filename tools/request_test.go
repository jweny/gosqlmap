package tools

import (
	"testing"
)

func Test(t *testing.T) {
	config := ReqConf{url: "http://192.168.102.135/Less-1/?id=1&name=haha&pass=hehe"}
	start(&config)
}
