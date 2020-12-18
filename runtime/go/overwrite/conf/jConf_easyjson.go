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

func easyjson72d594e2DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(in *jlexer.Lexer, out *JConf) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		*out = make(JConf)
		for !in.IsDelim('}') {
			key := string(in.String())
			in.WantColon()
			var v1 string
			v1 = string(in.String())
			(*out)[key] = v1
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson72d594e2EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(out *jwriter.Writer, in JConf) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v2First := true
		for v2Name, v2Value := range in {
			if v2First {
				v2First = false
			} else {
				out.RawByte(',')
			}
			out.String(string(v2Name))
			out.RawByte(':')
			out.String(string(v2Value))
		}
		out.RawByte('}')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v JConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson72d594e2EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson72d594e2EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson72d594e2DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson72d594e2DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(l, v)
}
