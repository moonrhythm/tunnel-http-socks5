package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/acoshift/configfile"
	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/logger"
	"github.com/moonrhythm/parapet/pkg/upstream"
)

var config = configfile.NewEnvReader()

var (
	port                 = config.IntDefault("port", 8080)
	upstreamAddr         = config.String("upstream_addr")  // comma split addr
	upstreamProto        = config.String("upstream_proto") // http, h2c, https, unix
	upstreamHeaderSet    = parseHeaders(config.String("upstream_header_set"))
	upstreamHeaderAdd    = parseHeaders(config.String("upstream_header_add"))
	upstreamHeaderDel    = parseHeaders(config.String("upstream_header_del"))
	upstreamOverrideHost = config.String("upstream_override_host")
	upstreamPath         = config.String("upstream_path") // prefix path
)

func main() {
	fmt.Println("tunnel-http-socks")
	fmt.Println()

	s := parapet.Server{}
	s.Use(logger.Stdout())

	var targets []*upstream.Target
	for _, addr := range strings.Split(upstreamAddr, ",") {
		targets = append(targets, &upstream.Target{
			Host:      addr,
			Transport: &upstream.HTTPTransport{},
		})
	}

	us := upstream.New(upstream.NewRoundRobinLoadBalancer(targets))
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

func parseHeaders(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	ss := strings.Split(s, ",")

	var rs []string
	for _, x := range ss {
		ps := strings.Split(x, ":")
		if len(ps) != 2 {
			continue
		}
		rs = append(rs, strings.TrimSpace(ps[0]), strings.TrimSpace(ps[1]))
	}

	return rs
}
