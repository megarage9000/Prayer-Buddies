package main

import (
	"net/http"
)

func helloWorld(resp http.ResponseWriter, req *http.Request) {
	helloWorld := "Hello World!"
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(helloWorld))
}
