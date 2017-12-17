package vk

//go:generate easyjson -all

type Groups struct {
	Count int     `json:"count"`
	Items []Group `json:"items"`
}

type Group struct {
	// TODO
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
