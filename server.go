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
	delay := "10"
	addrs := []request{}
	addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.177:80/cws/10.66.76.172/input/routeInputToOutput/1/0`,
		[]byte("{\"input\":\"1:0\"}"),
	})
	addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.177:80/cws/10.66.76.172/input/routeInputToOutput/1/2`,
		[]byte("{\"input\":\"1:2\"}"),
	})
	addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.177:80/cws/10.66.76.172/input/routeInputToOutput/2/1`,
		[]byte("{\"input\":\"2:1\"}"),
	})
	addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.177:80/cws/10.66.76.172/input/routeInputToOutput/1/0`,
		[]byte("{\"input\":\"1:0\"}"),
	})
	addrs = append(addrs, request{`http://10.5.34.21:8018/metered/` + delay + `/10.66.76.177:80/cws/10.66.76.172/input/routeInputToOutput/1/2`,
		[]byte("{\"input\":\"1:2\"}"),
	})

	interval := 1
	count := 100
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
