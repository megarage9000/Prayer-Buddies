package main

import (
	"fmt"
	"net/http"
	"time"

	// Importing the sql drivers
	_ "github.com/lib/pq"
)

const TOKEN_EXPIRY = time.Hour
const ISSUER = "Prayer Buddies"

func main() {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Port),
		Handler:           serverMux,
		ReadHeaderTimeout: time.Minute,
	}

	serverMux.HandleFunc("POST /api/users", middlewareCORS(config.CreateUser, config.FrontendURL))
	serverMux.HandleFunc("POST /api/login", middlewareCORS(config.LoginUser, config.FrontendURL))

	serverMux.HandleFunc("POST /api/sendprayer", middlewareCORS(config.SendPrayerRequest, config.FrontendURL))
	serverMux.HandleFunc("/api/receivedRequests", middlewareCORS(config.ListReceivedPrayerRequests, config.FrontendURL))
	serverMux.HandleFunc("/api/sentRequests", middlewareCORS(config.ListSentPrayerRequests, config.FrontendURL))

	serverMux.HandleFunc("POST /api/sendFriendReq", middlewareCORS(config.SendFriendRequest, config.FrontendURL))
	serverMux.HandleFunc("POST /api/updateFriendReq", middlewareCORS(config.UpdateFriendRequest, config.FrontendURL))

	fmt.Printf("Loading on localhost:%s", config.Port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
}
