package main

import (
  "math"
  "strings"
  "fmt"
  "net/http"
  "sync"
  "log"
)

var (
  urlStore    = make(map[string]string)
  urlStoreMux = sync.RWMutex{}
  id          = 0
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(num int) string{
  var result []byte
  
  for num > 0 {
    result = append([]byte{base62[num%62]}, result...)
    num /= 62
  }

  return string(result)
}

func DecodeBase62(str string) int {
  var result int

  for i, char := range str {
    result += strings.IndexByte(base62, byte(char)) * int(math.Pow(62, float64(len(str)-i-1)))
  }

  return result
}

func generateShortURL(url string )string{
  urlStoreMux.Lock()
  defer urlStoreMux.Unlock()
  id++
  shortURL := EncodeBase62(id)
  urlStore[shortURL] = url
  return shortURL
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    longURL := r.FormValue("url")
    if longURL == "" {
      http.Error(w,"Missing URL", http.StatusBadRequest)
      return
    }
    shortURL := generateShortURL(longURL)
    fmt.Fprintf(w, "Shortened URL: http://localhost:3000/%s\n", shortURL)
    return
  }

  http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
  shortURL := r.URL.Path[len("/"):]
  urlStoreMux.RLock()
  defer urlStoreMux.RUnlock()

  longURL, exists := urlStore[shortURL]
  if !exists {
    http.NotFound(w,r)
    return
  }
  http.Redirect(w, r, longURL, http.StatusFound)
}

func main(){
  http.HandleFunc("/short-url", shortenerHandler)
  http.HandleFunc("/", redirectHandler)
  
  fmt.Println("Starting server on http://localhost:3000")
  log.Fatal(http.ListenAndServe(":3000", nil))
}
