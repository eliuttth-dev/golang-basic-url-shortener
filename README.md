# Go URL Shortener

A simple URL shortener built with Go using Base62 encoding. 

## Usage
    go run main.go

    The server will start running at `http://localhost:3000`.
Send a `POST` request to `/short-url` with the long URL you want to shorten. For example:

## Bash
```
curl -X POST -d "url=https://www.example.com" http://localhost:3000/short-url

Shortened URL: http://localhost:3000/1
```


