package web

import (
	"encoding/json"
)

// View for a web interaction.
// Implementations are responsible for rendering content, and exposing the model.
type View interface {
	// ContentType returns HTTP Content-type
	ContentType() string
	// Render the view given the specified model.
	Render(ctx Context) error
}

type ActionView struct {
	local       string
	contentType string
}

func (v *ActionView) ContentType() string {
	if v.contentType == "" {
		return "application/octet-stream; charset=UTF-8"
	}
	return v.contentType
}

type ContentView struct {
	*ActionView
	Content string
}

func (ar *ContentView) Render(ctx Context) error {
	w := ctx.Response().Writer()
	// w.Header().Set("Content-Type", ar.ContentType())
	_, err := w.Write([]byte(ar.Content))
	return err
}

func (ctx *RequestContext) Content(content string, contentType ...string) View {

	view := &ContentView{
		ActionView: &ActionView{
			local: ctx.local,
		},
		Content: content,
	}
	if len(contentType) > 0 {
		view.contentType = contentType[0]
	}
	return view
}

type JsonView struct {
	*ActionView
	Data interface{}
}

func (v *JsonView) Render(ctx Context) error {
	content, err := json.Marshal(v.Data)
	if err != nil {
		return err
	}
	w := ctx.Response().Writer()
	// w.Header().Set("Content-Type", v.ContentType())
	_, err = w.Write([]byte(content))
	return err
}

func (ctx *RequestContext) JSON(data interface{}, contentType ...string) View {
	view := &JsonView{
		ActionView: &ActionView{
			local: ctx.local,
		},
		Data: data,
	}
	if len(contentType) > 0 {
		view.contentType = contentType[0]
	} else {
		view.contentType = "application/json; charset=UTF-8"
	}
	return view
}

type TemplateView struct {
	*ActionView
	Template string
}

func (v *TemplateView) Render(ctx Context) error {
	return ctx.App().viewEngine.Render(ctx.ViewData(), v.Template, ctx.Response().Writer())
}

func (ctx *RequestContext) View(template string, contentType ...string) View {
	view := &TemplateView{
		ActionView: &ActionView{
			local: ctx.local,
		},
		Template: template,
	}
	if len(contentType) > 0 {
		view.contentType = contentType[0]
	} else {
		view.contentType = "text/html; charset=UTF-8"
	}
	return view
}

type emptyView struct {
	*ActionView
}

func (*emptyView) Render(ctx Context) error {
	return nil
}

func (ctx *RequestContext) Empty() View {
	return &emptyView{
		ActionView: &ActionView{
			local: ctx.local,
		},
	}
}
