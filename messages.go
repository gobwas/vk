package vk

//go:generate easyjson -all

type Dialogs struct {
	Count int      `json:"count"`
	Items []Dialog `json:"items"`
}

type Dialog struct {
	Unread  int     `json:"unread"`
	Message Message `json:"message"`
	InRead  int     `json:"in_read"`
	OutRead int     `json:"out_read"`
}

type Messages struct {
	Count int       `json:"count"`
	Items []Message `json:"items"`
}

type Message struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	FromID      int           `json:"from_id"`
	Date        int64         `json:"date"`
	ReadState   int           `json:"read_state"`
	Out         int           `json:"out"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	Geo         Geo           `json:"geo"`
	Attachments []Attachement `json:"attachments"`
	FwdMessages []Message     `json:"fwd_messages"`
	Emoji       int           `json:"emoji"`
	Important   int           `json:"important"`
	Deleted     int           `json:"deleted"`
	RandomId    int           `json:"random_id"`
	// Chat fields.
	ChatID      int    `json:"chat_id"`
	ChatActive  []int  `json:"chat_active"`
	UsersCount  int    `json:"users_count"`
	AdminID     int    `json:"admin_id"`
	Action      string `json:"action"`
	ActionMid   int    `json:"action_mid"`
	ActionEmail string `json:"action_email"`
	ActionText  string `json:"action_text"`
	Photo50     string `json:"photo_50"`
	Photo100    string `json:"photo_100"`
	Photo200    string `json:"photo_200"`
}

type Attachement struct {
	Type  string `json:"type"`
	Photo Photo  `json:"photo"`
}
