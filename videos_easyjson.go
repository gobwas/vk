// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package vk

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC3e3fb8cDecodeGithubComGobwasVk(in *jlexer.Lexer, out *Videos) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "count":
			out.Count = int(in.Int())
		case "items":
			if in.IsNull() {
				in.Skip()
				out.Items = nil
			} else {
				in.Delim('[')
				if out.Items == nil {
					if !in.IsDelim(']') {
						out.Items = make([]Video, 0, 1)
					} else {
						out.Items = []Video{}
					}
				} else {
					out.Items = (out.Items)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Video
					(v1).UnmarshalEasyJSON(in)
					out.Items = append(out.Items, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3e3fb8cEncodeGithubComGobwasVk(out *jwriter.Writer, in Videos) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Count))
	}
	{
		const prefix string = ",\"items\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Items == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Items {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Videos) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3e3fb8cEncodeGithubComGobwasVk(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Videos) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3e3fb8cEncodeGithubComGobwasVk(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Videos) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3e3fb8cDecodeGithubComGobwasVk(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Videos) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3e3fb8cDecodeGithubComGobwasVk(l, v)
}
func easyjsonC3e3fb8cDecodeGithubComGobwasVk1(in *jlexer.Lexer, out *Video) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "owner_id":
			out.OwnerID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "duration":
			out.Duration = int(in.Int())
		case "photo_130":
			out.Photo130 = string(in.String())
		case "photo_320":
			out.Photo320 = string(in.String())
		case "photo_640":
			out.Photo640 = string(in.String())
		case "photo_800":
			out.Photo800 = string(in.String())
		case "date":
			out.Date = int(in.Int())
		case "adding_date":
			out.AddingDate = int(in.Int())
		case "views":
			out.Views = int(in.Int())
		case "comments":
			out.Comments = int(in.Int())
		case "player":
			out.Player = string(in.String())
		case "platform":
			out.Platform = string(in.String())
		case "can_edit":
			out.CanEdit = int(in.Int())
		case "can_add":
			out.CanAdd = int(in.Int())
		case "is_private":
			out.IsPrivate = int(in.Int())
		case "access_key":
			out.AccessKey = string(in.String())
		case "processing":
			out.Processing = int(in.Int())
		case "live":
			out.Live = int(in.Int())
		case "upcoming":
			out.Upcoming = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3e3fb8cEncodeGithubComGobwasVk1(out *jwriter.Writer, in Video) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"owner_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.OwnerID))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"duration\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Duration))
	}
	{
		const prefix string = ",\"photo_130\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Photo130))
	}
	{
		const prefix string = ",\"photo_320\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Photo320))
	}
	{
		const prefix string = ",\"photo_640\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Photo640))
	}
	{
		const prefix string = ",\"photo_800\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Photo800))
	}
	{
		const prefix string = ",\"date\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Date))
	}
	{
		const prefix string = ",\"adding_date\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.AddingDate))
	}
	{
		const prefix string = ",\"views\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Views))
	}
	{
		const prefix string = ",\"comments\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Comments))
	}
	{
		const prefix string = ",\"player\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Player))
	}
	{
		const prefix string = ",\"platform\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Platform))
	}
	{
		const prefix string = ",\"can_edit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanEdit))
	}
	{
		const prefix string = ",\"can_add\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanAdd))
	}
	{
		const prefix string = ",\"is_private\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.IsPrivate))
	}
	{
		const prefix string = ",\"access_key\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AccessKey))
	}
	{
		const prefix string = ",\"processing\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Processing))
	}
	{
		const prefix string = ",\"live\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Live))
	}
	{
		const prefix string = ",\"upcoming\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Upcoming))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Video) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3e3fb8cEncodeGithubComGobwasVk1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Video) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3e3fb8cEncodeGithubComGobwasVk1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Video) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3e3fb8cDecodeGithubComGobwasVk1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Video) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3e3fb8cDecodeGithubComGobwasVk1(l, v)
}