build:
	@go build -v -o ~/tmp/dumpsc \
	-ldflags "-s -w -X github.com/lowk3v/dumpsc/pkg/version.VERSION=$(cat VERSION)" \
	main.go;
