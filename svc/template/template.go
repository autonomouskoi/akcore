// Package template provides go text/template rendering as a service
package template

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"google.golang.org/protobuf/proto"
)

// Template service
type Template struct {
	modutil.ModuleBase
}

// New creates a new Template
func New(deps *modutil.Deps) (*Template, error) {
	t := &Template{}
	t.Log = deps.Log.NewForSource("svc.template")
	return t, nil
}

// HandleRequestRender handles a template rendering request
func (t *Template) HandleRequestRender(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	req := &svc.TemplateRenderRequest{}
	if reply.Error = t.UnmarshalMessage(msg, req); reply.Error != nil {
		return reply
	}

	if req.Template == "" {
		reply.Error = &bus.Error{
			Detail: proto.String("no template"),
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
		}
		return reply
	}

	tmpl, err := template.New("").Parse(req.GetTemplate())
	if err != nil {
		reply.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String(err.Error()),
		}
		return reply
	}

	var data any
	if len(req.GetJson()) != 0 {
		if err := json.Unmarshal(req.GetJson(), &data); err != nil {
			reply.Error = &bus.Error{
				Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
				Detail: proto.String(err.Error()),
			}
			return reply
		}
	}

	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, data); err != nil {
		reply.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String(err.Error()),
		}
		return reply
	}

	t.MarshalMessage(reply, &svc.TemplateRenderResponse{Output: buf.String()})
	return reply
}
