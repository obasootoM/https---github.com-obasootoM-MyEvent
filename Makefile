certificate:
	go run $GOROOT/usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
	
main:
	go run main.go

docker:
	sudo docker run -d --name rabbitmq -h rabbit-mq -p 8080:15672 rabbitmq:3-management
	
run:
	go test -v ./...

event:
	sudo docker run -d --name event-db --network myevents

booking:
	sudo docker run -d --name booking-db --network myevents	

dockerRun:
	sudo docker run \
	--detach \
	--name events \
	--network myevents \
	-e AMQPMESSAGEBROKER=amqp://guest:guest@rabbit-mq:5672/ \
	-e MONGO_URL=mongodb://event-db/events \
	-p 9191:9191 \
	myevents/eventServe


.PHONY:certificate main docker run event booking dockerRun