package endpoints

type ShortenResponse struct {
	ID       string `json:"id"`
	ShortUrl string `json:"shortUrl"`
	Err      error  `json:"-"`
}

type RedirectResponse struct {
	Res string `json:"res"` // res
	Err error  `json:"-"`   // err
}
