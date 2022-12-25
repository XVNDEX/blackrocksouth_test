gen:
	protoc --go_out=pb --go-grpc_out=pb proto/*.proto

clean:
	rm  pb/*

auth-server:
	go run cmd/server/auth.go

catalog-server:
	go run cmd/server/catalog.go

test:
	go run cmd/client/test-client.go