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

func easyjsonE3778cd9DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(in *jlexer.Lexer, out *ExtConf) {
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
		case "D":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.D).UnmarshalJSON(data))
			}
		case "T":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.T).UnmarshalJSON(data))
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
func easyjsonE3778cd9EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(out *jwriter.Writer, in ExtConf) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"D\":"
		out.RawString(prefix[1:])
		out.Raw((in.D).MarshalJSON())
	}
	{
		const prefix string = ",\"T\":"
		out.RawString(prefix)
		out.Raw((in.T).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExtConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE3778cd9EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExtConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE3778cd9EncodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExtConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE3778cd9DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExtConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE3778cd9DecodeGitlabEeFunplusIoWatcherWatcherArchivistCliArchivistExampleConf(l, v)
}
