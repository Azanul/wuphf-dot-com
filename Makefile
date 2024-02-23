init:
	sudo apt-get update
	sudo apt install protobuf-compiler
	sudo apt install golang-goprotobuf-dev
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
proto-gen:
	protoc --go_out=user --go-grpc_out=user user/api/auth.proto