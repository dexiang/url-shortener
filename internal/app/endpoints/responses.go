package endpoints

type ShortenResponse struct {
	Res string `json:"res"` // res
	Err error  `json:"-"`   // err
}

type RedirectResponse struct {
	Res string `json:"res"` // res
	Err error  `json:"-"`   // err
}
