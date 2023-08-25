gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs/logs.proto
run:
	cd ./cmd/web && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go run main.go