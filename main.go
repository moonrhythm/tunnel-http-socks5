package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/acoshift/configfile"
	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/logger"
	"github.com/moonrhythm/parapet/pkg/upstream"
	"golang.org/x/crypto/ssh"
)

var config = configfile.NewEnvReader()

var (
	port                 = config.IntDefault("port", 8080)
	upstreamAddr         = config.String("upstream_addr")
	upstreamOverrideHost = config.String("upstream_override_host")
	upstreamPath         = config.String("upstream_path") // prefix path
	tunnelAddr           = config.String("tunnel_addr")
	tunnelUser           = config.String("tunnel_user")
	tunnelSSHKEY         = config.Base64("tunnel_ssh_key")
)

func main() {
	fmt.Println("tunnel-http-socks5")
	fmt.Println()

	s := parapet.Server{}
	s.Use(logger.Stdout())

	priKey, err := ssh.ParsePrivateKey(tunnelSSHKEY)
	if err != nil {
		log.Fatalf("can not parse private key; %v", err)
		return
	}

	var (
		muClient  sync.RWMutex
		sshClient *ssh.Client
	)
	connectSSH := func() error {
		muClient.Lock()
		defer muClient.Unlock()

		var err error
		sshClient, err = ssh.Dial("tcp", tunnelAddr, &ssh.ClientConfig{
			User:            tunnelUser,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(priKey)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		return err
	}
	err = connectSSH()
	if err != nil {
		log.Printf("ssh: dial error; %v", err)
		return
	}
	defer sshClient.Close()

	us := upstream.New(upstream.SingleHost(upstreamAddr, &transport{&http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			muClient.RLock()
			defer muClient.RUnlock()

			conn, err = sshClient.Dial(network, addr)
			if err != nil {
				// re-connect
				muClient.RUnlock()
				connectSSH()
				muClient.RLock()

				conn, err = sshClient.Dial(network, addr)
			}
			return
		},
	}}))
	us.Host = upstreamOverrideHost
	us.Path = upstreamPath
	s.Use(us)

	fmt.Println("Upstream", upstreamAddr)

	s.Addr = fmt.Sprintf(":%d", port)
	fmt.Println("Starting parapet on port", port)
	fmt.Println()

	err = s.ListenAndServe()
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
