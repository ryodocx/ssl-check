package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

type SSLCheckRecord struct {
	Host      string
	Port      string
	Error     error
	IPs       []net.IP
	ExpiresOn *time.Time
	CheckTime *time.Time
}

func SSLCheck(hostPort string) (*SSLCheckRecord, error) {
	checkSet := &SSLCheckRecord{}

	// validate
	switch str := strings.Split(hostPort, ":"); len(str) {
	case 2: // host + port
		checkSet.Host = str[0]
		checkSet.Port = str[1]
	case 1: // hostname only
		checkSet.Host = str[0]
		checkSet.Port = "443"
	default:
		return nil, fmt.Errorf("invalid hostPort: %s", hostPort)
	}

	// check
	checkSet.check()

	return checkSet, nil
}

func (p *SSLCheckRecord) check() {
	defer func() {
		t := time.Now().Local()
		p.CheckTime = &t

		if p.ExpiresOn != nil {
			t := p.ExpiresOn.Local()
			p.ExpiresOn = &t
		}
	}()

	// ip lookup
	if ips, err := net.LookupIP(p.Host); err != nil {
		p.Error = fmt.Errorf("can't lookup ip: %v", err)
		return
	} else {
		p.IPs = ips
	}

	// check
	conn, err := tls.DialWithDialer(
		&net.Dialer{
			Timeout: time.Second * 5,
		},
		"tcp",
		fmt.Sprintf("%s:%s", p.Host, p.Port),
		&tls.Config{
			ServerName: p.Host,
		},
	)
	if err != nil {
		p.Error = err
		return
	}
	defer conn.Close()

	state := conn.ConnectionState()
	p.ExpiresOn = &state.PeerCertificates[0].NotAfter
}
