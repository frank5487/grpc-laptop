server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
            --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
            proto/*.proto

clean:
	rm -rf pb/*.go

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen, server, client, clean, test cert