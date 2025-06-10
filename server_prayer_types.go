package main

type PrayerRequest struct {
	Receiver string `json: "receiver"`
	Prayer   string `json: "prayer"`
}
