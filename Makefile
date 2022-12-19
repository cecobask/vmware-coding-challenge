USERNAME := baskski
TAG := latest

run:
	go run cmd/apiserver/main.go

test:
	go test ./...

deploy:
	kubectl apply -f deployment.yaml

docker-build:
	docker build -t $(USERNAME)/apiserver:$(TAG) .

docker-push:
	docker image push $(USERNAME)/apiserver:$(TAG)

docker-build-push: docker-build docker-push

docker-run:
	docker run -it -p 8080:8080 --name apiserver $(USERNAME)/apiserver:$(TAG)

docker-cleanup:
	docker stop apiserver
	docker rm apiserver

kubernetes-deploy:
	kubectl apply -f deployment.yml
