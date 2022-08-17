certificate:
	go run $GOROOT/usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
	
main:
	go run main.go

docker:
	sudo docker run -d --name rabbitmq -h rabbit-mq -p 8000:5672 -p 8080:15672 rabbitmq:3-management
run:
	go test -v ./...

.PHONY:certificate main docker run