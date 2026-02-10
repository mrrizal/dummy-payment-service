run:
	go run cmd/api/main.go

build:
	go build -o dummy-payment-service cmd/api/main.go

run-build:
	./dummy-payment-service
