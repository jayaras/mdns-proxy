module mdns-proxy

go 1.19

//replace github.com/pion/mdns => github.com/jaras/mdns [master]

require (
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/zapr v1.2.4
	github.com/miekg/dns v1.1.56
	github.com/pion/mdns v0.0.9
	github.com/spf13/cobra v1.7.0
	go.uber.org/zap v1.26.0
	golang.org/x/net v0.17.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
)
