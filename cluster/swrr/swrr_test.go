package swrr

import (
	"log"
	"strconv"
	"testing"

	"local/global"
)

func TestNext(t *testing.T) {
	var (
		obj    Cluster
		ip     string
		port   uint16
		weight int
	)
	for i := 1; i <= 10; i++ {
		if i == 6 {
			ip = "10.0.0.6"
			port = 80
			weight = 6
		} else {
			ip = "10.0.0." + strconv.Itoa(i)
			port = 80
			weight = 1
		}
		obj.Set(global.Node{
			IP:     ip,
			Port:   port,
			Weight: weight,
		})
		log.Println(ip, port, weight)
	}

	for i := 0; i < 10; i++ {
		t.Log(obj.Select())
	}
}
