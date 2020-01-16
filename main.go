package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/acoshift/configfile"
	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/logger"
	"github.com/moonrhythm/parapet/pkg/upstream"
)

var config = configfile.NewEnvReader()

var (
	port                 = config.IntDefault("port", 8080)
	upstreamAddr         = config.String("upstream_addr")
	upstreamOverrideHost = config.String("upstream_override_host")
	upstreamPath         = config.String("upstream_path") // prefix path
)

func main() {
	fmt.Println("tunnel-http-socks5")
	fmt.Println()

	s := parapet.Server{}
	s.Use(logger.Stdout())

	socks5URL, _ := url.Parse("socks5://127.0.0.1:5000")

	us := upstream.New(upstream.SingleHost(upstreamAddr, &transport{&http.Transport{
		Proxy: http.ProxyURL(socks5URL),
	}}))
	us.Host = upstreamOverrideHost
	us.Path = upstreamPath
	s.Use(us)

	fmt.Println("Upstream", upstreamAddr)

	s.Addr = fmt.Sprintf(":%d", port)
	fmt.Println("Starting parapet on port", port)
	fmt.Println()

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type transport struct {
	h *http.Transport
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	return t.h.RoundTrip(r)
}
