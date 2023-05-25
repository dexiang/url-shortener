# URL Shortener

## Usage

### Generate Shortened URL
```zsh
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{"url": "<original_url>","expireAt": "2021-02-08T09:20:41Z"}'
```
Response
```json
{
    "id": "<url_id>",
    "shortUrl": "http://localhost/<url_id>"
}
```

### Redirect URL API
```zsh
curl -L -X GET http://localhost/<url_id>
```
use 302 Redirect to original URL

## Project Layout
- single service put at `internal/app`
- `cmd` is service entry point
- `internal/pkg/<shared-pkg>` is shared pkg for other service

```
.
├── README.md
├── cmd
│   └── main.go
├── internal
│   ├── app
│   │   ├── endpoint
│   │   ├── service
│   │   └── transport
│   └── pkg 
├── go.mod
└── go.sum
```
