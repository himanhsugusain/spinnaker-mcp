package main

import (
	"context"
	"log"
	"net/http"

	server "github.com/himanhsugusain/go-mcp"
	"go.uber.org/zap"
	"spinnaker"
)

func main() {
	ctx := context.Background()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	backend, err := spinnaker.NewSpinnakerClient(ctx, logger)
	if err != nil {
		panic(err)
	}
	mcpHandler := server.NewApp(backend, logger)
	mux := http.DefaultServeMux
	mux.Handle("/mcp", mcpHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
