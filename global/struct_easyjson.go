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

func easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal(in *jlexer.Lexer, out *ServiceType) {
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
		case "id":
			out.ID = string(in.String())
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
func easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal(out *jwriter.Writer, in ServiceType) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.LoadBalance != "" {
		const prefix string = ",\"load_balance\":"
		out.RawString(prefix)
		out.String(string(in.LoadBalance))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ServiceType) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ServiceType) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ServiceType) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ServiceType) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal(l, v)
}
func easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal1(in *jlexer.Lexer, out *NodeType) {
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
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = int(in.Int())
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
func easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal1(out *jwriter.Writer, in NodeType) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"service_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ServiceID))
	}
	{
		const prefix string = ",\"ip\":"
		out.RawString(prefix)
		out.String(string(in.IP))
	}
	{
		const prefix string = ",\"port\":"
		out.RawString(prefix)
		out.Int(int(in.Port))
	}
	{
		const prefix string = ",\"weight\":"
		out.RawString(prefix)
		out.Int(int(in.Weight))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NodeType) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NodeType) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NodeType) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NodeType) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal1(l, v)
}
func easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal2(in *jlexer.Lexer, out *ModuleConfig) {
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
		case "name":
			out.Name = string(in.String())
		case "config":
			out.Config = string(in.String())
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
func easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal2(out *jwriter.Writer, in ModuleConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"config\":"
		out.RawString(prefix)
		out.String(string(in.Config))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ModuleConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ModuleConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeGithubComDxvgefTsingCenterGlobal2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ModuleConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ModuleConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeGithubComDxvgefTsingCenterGlobal2(l, v)
}
