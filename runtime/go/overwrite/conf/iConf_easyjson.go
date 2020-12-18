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

func easyjson809b5f19DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(in *jlexer.Lexer, out *IConf) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		*out = make(IConf)
		for !in.IsDelim('}') {
			key := string(in.String())
			in.WantColon()
			var v1 map[string]int64
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				v1 = make(map[string]int64)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 int64
					v2 = int64(in.Int64())
					(v1)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
			(*out)[key] = v1
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson809b5f19EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(out *jwriter.Writer, in IConf) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v3First := true
		for v3Name, v3Value := range in {
			if v3First {
				v3First = false
			} else {
				out.RawByte(',')
			}
			out.String(string(v3Name))
			out.RawByte(':')
			if v3Value == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
				out.RawString(`null`)
			} else {
				out.RawByte('{')
				v4First := true
				for v4Name, v4Value := range v3Value {
					if v4First {
						v4First = false
					} else {
						out.RawByte(',')
					}
					out.String(string(v4Name))
					out.RawByte(':')
					out.Int64(int64(v4Value))
				}
				out.RawByte('}')
			}
		}
		out.RawByte('}')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v IConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson809b5f19EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson809b5f19EncodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *IConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson809b5f19DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson809b5f19DecodeGitlabEeFunplusIoWatcherWatcherArchivistRuntimeGoOverwriteConf(l, v)
}
