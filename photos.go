package vk

import (
	"fmt"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

//go:generate easyjson -all

type Tags struct {
	Count int   `json:"count"`
	Items []Tag `json:"items"`
}

type Tag struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	PlacerID   int     `json:"placer_id"`
	TaggedName string  `json:"tagged_name"`
	Date       int     `json:"date"`
	X1         float64 `json:"x"`
	Y1         float64 `json:"y"`
	X2         float64 `json:"x2"`
	Y2         float64 `json:"y2"`
	Viewed     int     `json:"viewed"`
}

type PhotoAlbum struct {
	ID             int      `json:"id"`
	ThumbID        int      `json:"thumb_id"`
	ThumbSrc       string   `json:"thumb_src"`
	OwnerID        int      `json:"owner_id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Created        int      `json:"created"`
	Updated        int      `json:"updated"`
	Size           int      `json:"size"`
	ThumbIsLast    int      `json:"thumb_is_last"`
	PrivacyView    []string `json:"privacy_view"`
	PrivacyComment []string `json:"privacy_comment"`
}

type Photo struct {
	ID      int         `json:"id"`
	AlbumID int         `json:"album_id"`
	OwnerID int         `json:"owner_id"`
	UserID  int         `json:"user_id"`
	Text    string      `json:"text"`
	Date    int         `json:"date"`
	PostID  int         `json:"post_id"`
	Sizes   []PhotoSize `json:"sizes"`
}

type PhotoSize struct {
	Src    string   `json:"src"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Type   SizeType `json:"type"`
}

type PhotoAlbums struct {
	Count int          `json:"count"`
	Items []PhotoAlbum `json:"items"`
}

type Photos struct {
	Count int     `json:"count"`
	Items []Photo `json:"items"`
}

type SizeType byte

// Size types in increasing order.
const (
	SizeUnknown SizeType = iota
	SizeS
	SizeM
	SizeX
	SizeO
	SizeP
	SizeQ
	SizeR
	SizeY
	SizeZ
	SizeW
)

func (s SizeType) String() string {
	switch s {
	case SizeS:
		return "s"
	case SizeM:
		return "m"
	case SizeX:
		return "x"
	case SizeO:
		return "o"
	case SizeP:
		return "p"
	case SizeQ:
		return "q"
	case SizeR:
		return "r"
	case SizeY:
		return "y"
	case SizeZ:
		return "z"
	case SizeW:
		return "w"
	default:
		return "_"
	}
}

func (s *SizeType) UnmarshalEasyJSON(in *jlexer.Lexer) {
	str := in.String()
	if len(str) != 1 {
		in.AddError(fmt.Errorf("unexpected size type: %q", str))
		return
	}
	switch byte(str[0]) {
	case 's':
		*s = SizeS
	case 'm':
		*s = SizeM
	case 'x':
		*s = SizeX
	case 'o':
		*s = SizeO
	case 'p':
		*s = SizeP
	case 'q':
		*s = SizeQ
	case 'r':
		*s = SizeR
	case 'y':
		*s = SizeY
	case 'z':
		*s = SizeZ
	case 'w':
		*s = SizeW
	}
}

func (s *SizeType) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(s.String())
}

func (s SizeType) Less(b SizeType) bool {
	return s < b
}
