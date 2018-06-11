package main

import (
	"io/ioutil"
	"log"
	"math/rand"
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
	/*
		delay := "10"
		addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.176:80/cws/10.66.76.171/input/routeInputToOutput/1/0`,
			[]byte("{\"input\":\"1:0\"}"),
		})
		addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.176:80/cws/10.66.76.171/input/routeInputToOutput/1/2`,
			[]byte("{\"input\":\"1:2\"}"),
		})
		addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.176:80/cws/10.66.76.171/input/routeInputToOutput/0/1`,
			[]byte("{\"input\":\"0:1\"}"),
		})
		addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.176:80/cws/10.66.76.171/input/routeInputToOutput/1/0`,
			[]byte("{\"input\":\"1:0\"}"),
		})
		addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.176:80/cws/10.66.76.171/input/routeInputToOutput/1/2`,
			[]byte("{\"input\":\"1:2\"}"),
		})
	*/
	addrs = append(addrs, request{`http://localhost:8005/10.5.34.48/input/current`,
		[]byte("{\"input\":\"digital3\"}"),
	})
	addrs = append(addrs, request{`http://localhost:8005/10.5.34.48/power/status`,
		[]byte("{\"power\":\"standby\"}"),
	})

	interval := 500
	count := 1000000000000
	concurrent := 3
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

		cur := addr[rand.Intn(len(addr))]
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
			log.Printf(color.HiRedString("[%v] Incorrect response %s", id, body))
			log.Printf(color.HiRedString("[%v]  Expected response %s", id, cur.ExpectedResponse))
			continue
		}

		log.Printf(color.HiGreenString("[%v] ok - resp: %s", id, body))
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
