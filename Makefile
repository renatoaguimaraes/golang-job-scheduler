test:
	go test -race -v -timeout 30s -failfast -cover github.com/renatoaguimaraes/job-scheduler/...

api:
	go build -o ./bin/worker-api cmd/api/main.go
