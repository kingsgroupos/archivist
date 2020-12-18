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

func easyjsonB418addaDecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(in *jlexer.Lexer, out *BConf) {
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
		case "B1":
			out.B1 = int64(in.Int64())
		case "B2":
			out.B2 = int64(in.Int64())
		case "B3":
			out.B3 = int64(in.Int64())
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
func easyjsonB418addaEncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(out *jwriter.Writer, in BConf) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"B1\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.B1))
	}
	{
		const prefix string = ",\"B2\":"
		out.RawString(prefix)
		out.Int64(int64(in.B2))
	}
	{
		const prefix string = ",\"B3\":"
		out.RawString(prefix)
		out.Int64(int64(in.B3))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB418addaEncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB418addaEncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB418addaDecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB418addaDecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(l, v)
}
