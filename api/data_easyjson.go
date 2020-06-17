// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package api

import (
	json "encoding/json"
	global "github.com/dxvgef/tsing-center/global"
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

func easyjson794297d0DecodeGithubComDxvgefTsingCenterApi(in *jlexer.Lexer, out *Data) {
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
		case "services":
			if in.IsNull() {
				in.Skip()
				out.Services = nil
			} else {
				in.Delim('[')
				if out.Services == nil {
					if !in.IsDelim(']') {
						out.Services = make([]global.ServiceType, 0, 2)
					} else {
						out.Services = []global.ServiceType{}
					}
				} else {
					out.Services = (out.Services)[:0]
				}
				for !in.IsDelim(']') {
					var v1 global.ServiceType
					(v1).UnmarshalEasyJSON(in)
					out.Services = append(out.Services, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "nodes":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Nodes = make(map[string][]global.NodeType)
				} else {
					out.Nodes = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 []global.NodeType
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						in.Delim('[')
						if v2 == nil {
							if !in.IsDelim(']') {
								v2 = make([]global.NodeType, 0, 2)
							} else {
								v2 = []global.NodeType{}
							}
						} else {
							v2 = (v2)[:0]
						}
						for !in.IsDelim(']') {
							var v3 global.NodeType
							(v3).UnmarshalEasyJSON(in)
							v2 = append(v2, v3)
							in.WantComma()
						}
						in.Delim(']')
					}
					(out.Nodes)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
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
func easyjson794297d0EncodeGithubComDxvgefTsingCenterApi(out *jwriter.Writer, in Data) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Services) != 0 {
		const prefix string = ",\"services\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v4, v5 := range in.Services {
				if v4 > 0 {
					out.RawByte(',')
				}
				(v5).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.Nodes) != 0 {
		const prefix string = ",\"nodes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v6First := true
			for v6Name, v6Value := range in.Nodes {
				if v6First {
					v6First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v6Name))
				out.RawByte(':')
				if v6Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v7, v8 := range v6Value {
						if v7 > 0 {
							out.RawByte(',')
						}
						(v8).MarshalEasyJSON(out)
					}
					out.RawByte(']')
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Data) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComDxvgefTsingCenterApi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Data) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComDxvgefTsingCenterApi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Data) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComDxvgefTsingCenterApi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Data) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComDxvgefTsingCenterApi(l, v)
}
