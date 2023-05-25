package main

import (
	"context"
	"net/http"

	"github.com/dexiang/url-shortener/internal/endpoints"
	"github.com/dexiang/url-shortener/internal/service"
	"github.com/dexiang/url-shortener/internal/transports"
)

func main() {

	ctx := context.Background()
	service := service.New()
	endpoints := endpoints.New(service)
	handler := transports.NewHTTPServer(ctx, endpoints)

	http.ListenAndServe(":80", handler)
}
