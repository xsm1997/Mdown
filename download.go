package Ossdownloader

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"log"
	"os"
	"github.com/valyala/fasthttp"
	"time"
)

var wg sync.WaitGroup

func Download(src, target string, timeout time.Duration) (err error) {
	var res *http.Response
	for {
		res, err = http.Head(src)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])
	if err != nil {
		return
	}

	//根据文件大小分线程数量 10m=1,55m=5  100m=10 1g=100
	thread := length / 10485760
	if thread == 0 {
		thread = 1
	}

	len_sub := length / thread
	diff := length % thread
	body := make([][]byte, thread)
	for i := 0; i < thread; i++ {
		wg.Add(1)
		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range
		if (i == thread-1) {
			max += diff
		}
		go func(min int, max int, i int) {
			for {
				req := fasthttp.AcquireRequest()
				req.SetRequestURI(src)
				req.Header.SetByteRange(min, max-1)
				resp := fasthttp.AcquireResponse()
				client := &fasthttp.Client{}
				e := client.DoTimeout(req, resp, time.Second*timeout)
				if e != nil {
					log.Println("request error,retry:", e)
					continue
				}

				body[i] = resp.Body()
				e = ioutil.WriteFile(target+"."+strconv.Itoa(i), body[i], 777)
				if e != nil {
					log.Println("retry write file:", e)
					continue
				}
				break
			}
			wg.Done()
		}(min, max, i)
	}
	wg.Wait()

	//merge chunk files to target file
	os.Remove(target)
	f, _ := os.OpenFile(target, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 777)
	for j := 0; j < thread; j++ {
		content, _ := ioutil.ReadFile(target + "." + strconv.Itoa(j))
		f.Write(content)
		os.Remove(target + "." + strconv.Itoa(j))
	}
	f.Close()
	return
}
