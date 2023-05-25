package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/dexiang/url-shortener/internal/service"
)

type Endpoints struct {
	ShortenEndpoint  endpoint.Endpoint
	RedirectEndpoint endpoint.Endpoint
}

func (e Endpoints) Shorten(ctx context.Context, url string, expireAt time.Time) (string, error) {
	resp, err := e.ShortenEndpoint(ctx, ShortenRequest{URL: url, ExpireAt: expireAt})
	if err != nil {
		return "", err
	}

	response := resp.(ShortenResponse)
	return response.Res, nil
}

func (e Endpoints) Redirect(ctx context.Context, id string) (string, error) {
	resp, err := e.RedirectEndpoint(ctx, RedirectRequest{ID: id})
	if err != nil {
		return "", err
	}

	response := resp.(RedirectResponse)
	return response.Res, nil
}

func MakeShortenEndpoint(s service.URLShortenerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShortenRequest)
		tinyURL, err := s.ShortenURL(ctx, req.URL, req.ExpireAt)
		return ShortenResponse{Res: tinyURL}, err
	}
}

func MakeRedirectEndpoint(s service.URLShortenerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RedirectRequest)
		originalURL, err := s.GetOriginalURL(ctx, req.ID)
		return RedirectResponse{Res: originalURL}, err
	}
}

func New(s service.URLShortenerService) Endpoints {
	endpoints := Endpoints{
		ShortenEndpoint:  MakeShortenEndpoint(s),
		RedirectEndpoint: MakeRedirectEndpoint(s),
	}

	return endpoints
}
