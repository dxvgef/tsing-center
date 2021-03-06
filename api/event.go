package api

import (
	"net/http"
	"strconv"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"local/global"
)

// 事件处理器
func EventHandler(event tsing.Event) {
	event.ResponseWriter.WriteHeader(event.Status)

	switch event.Status {
	case 404:
		log.Error().Int("status", event.Status).
			Str("method", event.Request.Method).
			Str("uri", event.Request.RequestURI).Msg(http.StatusText(404))
	case 405:
		log.Error().Int("status", event.Status).
			Str("method", event.Request.Method).
			Str("uri", event.Request.RequestURI).Msg(http.StatusText(405))
	case 500:
		e := log.Err(event.Message)
		e.Str("caller", " "+event.Source.File+":"+strconv.Itoa(event.Source.Line)+" ").
			Str("func", event.Source.Func)
		var trace []string
		for k := range event.Trace {
			trace = append(trace, event.Trace[k])
		}
		e.Strs("trace", trace)

		e.Send()
	}

	if _, err := event.ResponseWriter.Write(global.StrToBytes(event.Message.Error())); err != nil {
		log.Err(err).Caller().Send()
	}
}
