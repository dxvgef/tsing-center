// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package global

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

func easyjson9f2eff5fDecodeLocalGlobal(in *jlexer.Lexer, out *ServiceConfig) {
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
		case "service_id":
			out.ServiceID = string(in.String())
		case "load_balance":
			out.LoadBalance = string(in.String())
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
func easyjson9f2eff5fEncodeLocalGlobal(out *jwriter.Writer, in ServiceConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"service_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ServiceID))
	}
	{
		const prefix string = ",\"load_balance\":"
		out.RawString(prefix)
		out.String(string(in.LoadBalance))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ServiceConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeLocalGlobal(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ServiceConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeLocalGlobal(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ServiceConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeLocalGlobal(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ServiceConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeLocalGlobal(l, v)
}
func easyjson9f2eff5fDecodeLocalGlobal1(in *jlexer.Lexer, out *Node) {
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
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = uint16(in.Uint16())
		case "weight":
			out.Weight = int(in.Int())
		case "expires":
			out.Expires = int64(in.Int64())
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
func easyjson9f2eff5fEncodeLocalGlobal1(out *jwriter.Writer, in Node) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ip\":"
		out.RawString(prefix[1:])
		out.String(string(in.IP))
	}
	{
		const prefix string = ",\"port\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Port))
	}
	{
		const prefix string = ",\"weight\":"
		out.RawString(prefix)
		out.Int(int(in.Weight))
	}
	{
		const prefix string = ",\"expires\":"
		out.RawString(prefix)
		out.Int64(int64(in.Expires))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Node) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeLocalGlobal1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Node) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeLocalGlobal1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Node) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeLocalGlobal1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Node) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeLocalGlobal1(l, v)
}
