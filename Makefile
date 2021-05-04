test:
	go test -race -v -timeout 30s -failfast -cover github.com/renatoaguimaraes/job-scheduler/...

api:
	go build -o ./bin/worker-api cmd/api/main.go

proto:
	protoc --go_out=internal/worker --go_opt=paths=source_relative \
		--go-grpc_out=internal/worker --go-grpc_opt=paths=source_relative \
		proto/worker.proto
