package endpoints

import "time"

type ShortenRequest struct {
	URL      string    `json:"url" validate:"required,url"`
	ExpireAt time.Time `json:"expireAt"` // ISO8601
}

type RedirectRequest struct {
	ID string `json:"id"`
}
