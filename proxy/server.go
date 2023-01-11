// Package server implments the mdns-proxy server.
// This server implements a DNS listner that proxies all requests
// to the local mdns service.
package proxy

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/miekg/dns"
	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
)

const (
	ip   = ""
	port = 53
	zone = "local."
)

type (
	// Server is the mDNS proxy server
	Server struct {
		mDNS    *mdns.Conn
		dns     *dns.Server
		Log     logr.Logger
		Timeout time.Duration
		IP      string
		Port    int
		Zone    string
	}
)

func (s *Server) ListenAndServe() error {
	addr := netip.MustParseAddrPort(mdns.DefaultAddress)
	udpAddr := net.UDPAddrFromAddrPort(addr)

	s.Log.Info("starting mDNS listener ", "mDNS Listener", addr)
	l, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}

	mDNSSrv, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		QueryInterval: s.Timeout,
	})
	if err != nil {
		return fmt.Errorf("mdns server error: %w", err)
	}

	if s.Port == 0 {
		s.Port = port
	}

	if s.IP == "" {
		s.IP = ip
	}

	s.dns = &dns.Server{
		Net:  "udp",
		Addr: fmt.Sprintf("%s:%d", s.IP, s.Port),
		//	IdleTimeout: func() time.Duration { return s.Timeout },
	}

	s.mDNS = mDNSSrv

	if s.Zone == "" {
		s.Zone = zone
	}

	s.Log.Info("starting DNS listener", "addr", s.dns.Addr, "zone", s.Zone)
	mux := dns.NewServeMux()
	mux.HandleFunc(s.Zone, s.dnsHandler)
	s.dns.Handler = mux

	err = s.dns.ListenAndServe()
	if err != nil {
		return fmt.Errorf("dns server error: %w", err)
	}

	return nil
}

func (s *Server) Close() error {
	if err := s.mDNS.Close(); err != nil {
		return fmt.Errorf("mdns agent shutdown failed: %w", err)
	}

	if err := s.dns.Shutdown(); err != nil {
		return fmt.Errorf("dns server shutdown failed: %w", err)
	}

	return nil
}

func (s *Server) dnsHandler(w dns.ResponseWriter, req *dns.Msg) {
	var resp dns.Msg

	resp.SetReply(req)

	for _, host := range req.Question {
		s.Log.Info("querying for host", "host", host.Name)

		ctx, can := context.WithTimeout(context.Background(), s.Timeout)

		defer can()

		newHost := s.rewriteHostname(host.Name)

		// TODO need to pull out context errors here if we can
		// as that means we timed out resolving hostname
		// not everything is on fire.
		ip, err := s.mDNSResolveHostname(ctx, newHost)
		if err != nil {
			s.Log.Error(err, "mdns resolve hostname", "host", newHost)
			continue
		}

		s.Log.Info("found ip", "ip", ip)

		rec := dns.A{
			Hdr: dns.RR_Header{
				Name:   host.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: ip,
		}

		resp.Answer = append(resp.Answer, &rec)
	}

	if len(resp.Answer) == 0 {
		s.Log.Info("no records found")
	}

	if err := w.WriteMsg(&resp); err != nil {
		s.Log.Error(err, "could not write packet", "resp", resp)
	}
}

func (s *Server) rewriteHostname(host string) string {
	h := host

	if s.Zone == zone {
		return h
	}

	i := strings.LastIndex(h, s.Zone)
	if i == -1 {
		return h
	}

	h = fmt.Sprintf("%s%s", host[0:i], zone)

	return h
}

func (s *Server) mDNSResolveHostname(ctx context.Context, hostname string) (net.IP, error) {
	h := strings.TrimSuffix(hostname, ".")

	_, addr, err := s.mDNS.Query(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dns hostname: %w", err)
	}

	return net.ParseIP(addr.String()), nil
}
