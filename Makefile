all: build

run-server:
	cd cmd; go run server/main.go

run-client:
	cd cmd; go run client/main.go

build: bin
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

server:
	go build -o bin/server cmd/server/main.go
	bin/server

client:
	go build -o bin/client cmd/client/main.go
	bin/client

compile: bin
	echo "Compiling for every OS and Platform"
	echo "Compiling server..."
	GOOS=freebsd GOARCH=386 go build -o bin/server-freebsd-386 cmd/server/main.go
	GOOS=linux GOARCH=386 go build -o bin/server-linux-386 cmd/server/main.go
	GOOS=windows GOARCH=386 go build -o bin/server-windows-386 cmd/server/main.go
	echo "Compiling client..."
	GOOS=freebsd GOARCH=386 go build -o bin/client-freebsd-386 cmd/client/main.go
	GOOS=linux GOARCH=386 go build -o bin/client-linux-386 cmd/client/main.go
	GOOS=windows GOARCH=386 go build -o bin/client-windows-386 cmd/client/main.go
	echo "Done"

bin:
	mkdir -p bin

clean:
	rm -rf bin/*
