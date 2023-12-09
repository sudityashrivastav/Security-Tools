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
	DomainMutex         sync.Mutex
	DomainWG            sync.WaitGroup
	DomainSemaphore     = make(chan struct{}, 100)
	availableSubDomains []string
)

func HandleDomain(URL string, fileLocation string) ([]string, error) {

	file, err := os.Open(fileLocation)
	if err != nil {
		fmt.Println(err)
		return availableSubDomains, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()

		DomainWG.Add(1)
		go checkDomain(URL, word)

	}

	DomainWG.Wait()

	return availableSubDomains, err
}

func checkDomain(domain string, word string) {

	defer DomainWG.Done()

	DomainSemaphore <- struct{}{}

	defer func() { <-DomainSemaphore }()

	mainDomain := strings.TrimPrefix(domain, "https://")

	subDomain := "https://" + word + "."

	url := subDomain + mainDomain

	client := &http.Client{Timeout: 5 * time.Second}

	res, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	fmt.Println("Checking: ", url)

	if res.StatusCode != 404 || res.StatusCode != 406 {
		DomainMutex.Lock()
		availableSubDomains = append(availableSubDomains, url)
		DomainMutex.Unlock()
	}
}
