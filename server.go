package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

type request struct {
	Address          string
	ExpectedResponse []byte
}

func main() {
	newTest()
}

func newTest() {
	addrs := []request{}
	addrs = append(addrs, request{`http://10.66.9.7:8017/10.66.9.200/input/1/0`,
		[]byte("{\r\n  \"input\": \"1:0\"\r\n}"),
	})
	addrs = append(addrs, request{`http://10.66.9.7:8017/10.66.9.200/input/1/2`,
		[]byte("{\r\n  \"input\": \"1:2\"\r\n}"),
	})
	addrs = append(addrs, request{`http://10.66.9.7:8017/10.66.9.200/input/3/1`,
		[]byte("{\r\n  \"input\": \"3:1\"\r\n}"),
	})
	addrs = append(addrs, request{`http://10.66.9.7:8017/10.66.9.200/input/3/0`,
		[]byte("{\r\n  \"input\": \"3:0\"\r\n}"),
	})
	addrs = append(addrs, request{`http://10.66.9.7:8017/10.66.9.200/input/3/2`,
		[]byte("{\r\n  \"input\": \"3:2\"\r\n}"),
	})

	interval := 1000
	count := 10
	concurrent := 8
	wg := sync.WaitGroup{}

	wg.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go newMakeRequests(addrs, interval, count, i, &wg)
	}
	wg.Wait()

	log.Printf("done")

}

func newMakeRequests(addr []request, interval int, count int, id int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		time.Sleep(time.Duration(interval) * time.Millisecond)

		cur := addr[i%len(addr)]
		resp, err := http.Get(cur.Address)
		if err != nil {
			log.Printf(color.HiRedString("[%v] Error on request: %v", id, err.Error()))
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf(color.HiRedString("[%v] Error reading body: %v", id, err.Error()))
			continue
		}

		if resp.StatusCode != 200 {
			log.Printf(color.HiRedString("[%v] Error on request: Non 200 %s", id, body))
			continue
		}

		//check the response
		if !compareSlices(body, cur.ExpectedResponse) {
			log.Printf(color.HiRedString("[%v] Incorrect response %v", id, body))
			log.Printf(color.HiRedString("[%v]  Expected response %v", id, cur.ExpectedResponse))
			continue
		}

		log.Printf(color.HiGreenString("[%v] ok.", id))
		resp.Body.Close()
	}
	log.Printf(color.HiGreenString("[%v] Done.", id))
	wg.Done()
}

func compareSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
