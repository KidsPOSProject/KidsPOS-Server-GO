package handler

import (
	"fmt"
	"net/http"
)

// Handler はVercelがリクエストを処理するために呼び出す関数です
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go on Vercel!</h1><p>Path: %s</p>", r.URL.Path)
}
