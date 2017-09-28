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

func Download(src, target string, thread int) {
	res, _ := http.Head(src)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0]) // Get the content length from the header request
	len_sub := length / thread // Bytes for each Go-routine
	diff := length % thread // Get the remaining for the last request
	body := make([][]byte, thread) // Make up a temporary array to hold the data to be written to the file
	for i := 0; i < thread; i++ {
		wg.Add(1)
		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range
		if (i == thread-1) {
			max += diff // Add the remaining bytes in the last request
		}
		go func(min int, max int, i int) {
			for {
				req := fasthttp.AcquireRequest()
				req.SetRequestURI(src)
				req.Header.SetByteRange(min, max-1)
				resp := fasthttp.AcquireResponse()
				client := &fasthttp.Client{}
				e := client.DoTimeout(req, resp, time.Second*30)
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
}
