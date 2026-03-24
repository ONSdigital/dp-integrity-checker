module github.com/ONSdigital/dp-integrity-checker

go 1.26.0

require (
	github.com/ONSdigital/log.go/v2 v2.5.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/slack-go/slack v0.20.0
	github.com/smartystreets/goconvey v1.8.1
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.277.0 // indirect
	github.com/ONSdigital/dp-net/v3 v3.10.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.42.0 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

// [CVE-2022-27191] CWE-noinfo: Allows an attacker to crash a server in certain circumstances involving AddHostKey
tool golang.org/x/crypto
