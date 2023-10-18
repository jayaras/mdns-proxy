module mdns-proxy

go 1.19

//replace github.com/pion/mdns => github.com/jaras/mdns [master]

require (
	github.com/go-logr/logr v1.2.3
	github.com/go-logr/zapr v1.2.3
	github.com/miekg/dns v1.1.50
	github.com/pion/mdns v0.0.6-0.20230110082909-9dd554ad1ce5
	github.com/spf13/cobra v1.7.0
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.4.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
)
