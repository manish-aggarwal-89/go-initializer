.PHONY : test sonar mock swagger

PACKAGES = $(shell go list ./... | grep -v -e . -e mocks | tr '\n' ',')

add-env:
	@export CONFIG_FILE = application.env

# unit test & calculate code coverage
test:
	@if [ -f coverage.out ]; then rm coverage.out; fi;
	@echo ">> running unit test and calculate coverage"
	@go test ./... -cover -coverprofile=coverage.out -covermode=count -coverpkg=$(PACKAGES)
	@go tool cover -func=coverage.out

# running test and sonar
sonar:
	@make test
	@sonar-scanner -Dproject.settings=./sonar-local.properties

mock:
	@mockery --all  --recursive=true --inpackage --case snake --with-expecter=true

# this will use .mockery.yml template, make sure your mockery version is 3.
mockV3:
	@mockery

swagger:
	@swag init -g ./cmd/{{PROVIDER}}/main.go --parseInternal --parseDependency
	@swag fmt -g ./cmd/{{PROVIDER}}/main.go

build:
	@go build -ldflags "-s -w" -o {{SERVICE_NAME}} ./cmd/{{PROVIDER}}
	
run:
	@go run cmd/{{PROVIDER}}/main.go