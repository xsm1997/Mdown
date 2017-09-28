package Ossdownloader

import (
	"testing"
	"strings"
)

func Test_Download(t *testing.T) {
	//object地址
	src := "http://oss.aliyuncs.com/607afc11/7a70d025-ec54-4fff-ab4e-aef080305645.zip"

	//将http更换为https 防止运营商劫持
	url:=strings.Replace(src,`http://`, `https://`, -1)


	Download(url, "a.zip", 10)
}
