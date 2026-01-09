build:
	@go build -o librenotes .
run: build
	@./librenotes
