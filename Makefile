run:
	@go run
build:
	@go build
build-opt:
	@go build -ldflags "-s -w"  -o csv2parquet