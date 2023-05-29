package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"os"

	"github.com/dexiang/url-shortener/internal/app/service"
)

type Endpoints struct {
	ShortenEndpoint  endpoint.Endpoint
	RedirectEndpoint endpoint.Endpoint
}

func MakeShortenEndpoint(s service.URLShortenerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShortenRequest)
		urlID, err := s.ShortenURL(ctx, req.URL, req.ExpireAt)
		return ShortenResponse{ID: urlID, ShortUrl: "http://" + os.Getenv("SERVICE_HOST") + ":" + os.Getenv("SERVICE_PORT") + "/" + urlID}, err
	}
}

func MakeRedirectEndpoint(s service.URLShortenerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RedirectRequest)
		originalURL, err := s.GetOriginalURL(ctx, req.ID)
		if err == service.ErrIDNotFound {
			return RedirectResponse{Res: "", Err: err}, nil
		}
		return RedirectResponse{Res: originalURL}, err
	}
}

func New(s service.URLShortenerService, logger log.Logger) Endpoints {

	var shortenEndpoint endpoint.Endpoint
	{
		shortenEndpoint = MakeShortenEndpoint(s)
		shortenEndpoint = LoggingMiddleware(log.With(logger, "method", "shorten"))(shortenEndpoint)
	}
	var redirectEndpoint endpoint.Endpoint
	{
		redirectEndpoint = MakeRedirectEndpoint(s)
		redirectEndpoint = LoggingMiddleware(log.With(logger, "method", "redirect"))(redirectEndpoint)
	}

	return Endpoints{
		ShortenEndpoint:  shortenEndpoint,
		RedirectEndpoint: redirectEndpoint,
	}
}
