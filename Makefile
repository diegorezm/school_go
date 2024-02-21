run:
	@go run main.go

build: 
	@rm -rf dist
	@mkdir dist
	@go build -o dist/main.bin main.go
