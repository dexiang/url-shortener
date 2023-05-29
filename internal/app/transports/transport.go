package transports

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/dexiang/url-shortener/internal/app/endpoints"
)

func decodeShortenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ShortenRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	validate := validator.New()
	if e := validate.Struct(req); e != nil {
		return nil, e
	}

	return req, nil
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

func NewHTTPHandler(ctx context.Context, endpoints endpoints.Endpoints, logger log.Logger) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/api/v1/urls").Handler(httptransport.NewServer(
		endpoints.ShortenEndpoint,
		decodeShortenRequest,
		encodeShortenResponse,
		options...,
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

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if errors, ok := err.(validator.ValidationErrors); ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errors.Error(),
		})
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
}
