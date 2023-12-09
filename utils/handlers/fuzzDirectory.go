package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mutex           sync.Mutex
	wg              sync.WaitGroup
	semaphore       = make(chan struct{}, 100)
	availableRoutes []string
)

func HandleURL(URL string, fileLocation string) []string {

	file, err := os.Open(fileLocation)
	if err != nil {
		fmt.Println("error at os.open: ", err)
		return availableRoutes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		wg.Add(1)
		go CheckURL(URL, word)
	}

	wg.Wait()

	return availableRoutes
}

func CheckURL(url string, word string) {
	defer wg.Done()

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	if !strings.HasSuffix(url, "/") {
		url = strings.TrimSpace(url) + "/"
	}

	word = strings.TrimPrefix(word, "/")

	address := url + word

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(address)

	if err != nil {
		fmt.Println("error at http get", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 404 && res.StatusCode != 406 {
		mutex.Lock()
		availableRoutes = append(availableRoutes, address)
		mutex.Unlock()
	}
}
