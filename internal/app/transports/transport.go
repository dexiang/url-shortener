package transports

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/dexiang/url-shortener/internal/app/endpoints"
)

func decodeShortenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeShortenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeRedirectRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoints.RedirectRequest{ID: path.Base(r.URL.RequestURI())}, nil
}

func encodeRedirectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	resp := response.(endpoints.RedirectResponse)

	if resp.Err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Location", resp.Res)
		w.WriteHeader(http.StatusFound)
	}

	return nil
}

func NewHTTPHandler(ctx context.Context, endpoints endpoints.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/api/v1/urls").Handler(httptransport.NewServer(
		endpoints.ShortenEndpoint,
		decodeShortenRequest,
		encodeShortenResponse,
	))

	// {id:[a-zA-Z0-9]+}
	// {id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}
	r.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.RedirectEndpoint,
		decodeRedirectRequest,
		encodeRedirectResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
