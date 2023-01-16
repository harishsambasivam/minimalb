package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	servers = []*Server{
		createServer("http://localhost:5001"),
		createServer("http://localhost:5002"),
	}
	lastServed = 0
)

func main() {
	http.HandleFunc("/", forwardRequest)
	go healthCheck()
	fmt.Println("staring load balancer on " + GetLocalIP() + ":" + os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	reverseProxy, err := getHealthyReverseProxy()
	if err != nil {
		io.WriteString(res, "No healthy hosts found")
	}
	reverseProxy.ServeHTTP(res, req)
}

func getHealthyReverseProxy() (*httputil.ReverseProxy, error) {
	for i := 0; i < len(servers); i++ {
		server := servers[lastServed%len(servers)]
		lastServed += 1
		if server.health {
			return server.rp, nil
		} else {
			log.Print("unhealthy server found")
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}
