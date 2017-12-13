package vk

//go:generate easyjson -all

type Geo struct {
	Type        string   `json:"type"`
	Coordinates string   `json:"coordinates"`
	Place       GeoPlace `json:"place"`
}

type GeoPlace struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Created   int     `json:"created"`
	Icon      string  `json:"icon"`
	Country   string  `json:"country"`
	City      string  `json:"city"`

	// Checkin additional fields.
	Type       int    `json:"type"`
	GroupID    int    `json:"group_id"`
	GroupPhoto string `json:"group_photo"`
	Checkins   int    `json:"checkins"`
	Updated    int    `json:"updated"`
	Address    string `json:"address"`
}
