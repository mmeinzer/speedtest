package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kylegrantlucas/speedtest"
)

func main() {
	endpoint := flag.String("e", "", "POST request endpoint url")
	username := flag.String("u", "", "basic auth username")
	password := flag.String("p", "", "basic auth password")
	flag.Parse()

	if *endpoint == "" {
		fmt.Printf("error endpoint url not found (-e)\n")
		os.Exit(1)
	}

	if *username == "" {
		fmt.Printf("error basic auth username not found (-u)\n")
		os.Exit(1)
	}

	if *password == "" {
		fmt.Printf("error basic auth password not found (-p)\n")
		os.Exit(1)
	}

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		fmt.Printf("error creating client: %v\n", err)
		os.Exit(1)
	}

	// Pass an empty string to select the fastest server
	server, err := client.GetServer("")
	if err != nil {
		fmt.Printf("error getting server: %v\n", err)
		os.Exit(1)
	}

	dmbps, err := client.Download(server)
	if err != nil {
		fmt.Printf("error getting download: %v\n", err)
		os.Exit(1)
	}

	umbps, err := client.Upload(server)
	if err != nil {
		fmt.Printf("error getting upload: %v\n", err)
		os.Exit(1)
	}

	ping := int(math.Round(server.Latency))

	buf := strings.NewReader(fmt.Sprintf(`{"ping":%d,"download":%3.2f,"upload":%3.2f}`, ping, dmbps, umbps))
	req, err := http.NewRequest("POST", *endpoint, buf)
	if err != nil {
		fmt.Printf("error creating POST request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(*username, *password)

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("error submitting POST request: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode == 401 {
		fmt.Println("error unauthorized - check basic auth credentials")
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("error non-200 response code from server: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Printf("Ping: %v ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", ping, dmbps, umbps)
	os.Exit(0)
}
