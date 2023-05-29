# URL Shortener

## Usage

## start service(Local)
```zsh
skaffold dev
```

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

# Features
- [Go-Kit](https://github.com/go-kit/kit)
- DevOps: The app is designed to run on Kubernetes (both locally on "Docker for Desktop", as well as on the cloud with GKE).
  - [K8s](https://kubernetes.io/)
  - [Skaffold](https://skaffold.dev/)
  - [Buildpacks](https://buildpacks.io/)
- Redis
- random + Base62

# Enhancement
- Algorithm
  - Solution 0: UUID
  - Solution 1: md5
    - may have collision
    - Extract first 6 characters
  - Solution 2: random + Base62
    - To prevent collisions, it is needs to constantly ensure that unique
    - Keep length to 6 characters
    - prevent collision
  - Solution 3: counter + Base62 
    - Predictable
    - In a decentralized architecture is a challenge
  - Solution 4: Zookeeper