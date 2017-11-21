## Introduce
使用阿里云OSS的官方golang库下载文件,总是出现莫名其妙的问题，unexpected EOF、io error、下载完不能合并块、不能重试下载等,官方不能及时修复，
所以简单实现了下,具备以下功能：  <br />
1、支持任意http的分片多线程下载，不局限于oss  <br />
2、可设置协程数  <br />
3、协程出错30秒超时自动重试(防止unexpected EOF等)  <br />
4、使用fasthttp作为http client


### 注意:<br />
因个别运营商有劫持现象，会把oss的文件缓存起来，下载时发生302跳转，直接导致http头中的range无法生效，分块下载失败  <br />
所以需要使用https协议来下载文件


## Install
```
go get -u https://github.com/bryant24/Mdown
```

## How to use
```
package main

import (
	"strings"
	"github.com/bryant24/Mdown"
)

func Test_Download(t *testing.T) {

	src := "http://oss.aliyuncs.com/abc.zip"

	//将http更换为https 防止运营商劫持使用缓存
	url:=strings.Replace(src,`http://`, `https://`, -1)

	var timeout time.Duration
	timeout = 30 //分片超时时间
	Download(url, "a.zip", timeout)
}

```
