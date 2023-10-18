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
		mDNS      *mdns.Conn
		dns       *dns.Server
		Log       logr.Logger
		Timeout   time.Duration
		IP        string
		Port      int
		Zone      string
		Recusrive bool
		client    *dns.Client
		Upstream  string
	}
)

// ListenAndServe Start the mDNS and DNS systems
func (srv *Server) ListenAndServe() error {
	addr := netip.MustParseAddrPort(mdns.DefaultAddress)
	udpAddr := net.UDPAddrFromAddrPort(addr)

	srv.Log.Info("starting mDNS listener ", "mDNS Listener", addr)
	l, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}

	mDNSSrv, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		QueryInterval: srv.Timeout,
	})
	if err != nil {
		return fmt.Errorf("mdns server error: %w", err)
	}

	if srv.Port == 0 {
		srv.Port = port
	}

	if srv.IP == "" {
		srv.IP = ip
	}

	srv.dns = &dns.Server{
		Net:  "udp",
		Addr: fmt.Sprintf("%s:%d", srv.IP, srv.Port),
		//	IdleTimeout: func() time.Duration { return s.Timeout },
	}

	srv.mDNS = mDNSSrv

	if srv.Zone == "" {
		srv.Zone = zone
	}

	srv.Log.Info("starting DNS listener",
		"addr", srv.dns.Addr, "zone", srv.Zone)

	mux := dns.NewServeMux()
	mux.HandleFunc(srv.Zone, srv.dnsHandler)

	if srv.Recusrive {
		srv.client = &dns.Client{}
		mux.HandleFunc(".", srv.recursiveHandler)
	}

	srv.dns.Handler = mux

	err = srv.dns.ListenAndServe()
	if err != nil {
		return fmt.Errorf("dns server error: %w", err)
	}

	return nil
}

func (srv *Server) Close() error {
	if err := srv.mDNS.Close(); err != nil {
		return fmt.Errorf("mdns agent shutdown failed: %w", err)
	}

	if err := srv.dns.Shutdown(); err != nil {
		return fmt.Errorf("dns server shutdown failed: %w", err)
	}

	return nil
}

func (srv *Server) recursiveHandler(w dns.ResponseWriter, req *dns.Msg) {
	for _, host := range req.Question {
		srv.Log.Info("querying for dns", "host", host.Name)

	}
	m, _, err := srv.client.Exchange(req, srv.Upstream)

	if err != nil {
		srv.Log.Error(err, "host exchange")
	}
	w.WriteMsg(m)

}

func (srv *Server) dnsHandler(w dns.ResponseWriter, req *dns.Msg) {
	var resp dns.Msg

	resp.SetReply(req)

	for _, host := range req.Question {
		srv.Log.Info("querying for mdns", "host", host.Name)

		ctx, can := context.WithTimeout(context.Background(), srv.Timeout)

		defer can()

		newHost := srv.rewriteHostname(host.Name)

		// TODO need to pull out context errors here if we can
		// as that means we timed out resolving hostname
		// not everything is on fire.
		// mdns seems to wrap this it might be tricky.
		ip, err := srv.mDNSResolveHostname(ctx, newHost)
		if err != nil {
			srv.Log.Error(err, "mdns resolve hostname", "host", newHost)
			continue
		}

		srv.Log.Info("found ip", "ip", ip)

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
		srv.Log.Info("no records found")
	}

	if err := w.WriteMsg(&resp); err != nil {
		srv.Log.Error(err, "could not write packet", "resp", resp)
	}
}

func (srv *Server) rewriteHostname(host string) string {
	h := host

	if srv.Zone == zone {
		return h
	}

	i := strings.LastIndex(h, srv.Zone)
	if i == -1 {
		return h
	}

	h = fmt.Sprintf("%s%s", host[0:i], zone)

	return h
}

func (srv *Server) mDNSResolveHostname(ctx context.Context, hostname string) (net.IP, error) {
	h := strings.TrimSuffix(hostname, ".")

	_, addr, err := srv.mDNS.Query(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dns hostname: %w", err)
	}

	return net.ParseIP(addr.String()), nil
}
