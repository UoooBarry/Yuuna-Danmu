package main

import (
	"embed"
	"io"
	"net/http"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

func proxyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/proxy" {
			targetURL := r.URL.Query().Get("url")
			if targetURL == "" {
				http.Error(w, "URL is required", http.StatusBadRequest)
				return
			}

			if strings.HasPrefix(targetURL, "//") {
				targetURL = "https:" + targetURL
			}

			req, _ := http.NewRequest("GET", targetURL, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")
			req.Header.Set("Referer", "https://www.bilibili.com/")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.Header().Set("Cache-Control", "public, max-age=86400")
			io.Copy(w, resp.Body)
			return
		}
		next.ServeHTTP(w, r)
	})
}
