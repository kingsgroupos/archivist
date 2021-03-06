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

func easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(in *jlexer.Lexer, out *ComplexConf_241695546) {
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
		case "x5a":
			if in.IsNull() {
				in.Skip()
				out.X5A = nil
			} else {
				in.Delim('[')
				if out.X5A == nil {
					if !in.IsDelim(']') {
						out.X5A = make([][]int64, 0, 2)
					} else {
						out.X5A = [][]int64{}
					}
				} else {
					out.X5A = (out.X5A)[:0]
				}
				for !in.IsDelim(']') {
					var v1 []int64
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if v1 == nil {
							if !in.IsDelim(']') {
								v1 = make([]int64, 0, 8)
							} else {
								v1 = []int64{}
							}
						} else {
							v1 = (v1)[:0]
						}
						for !in.IsDelim(']') {
							var v2 int64
							v2 = int64(in.Int64())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.X5A = append(out.X5A, v1)
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
func easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(out *jwriter.Writer, in ComplexConf_241695546) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"x5a\":"
		out.RawString(prefix[1:])
		if in.X5A == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.X5A {
				if v3 > 0 {
					out.RawByte(',')
				}
				if v4 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v5, v6 := range v4 {
						if v5 > 0 {
							out.RawByte(',')
						}
						out.Int64(int64(v6))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ComplexConf_241695546) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ComplexConf_241695546) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ComplexConf_241695546) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ComplexConf_241695546) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConfComplexConf(l, v)
}
func easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf(in *jlexer.Lexer, out *ComplexConfItem) {
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
		case "x3":
			if in.IsNull() {
				in.Skip()
				out.X3 = nil
			} else {
				in.Delim('[')
				if out.X3 == nil {
					if !in.IsDelim(']') {
						out.X3 = make([][]int64, 0, 2)
					} else {
						out.X3 = [][]int64{}
					}
				} else {
					out.X3 = (out.X3)[:0]
				}
				for !in.IsDelim(']') {
					var v7 []int64
					if in.IsNull() {
						in.Skip()
						v7 = nil
					} else {
						in.Delim('[')
						if v7 == nil {
							if !in.IsDelim(']') {
								v7 = make([]int64, 0, 8)
							} else {
								v7 = []int64{}
							}
						} else {
							v7 = (v7)[:0]
						}
						for !in.IsDelim(']') {
							var v8 int64
							v8 = int64(in.Int64())
							v7 = append(v7, v8)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.X3 = append(out.X3, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "x5":
			if in.IsNull() {
				in.Skip()
				out.X5 = nil
			} else {
				in.Delim('[')
				if out.X5 == nil {
					if !in.IsDelim(']') {
						out.X5 = make([]*ComplexConf_241695546, 0, 8)
					} else {
						out.X5 = []*ComplexConf_241695546{}
					}
				} else {
					out.X5 = (out.X5)[:0]
				}
				for !in.IsDelim(']') {
					var v9 *ComplexConf_241695546
					if in.IsNull() {
						in.Skip()
						v9 = nil
					} else {
						if v9 == nil {
							v9 = new(ComplexConf_241695546)
						}
						(*v9).UnmarshalEasyJSON(in)
					}
					out.X5 = append(out.X5, v9)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "x6":
			if in.IsNull() {
				in.Skip()
				out.X6 = nil
			} else {
				in.Delim('[')
				if out.X6 == nil {
					if !in.IsDelim(']') {
						out.X6 = make([]map[string][][]int64, 0, 8)
					} else {
						out.X6 = []map[string][][]int64{}
					}
				} else {
					out.X6 = (out.X6)[:0]
				}
				for !in.IsDelim(']') {
					var v10 map[string][][]int64
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v10 = make(map[string][][]int64)
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v11 [][]int64
							if in.IsNull() {
								in.Skip()
								v11 = nil
							} else {
								in.Delim('[')
								if v11 == nil {
									if !in.IsDelim(']') {
										v11 = make([][]int64, 0, 2)
									} else {
										v11 = [][]int64{}
									}
								} else {
									v11 = (v11)[:0]
								}
								for !in.IsDelim(']') {
									var v12 []int64
									if in.IsNull() {
										in.Skip()
										v12 = nil
									} else {
										in.Delim('[')
										if v12 == nil {
											if !in.IsDelim(']') {
												v12 = make([]int64, 0, 8)
											} else {
												v12 = []int64{}
											}
										} else {
											v12 = (v12)[:0]
										}
										for !in.IsDelim(']') {
											var v13 int64
											v13 = int64(in.Int64())
											v12 = append(v12, v13)
											in.WantComma()
										}
										in.Delim(']')
									}
									v11 = append(v11, v12)
									in.WantComma()
								}
								in.Delim(']')
							}
							(v10)[key] = v11
							in.WantComma()
						}
						in.Delim('}')
					}
					out.X6 = append(out.X6, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "x7":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.X7 = make(map[string]int64)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v14 int64
					v14 = int64(in.Int64())
					(out.X7)[key] = v14
					in.WantComma()
				}
				in.Delim('}')
			}
		case "xx":
			if in.IsNull() {
				in.Skip()
				out.Xx = nil
			} else {
				in.Delim('[')
				if out.Xx == nil {
					if !in.IsDelim(']') {
						out.Xx = make([][]int64, 0, 2)
					} else {
						out.Xx = [][]int64{}
					}
				} else {
					out.Xx = (out.Xx)[:0]
				}
				for !in.IsDelim(']') {
					var v15 []int64
					if in.IsNull() {
						in.Skip()
						v15 = nil
					} else {
						in.Delim('[')
						if v15 == nil {
							if !in.IsDelim(']') {
								v15 = make([]int64, 0, 8)
							} else {
								v15 = []int64{}
							}
						} else {
							v15 = (v15)[:0]
						}
						for !in.IsDelim(']') {
							var v16 int64
							v16 = int64(in.Int64())
							v15 = append(v15, v16)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Xx = append(out.Xx, v15)
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
func easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf(out *jwriter.Writer, in ComplexConfItem) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"x3\":"
		out.RawString(prefix[1:])
		if in.X3 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.X3 {
				if v17 > 0 {
					out.RawByte(',')
				}
				if v18 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v19, v20 := range v18 {
						if v19 > 0 {
							out.RawByte(',')
						}
						out.Int64(int64(v20))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"x5\":"
		out.RawString(prefix)
		if in.X5 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v21, v22 := range in.X5 {
				if v21 > 0 {
					out.RawByte(',')
				}
				if v22 == nil {
					out.RawString("null")
				} else {
					(*v22).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"x6\":"
		out.RawString(prefix)
		if in.X6 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v23, v24 := range in.X6 {
				if v23 > 0 {
					out.RawByte(',')
				}
				if v24 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v25First := true
					for v25Name, v25Value := range v24 {
						if v25First {
							v25First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v25Name))
						out.RawByte(':')
						if v25Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
							out.RawString("null")
						} else {
							out.RawByte('[')
							for v26, v27 := range v25Value {
								if v26 > 0 {
									out.RawByte(',')
								}
								if v27 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
									out.RawString("null")
								} else {
									out.RawByte('[')
									for v28, v29 := range v27 {
										if v28 > 0 {
											out.RawByte(',')
										}
										out.Int64(int64(v29))
									}
									out.RawByte(']')
								}
							}
							out.RawByte(']')
						}
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"x7\":"
		out.RawString(prefix)
		if in.X7 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v30First := true
			for v30Name, v30Value := range in.X7 {
				if v30First {
					v30First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v30Name))
				out.RawByte(':')
				out.Int64(int64(v30Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"xx\":"
		out.RawString(prefix)
		if in.Xx == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v31, v32 := range in.Xx {
				if v31 > 0 {
					out.RawByte(',')
				}
				if v32 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v33, v34 := range v32 {
						if v33 > 0 {
							out.RawByte(',')
						}
						out.Int64(int64(v34))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ComplexConfItem) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ComplexConfItem) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ComplexConfItem) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ComplexConfItem) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf(l, v)
}
func easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(in *jlexer.Lexer, out *ComplexConf) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		*out = make(ComplexConf)
		for !in.IsDelim('}') {
			key := int64(in.Int64Str())
			in.WantColon()
			var v35 *ComplexConfItem
			if in.IsNull() {
				in.Skip()
				v35 = nil
			} else {
				if v35 == nil {
					v35 = new(ComplexConfItem)
				}
				(*v35).UnmarshalEasyJSON(in)
			}
			(*out)[key] = v35
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(out *jwriter.Writer, in ComplexConf) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v36First := true
		for v36Name, v36Value := range in {
			if v36First {
				v36First = false
			} else {
				out.RawByte(',')
			}
			out.Int64Str(int64(v36Name))
			out.RawByte(':')
			if v36Value == nil {
				out.RawString("null")
			} else {
				(*v36Value).MarshalEasyJSON(out)
			}
		}
		out.RawByte('}')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ComplexConf) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ComplexConf) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85cf14f0EncodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ComplexConf) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ComplexConf) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85cf14f0DecodeGithubComKingsgrouposArchivistCliArchivistExampleConf1(l, v)
}
