package template_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/autonomouskoi/akcore/svc/template"
)

func TestTemplate(t *testing.T) {
	t.Parallel()

	deps := &modutil.ModuleDeps{
		Log: nil,
	}

	tmpl, err := template.New(deps)
	require.NoError(t, err, "getting template")

	makeReq := func(t *testing.T, template, js string) *bus.BusMessage {
		t.Helper()
		msg := &bus.BusMessage{
			Type: int32(svc.MessageType_TEMPLATE_RENDER_REQ),
		}
		b, err := proto.Marshal(&svc.TemplateRenderRequest{
			Template: template,
			Json:     []byte(js),
		})
		require.NoError(t, err, "marshalling request")
		msg.Message = b
		return msg
	}

	handleReply := func(t *testing.T, reply *bus.BusMessage) string {
		t.Helper()
		resp := &svc.TemplateRenderResponse{}
		require.NoError(t, proto.Unmarshal(reply.GetMessage(), resp), "unmarshalling reply")
		return resp.GetOutput()
	}

	t.Run("no template", func(t *testing.T) {
		t.Parallel()
		msg := makeReq(t, "", "")
		reply := tmpl.HandleRequestRender(msg)
		require.NotNil(t, reply.Error)
	})

	t.Run("invalid template", func(t *testing.T) {
		t.Parallel()
		msg := makeReq(t, "{{ }", "")
		reply := tmpl.HandleRequestRender(msg)
		require.NotNil(t, reply.Error)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		msg := makeReq(t, "A {{ .Foo }} C", `{ Foo: "B" }`)
		reply := tmpl.HandleRequestRender(msg)
		require.NotNil(t, reply.Error)
	})

	t.Run("valid missing data", func(t *testing.T) {
		t.Parallel()
		msg := makeReq(t, "A {{ .Foo }} C", `{ "Bar": "B" }`)
		reply := tmpl.HandleRequestRender(msg)
		require.Nil(t, reply.Error, reply.Error.GetUserMessage())
		output := handleReply(t, reply)
		require.Equal(t, "A <no value> C", output)
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		msg := makeReq(t, "A {{ .Foo }} C", `{ "Foo": "B" }`)
		reply := tmpl.HandleRequestRender(msg)
		require.Nil(t, reply.Error, reply.Error.GetUserMessage())
		output := handleReply(t, reply)
		require.Equal(t, "A B C", output)
	})
}
