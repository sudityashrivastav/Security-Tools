package handlers

import (
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	portWG        sync.WaitGroup
	PortSemaphore = make(chan struct{}, 1000)
	OpenPorts     []int
	PortsMutex    sync.Mutex
)

func ScanPortsMain(host string) ([]int, error) {

	for port := 1; port <= 65535; port++ {
		portWG.Add(1)
		go scanPorts(host, port)
	}

	portWG.Wait()

	return OpenPorts, nil
}

func scanPorts(host string, port int) {
	defer portWG.Done()

	semaphore <- struct{}{}

	defer func() { <-semaphore }()

	address := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		// fmt.Println("Error at dialtimout: ", err)
		return

	}

	PortsMutex.Lock()
	OpenPorts = append(OpenPorts, port)
	PortsMutex.Unlock()

	conn.Close()

}
