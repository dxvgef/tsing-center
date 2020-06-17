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
			for v2, v3 := range in.Services {
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
