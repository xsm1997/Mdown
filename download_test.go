package Mdown

import (
	"testing"
	"strings"
	"time"
)

func Test_Download(t *testing.T) {
	//object地址

	src := "http://oss.aliyuncs.com/abc.zip"

	//将http更换为https 防止运营商劫持使用缓存
	url:=strings.Replace(src,`http://`, `https://`, -1)
	var timeout time.Duration
	timeout = 30
	Download(url, "a.zip", timeout)
}
