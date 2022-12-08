package gask

//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.4 fmt
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.4 init --parseDependency --parseInternal -g ./cmd/gask/main.go
