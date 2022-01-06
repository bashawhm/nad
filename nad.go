package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

type Clients struct {
	lock      sync.Mutex
	endpoints []mdns.ServiceEntry
}

var clients Clients

func main() {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	info := []string{"NAD File Drop Service"}
	service, err := mdns.NewMDNSService(host, "_nad._tcp", "", host, 9000, nil, info)
	if err != nil {
		panic(err)
	}

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	mdnsEntries := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for {
			for entry := range mdnsEntries {
				clients.lock.Lock()

				found := false
				for i := 0; i < len(clients.endpoints); i++ {
					if entry.Name == clients.endpoints.Name {
						found = true
					}
				}
				if !found {
					clients.endpoints = append(clients.endpoints, entry)
					fmt.Println("entry.Name =", entry.Name)
				}

				clients.lock.Unlock()
			}
			time.Sleep(1)
		}
	}()

	mdns.Lookup("_nad._tcp", mdnsEntries)

	for {
	}
}
