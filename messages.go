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
}

type Attachement struct {
	Type  string `json:"type"`
	Photo Photo  `json:"photo"`
}
