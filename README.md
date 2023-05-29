# URL Shortener

Table of **KEY** Contents
- [System Design Concept](#system-design-concept)
- [Check and Acceptance](#check-and-acceptance)

---

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


## System Design Concept

### 定義需求 

- 基本要求
  - 兩隻 API
    - 產生有期限的短網址
    - 將短網址轉回長網址
      - 找不到短網址則回傳 404
- 額外限制
  - 期限最長三年
  - 短網址長度為 6 碼
  - Validation
  - Uint Test
  - Performance Test
- 附加條件
  - 使用 Go
    - 使用 Go-Kit

### 評估用量

- 假設短網址的效期，最長效期為 3 年
- 推估讀取比建立的 request 多，約 100:1
- 建立 request 的 QPS 預估為 100 request/s，存取的 QPS 為 10,000 request/s
- 短網址的長度為 6 個字元
  - 三年的 url mapping 總共有： $100 * 3600 * 24 * 365 * 3 \approx 10M$ 筆
  - 每一個字元可能有 62 種(0-9a-zA-Z)，所以一共可以有 $62^6 \approx 56M$
  - $56M > 10M$，所以長度 6 就足夠了
- 儲存空間至少需要 20TB
  - 依照各大瀏覽器支援的 URL 長度以及各 Web Service 的預設值，取最大相容長度為 2048
  - $(2048+6) * 10M \approx 20TB$

### 系統設計

- Algorithm
  - Solution 0: UUID
    - Pros: 
      - simple, locally generated
    - Cons: 
      - 16 characters, too long, affects DB performance
  - Solution 1: md5
    - may have collision
    - too long, need to extract first 6 characters
  - Solution 2: counter + Base62
    - Predictable
    - In a decentralized architecture is a challenge
      - a. generate within dividing ranges
      - b. pre-generated
  - Solution 3: random + Base62
    - To prevent collisions, it is needs to constantly ensure that unique
    - Keep length to 6 characters
  - Solution 4: Zookeeper
- Storage: NoSQL
  - There is no relationship between each url mapping
  - Billions of data
  - Read-heavy

### conclusion

因時程考量，使用 random integer and base62 encoded + Redis

## Check and Acceptance

- 基本功能
  - 啟動服務
  - 兩隻 API 的 happy path
- error handling
  - 錯誤的短網址
  - 錯誤的參數
- Unit Test
- Performance Test
- Technical Stack
  - [Go-Kit](https://github.com/go-kit/kit)
  - [Skaffold](https://skaffold.dev/)
  - [Buildpacks](https://buildpacks.io/)
  - [K8s](https://kubernetes.io/)
    - minikube
  - Redis
  - Base62 Encode/Decode

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
