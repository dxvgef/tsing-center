// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package etcd

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

func easyjsonCdfae1c8DecodeGithubComDxvgefTsingCenterStorageEtcd(in *jlexer.Lexer, out *NodeValue) {
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
		case "expires":
			out.Expires = int64(in.Int64())
		case "weight":
			out.Weight = int(in.Int())
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
func easyjsonCdfae1c8EncodeGithubComDxvgefTsingCenterStorageEtcd(out *jwriter.Writer, in NodeValue) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Expires != 0 {
		const prefix string = ",\"expires\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.Expires))
	}
	if in.Weight != 0 {
		const prefix string = ",\"weight\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Weight))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NodeValue) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCdfae1c8EncodeGithubComDxvgefTsingCenterStorageEtcd(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NodeValue) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCdfae1c8EncodeGithubComDxvgefTsingCenterStorageEtcd(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NodeValue) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCdfae1c8DecodeGithubComDxvgefTsingCenterStorageEtcd(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NodeValue) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCdfae1c8DecodeGithubComDxvgefTsingCenterStorageEtcd(l, v)
}