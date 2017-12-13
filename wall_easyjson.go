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

func easyjson783c1624DecodeGithubComGobwasVk(in *jlexer.Lexer, out *Posts) {
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
						out.Items = make([]Post, 0, 1)
					} else {
						out.Items = []Post{}
					}
				} else {
					out.Items = (out.Items)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Post
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
func easyjson783c1624EncodeGithubComGobwasVk(out *jwriter.Writer, in Posts) {
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
func (v Posts) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Posts) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Posts) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Posts) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk1(in *jlexer.Lexer, out *PostViews) {
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
func easyjson783c1624EncodeGithubComGobwasVk1(out *jwriter.Writer, in PostViews) {
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostViews) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostViews) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostViews) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostViews) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk1(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk2(in *jlexer.Lexer, out *PostSource) {
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
		case "type":
			out.Type = string(in.String())
		case "platform":
			out.Platform = string(in.String())
		case "data":
			out.Data = string(in.String())
		case "url":
			out.URL = string(in.String())
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
func easyjson783c1624EncodeGithubComGobwasVk2(out *jwriter.Writer, in PostSource) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
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
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Data))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostSource) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostSource) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostSource) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostSource) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk2(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk3(in *jlexer.Lexer, out *PostReposts) {
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
		case "user_reposted":
			out.UserReposted = int(in.Int())
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
func easyjson783c1624EncodeGithubComGobwasVk3(out *jwriter.Writer, in PostReposts) {
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
		const prefix string = ",\"user_reposted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.UserReposted))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostReposts) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostReposts) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostReposts) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostReposts) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk3(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk4(in *jlexer.Lexer, out *PostLikes) {
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
		case "user_likes":
			out.UserLikes = int(in.Int())
		case "can_like":
			out.CanLike = int(in.Int())
		case "can_publish":
			out.CanPublish = int(in.Int())
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
func easyjson783c1624EncodeGithubComGobwasVk4(out *jwriter.Writer, in PostLikes) {
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
		const prefix string = ",\"user_likes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.UserLikes))
	}
	{
		const prefix string = ",\"can_like\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanLike))
	}
	{
		const prefix string = ",\"can_publish\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanPublish))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostLikes) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostLikes) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostLikes) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostLikes) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk4(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk5(in *jlexer.Lexer, out *PostComments) {
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
		case "can_post":
			out.CanPost = int(in.Int())
		case "groups_can_post":
			out.GroupsCanPost = bool(in.Bool())
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
func easyjson783c1624EncodeGithubComGobwasVk5(out *jwriter.Writer, in PostComments) {
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
		const prefix string = ",\"can_post\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanPost))
	}
	{
		const prefix string = ",\"groups_can_post\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.GroupsCanPost))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostComments) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostComments) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostComments) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostComments) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk5(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk6(in *jlexer.Lexer, out *PostAttachement) {
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
		case "type":
			out.Type = string(in.String())
		case "photo":
			(out.Photo).UnmarshalEasyJSON(in)
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
func easyjson783c1624EncodeGithubComGobwasVk6(out *jwriter.Writer, in PostAttachement) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"photo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Photo).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostAttachement) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostAttachement) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostAttachement) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostAttachement) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk6(l, v)
}
func easyjson783c1624DecodeGithubComGobwasVk7(in *jlexer.Lexer, out *Post) {
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
		case "from_id":
			out.FromID = int(in.Int())
		case "created_by":
			out.CreatedBy = int(in.Int())
		case "date":
			out.Date = int(in.Int())
		case "text":
			out.Text = string(in.String())
		case "reply_owner_id":
			out.ReplyOwnerID = int(in.Int())
		case "reply_post_id":
			out.ReplyPostID = int(in.Int())
		case "friends_only":
			out.FriendsOnly = int(in.Int())
		case "post_type":
			out.PostType = string(in.String())
		case "can_pin":
			out.CanPin = int(in.Int())
		case "can_delete":
			out.CanDelete = int(in.Int())
		case "can_edit":
			out.CanEdit = int(in.Int())
		case "is_pinned":
			out.IsPinned = int(in.Int())
		case "marked_as_ads":
			out.MarkedAsAds = int(in.Int())
		case "signer_id":
			out.SignerID = int(in.Int())
		case "comments":
			(out.Comments).UnmarshalEasyJSON(in)
		case "likes":
			(out.Likes).UnmarshalEasyJSON(in)
		case "reposts":
			(out.Reposts).UnmarshalEasyJSON(in)
		case "views":
			(out.Views).UnmarshalEasyJSON(in)
		case "post_source":
			(out.PostSource).UnmarshalEasyJSON(in)
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]PostAttachement, 0, 1)
					} else {
						out.Attachments = []PostAttachement{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v4 PostAttachement
					(v4).UnmarshalEasyJSON(in)
					out.Attachments = append(out.Attachments, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "geo":
			(out.Geo).UnmarshalEasyJSON(in)
		case "copy_history":
			if in.IsNull() {
				in.Skip()
				out.CopyHistory = nil
			} else {
				in.Delim('[')
				if out.CopyHistory == nil {
					if !in.IsDelim(']') {
						out.CopyHistory = make([]Post, 0, 1)
					} else {
						out.CopyHistory = []Post{}
					}
				} else {
					out.CopyHistory = (out.CopyHistory)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Post
					(v5).UnmarshalEasyJSON(in)
					out.CopyHistory = append(out.CopyHistory, v5)
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
func easyjson783c1624EncodeGithubComGobwasVk7(out *jwriter.Writer, in Post) {
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
		const prefix string = ",\"from_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.FromID))
	}
	{
		const prefix string = ",\"created_by\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CreatedBy))
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
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"reply_owner_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.ReplyOwnerID))
	}
	{
		const prefix string = ",\"reply_post_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.ReplyPostID))
	}
	{
		const prefix string = ",\"friends_only\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.FriendsOnly))
	}
	{
		const prefix string = ",\"post_type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PostType))
	}
	{
		const prefix string = ",\"can_pin\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanPin))
	}
	{
		const prefix string = ",\"can_delete\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.CanDelete))
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
		const prefix string = ",\"is_pinned\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.IsPinned))
	}
	{
		const prefix string = ",\"marked_as_ads\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.MarkedAsAds))
	}
	{
		const prefix string = ",\"signer_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.SignerID))
	}
	{
		const prefix string = ",\"comments\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Comments).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"likes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Likes).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"reposts\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Reposts).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"views\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Views).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"post_source\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.PostSource).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"attachments\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.Attachments {
				if v6 > 0 {
					out.RawByte(',')
				}
				(v7).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"geo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Geo).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"copy_history\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.CopyHistory == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.CopyHistory {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson783c1624EncodeGithubComGobwasVk7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson783c1624EncodeGithubComGobwasVk7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson783c1624DecodeGithubComGobwasVk7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson783c1624DecodeGithubComGobwasVk7(l, v)
}
