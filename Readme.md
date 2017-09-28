## Introduce
使用阿里云OSS的官方golang库,总是出现莫名其妙的问题，unexpected EOF、io error、下载完不能合并块等,官方不能及时修复，
所以简单实现了下,具备以下功能：  <br />
1、可设置协程数  <br />
2、协程出错自动重试(防止unexpected EOF等)

### 注意:<br />
因各运营商有劫持现象，会把oss的文件缓存起来，下载时发生302跳转，直接导致http头中的range无法生效，分块下载失败  <br />
所以需要使用https协议来下载文件


## Install
```
go get -u https://github.com/bryant24/Ossdownloader
```

## How to use
```
package main

import (
	"strings"
	"github.com/bryant24/Ossdownloader"
)

func main() {

	//object地址
	src := "http://example.oss.aliyuncs.com/607afc11/7a70d025-ec54-4fff-ab4e-aef080305645.zip"

	//将http更换为https 防止运营商劫持
	//如果是https协议可以忽略
	url:=strings.Replace(src,`http://`, `https://`, -1)

    //设置为10个协程
	Ossdownloader.Download(url, "a.zip", 10)
}

```
