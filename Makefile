run:
	go run -tags=jsoniter ./cmd/gask server --debug

generate:
	go generate ./...
