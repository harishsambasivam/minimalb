package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	url       *url.URL
	urlString string
	rp        *httputil.ReverseProxy
	health    bool
}

func createServer(urlString string) *Server {
	url, _ := url.Parse(urlString)
	return &Server{
		urlString: urlString,
		url:       url,
		rp:        httputil.NewSingleHostReverseProxy(url),
		health:    true,
	}
}

func (s *Server) healthCheck() {
	resp, err := http.Head(s.urlString)
	if err != nil {
		log.Println(err)
		s.health = false
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Print("server is unhealthy")
		s.health = false
		return
	}

	s.health = true

}
