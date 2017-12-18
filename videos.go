package vk

//go:generate easyjson -all

type Videos struct {
	Count int     `json:"count"`
	Items []Video `json:"items"`
}

type Video struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Photo130    string `json:"photo_130"`
	Photo320    string `json:"photo_320"`
	Photo640    string `json:"photo_640"`
	Photo800    string `json:"photo_800"`
	Date        int    `json:"date"`
	AddingDate  int    `json:"adding_date"`
	Views       int    `json:"views"`
	Comments    int    `json:"comments"`
	Player      string `json:"player"`
	Platform    string `json:"platform"`
	CanEdit     int    `json:"can_edit"`
	CanAdd      int    `json:"can_add"`
	IsPrivate   int    `json:"is_private"`
	AccessKey   string `json:"access_key"`
	Processing  int    `json:"processing"`
	Live        int    `json:"live"`
	Upcoming    int    `json:"upcoming"`
}
