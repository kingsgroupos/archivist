// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package conf

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

func easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(in *jlexer.Lexer, out *ExampleArrayConf_1771430469) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "avatar_url":
			out.AvatarUrl = string(in.String())
		case "blog":
			out.Blog = string(in.String())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "events_url":
			out.EventsUrl = string(in.String())
		case "followers":
			out.Followers = int64(in.Int64())
		case "followers_url":
			out.FollowersUrl = string(in.String())
		case "following":
			out.Following = int64(in.Int64())
		case "following_url":
			out.FollowingUrl = string(in.String())
		case "gists_url":
			out.GistsUrl = string(in.String())
		case "gravatar_id":
			out.GravatarId = string(in.String())
		case "hireable":
			out.Hireable = bool(in.Bool())
		case "html_url":
			out.HtmlUrl = string(in.String())
		case "id":
			out.Id = int64(in.Int64())
		case "location":
			out.Location = string(in.String())
		case "login":
			out.Login = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "organizations_url":
			out.OrganizationsUrl = string(in.String())
		case "public_gists":
			out.PublicGists = int64(in.Int64())
		case "public_repos":
			out.PublicRepos = int64(in.Int64())
		case "received_events_url":
			out.ReceivedEventsUrl = string(in.String())
		case "repos_url":
			out.ReposUrl = string(in.String())
		case "site_admin":
			out.SiteAdmin = bool(in.Bool())
		case "starred_url":
			out.StarredUrl = string(in.String())
		case "subscriptions_url":
			out.SubscriptionsUrl = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "updated_at":
			out.UpdatedAt = string(in.String())
		case "url":
			out.Url = string(in.String())
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
func easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(out *jwriter.Writer, in ExampleArrayConf_1771430469) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"avatar_url\":"
		out.RawString(prefix[1:])
		out.String(string(in.AvatarUrl))
	}
	{
		const prefix string = ",\"blog\":"
		out.RawString(prefix)
		out.String(string(in.Blog))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"events_url\":"
		out.RawString(prefix)
		out.String(string(in.EventsUrl))
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		out.Int64(int64(in.Followers))
	}
	{
		const prefix string = ",\"followers_url\":"
		out.RawString(prefix)
		out.String(string(in.FollowersUrl))
	}
	{
		const prefix string = ",\"following\":"
		out.RawString(prefix)
		out.Int64(int64(in.Following))
	}
	{
		const prefix string = ",\"following_url\":"
		out.RawString(prefix)
		out.String(string(in.FollowingUrl))
	}
	{
		const prefix string = ",\"gists_url\":"
		out.RawString(prefix)
		out.String(string(in.GistsUrl))
	}
	{
		const prefix string = ",\"gravatar_id\":"
		out.RawString(prefix)
		out.String(string(in.GravatarId))
	}
	{
		const prefix string = ",\"hireable\":"
		out.RawString(prefix)
		out.Bool(bool(in.Hireable))
	}
	{
		const prefix string = ",\"html_url\":"
		out.RawString(prefix)
		out.String(string(in.HtmlUrl))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix)
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"organizations_url\":"
		out.RawString(prefix)
		out.String(string(in.OrganizationsUrl))
	}
	{
		const prefix string = ",\"public_gists\":"
		out.RawString(prefix)
		out.Int64(int64(in.PublicGists))
	}
	{
		const prefix string = ",\"public_repos\":"
		out.RawString(prefix)
		out.Int64(int64(in.PublicRepos))
	}
	{
		const prefix string = ",\"received_events_url\":"
		out.RawString(prefix)
		out.String(string(in.ReceivedEventsUrl))
	}
	{
		const prefix string = ",\"repos_url\":"
		out.RawString(prefix)
		out.String(string(in.ReposUrl))
	}
	{
		const prefix string = ",\"site_admin\":"
		out.RawString(prefix)
		out.Bool(bool(in.SiteAdmin))
	}
	{
		const prefix string = ",\"starred_url\":"
		out.RawString(prefix)
		out.String(string(in.StarredUrl))
	}
	{
		const prefix string = ",\"subscriptions_url\":"
		out.RawString(prefix)
		out.String(string(in.SubscriptionsUrl))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"updated_at\":"
		out.RawString(prefix)
		out.String(string(in.UpdatedAt))
	}
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		out.String(string(in.Url))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExampleArrayConf_1771430469) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExampleArrayConf_1771430469) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExampleArrayConf_1771430469) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExampleArrayConf_1771430469) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConfExampleArrayConf(l, v)
}
func easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(in *jlexer.Lexer, out *ExampleArrayConf) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ExampleArrayConf, 0, 8)
			} else {
				*out = ExampleArrayConf{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *ExampleArrayConf_1771430469
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(ExampleArrayConf_1771430469)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(out *jwriter.Writer, in ExampleArrayConf) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ExampleArrayConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExampleArrayConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson957c5883EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExampleArrayConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExampleArrayConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson957c5883DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(l, v)
}