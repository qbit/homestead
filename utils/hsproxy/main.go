package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	to := flag.String("to", "", "Where to proxy to")
	listen := flag.String("http", ":8080", "what to listen on")

	flag.Parse()

	u, err := url.Parse(*to)
	if err != nil {
		log.Fatal(err)
	}

	dir := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			fmt.Printf("'%s://%s' => '%s://%s'", req.URL.Scheme, req.URL.Host, u.Scheme, u.Host)
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
		},
	}

	http.Handle("/", dir)

	if err = http.ListenAndServe(*listen, nil); err != nil {
		log.Fatal(err)
	}
}
