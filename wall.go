package vk

//go:generate easyjson -all

type Posts struct {
	Count int    `json:"count"`
	Items []Post `json:"items"`
}

type Post struct {
	ID           int    `json:"id"`
	OwnerID      int    `json:"owner_id"`
	FromID       int    `json:"from_id"`
	CreatedBy    int    `json:"created_by"`
	Date         int    `json:"date"`
	Text         string `json:"text"`
	ReplyOwnerID int    `json:"reply_owner_id"`
	ReplyPostID  int    `json:"reply_post_id"`
	FriendsOnly  int    `json:"friends_only"`
	PostType     string `json:"post_type"`
	CanPin       int    `json:"can_pin"`
	CanDelete    int    `json:"can_delete"`
	CanEdit      int    `json:"can_edit"`
	IsPinned     int    `json:"is_pinned"`
	MarkedAsAds  int    `json:"marked_as_ads"`
	SignerID     int    `json:"signer_id"`

	Comments    PostComments      `json:"comments"`
	Likes       PostLikes         `json:"likes"`
	Reposts     PostReposts       `json:"reposts"`
	Views       PostViews         `json:"views"`
	PostSource  PostSource        `json:"post_source"`
	Attachments []PostAttachement `json:"attachments"`
	Geo         Geo               `json:"geo"`
	CopyHistory []Post            `json:"copy_history"`
}

type PostComments struct {
	Count         int  `json:"count"`
	CanPost       int  `json:"can_post"`
	GroupsCanPost bool `json:"groups_can_post"`
}

type PostLikes struct {
	Count      int `json:"count"`
	UserLikes  int `json:"user_likes"`
	CanLike    int `json:"can_like"`
	CanPublish int `json:"can_publish"`
}

type PostReposts struct {
	Count        int `json:"count"`
	UserReposted int `json:"user_reposted"`
}

type PostViews struct {
	Count int `json:"count"`
}

type PostAttachement struct {
	//TODO
	Type  string `json:"type"`
	Photo Photo  `json:"photo"`
}

type PostSource struct {
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Data     string `json:"data"`
	URL      string `json:"url"`
}
