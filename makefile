build:
	@go build -o librenotes ./cmd/librenotes/
run: build
	@./librenotes
