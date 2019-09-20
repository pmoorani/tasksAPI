.PHONY: all # All targets are accessible for user
.DEFAULT: help # Running Make will run the help target


TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

TMS_ENV=development
export TMS_ENV

GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/tasksAPI


$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web

test:
	go test ./...

build:
	go build -ldflags "-X main.version=latest" -o tms .

pack: #build
	docker build -t pmoorani/tms-service:latest -f BinaryDockerfile .

upload:
	docker push pmoorani/tms-service:latest

deploy:
	kubectl apply -f k8s/deployment.yaml

ship: test pack upload deploy

dev:
	go run .
