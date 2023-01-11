// Package server implments the mdns-proxy server.
// This server implements a DNS listner that proxies all requests
// to the local mdns service.
package proxy

import (
	"testing"
)

func TestServer_rewriteHostname(t *testing.T) {
	type fields struct {
		Zone string
	}
	type args struct {
		host string
	}
	tests := []struct {
		name string
		Zone string
		host string
		want string
	}{
		{
			name: "default local",
			Zone: "local.",
			host: "foo.local.",
			want: "foo.local.",
		},
		{
			name: "rewrite the domain",
			Zone: "dev.",
			host: "foo.dev.",
			want: "foo.local.",
		},
		{
			name: "duplicate domains",
			Zone: "dev.",
			host: "foo.dev.dev.",
			want: "foo.dev.local.",
		},
		{
			name: "duplicate domains",
			Zone: "dev.",
			host: "foo.bar.dev.",
			want: "foo.bar.local.",
		},
		{
			name: "nested domains",
			Zone: "foo.bar.",
			host: "blarg.foo.bar.",
			want: "blarg.local.",
		},
		{
			name: "not in domain",
			Zone: "foo.bar.",
			host: "beep.boop.bop.",
			want: "beep.boop.bop.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Zone: tt.Zone,
			}
			if got := s.rewriteHostname(tt.host); got != tt.want {
				t.Errorf("Server.rewriteHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}
