package Mdown

import (
	"net/http"
	"strconv"
	"sync"
	"log"
	"os"
	"time"
	"io"
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
	for i := 0; i < thread; i++ {
		wg.Add(1)
		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range
		if (i == thread-1) {
			max += diff
		}
		req, _ := http.NewRequest("GET", src, nil)
		var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: true,
		}
		go func(min int, max int, i int) {
			for {
				os.Remove(target + "." + strconv.Itoa(i))
				bytesrange := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
				req.Header.Set("Range", bytesrange)
				client := http.Client{
					Transport: transport,
					Timeout:   timeout * time.Second,
				}
				resp, e := client.Do(req)

				if e != nil {
					log.Println("[request error],retry:", e)
					timeout += 5
					continue
				}

				ff, _ := os.Create(target + "." + strconv.Itoa(i))
				_, e = io.Copy(ff, resp.Body)
				if e != nil {
					log.Println("[copy error],retry", e)
					ff.Close()
					resp.Body.Close()
					continue
				}
				ff.Close()
				resp.Body.Close()
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
		chunkf, _ := os.Open(target + "." + strconv.Itoa(j))
		io.Copy(f, chunkf)
		chunkf.Close()
		os.Remove(target + "." + strconv.Itoa(j))
	}
	f.Close()
	return
}
