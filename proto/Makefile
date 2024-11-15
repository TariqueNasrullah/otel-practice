book:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go_opt=paths=source_relative book.proto
all_with_docker:
	docker run -it --rm -v $(CURDIR):/proto/src tarique/buffer
all:
	make book