all: clean build 
build: build_css build_go
build_go:
	go build -o bin/hummingbard cmd/hummingbard/main.go
build_css:
	./build.sh
vendor: clean vendorbuild 
vendorbuild:
	go build -mod=vendor -o bin/hummingbard cmd/hummingbard/main.go
clean: 
	rm -f bin/hummingbard
	rm -rf static/css
	mkdir static/css
