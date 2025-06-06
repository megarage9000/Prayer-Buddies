package main

import "net/http"

/*
To enable cors
- Helped by ChatGPT
*/
func middlewareCORS(next http.HandlerFunc, allowedOrigin string) http.HandlerFunc {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")

		if allowedOrigin == origin {
			resp.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		resp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		next.ServeHTTP(resp, req)
	})
}
