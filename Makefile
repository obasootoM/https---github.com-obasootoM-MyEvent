certificate:
	go run $GOROOT/usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
	
main:
	go run main.go

.PHONY:certificate main