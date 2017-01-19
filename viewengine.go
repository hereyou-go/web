package web

import (
	"io"

	"github.com/hereyou-go/otpl"
	"github.com/hereyou-go/otpl/runtime"
)

// ViewEngine 提供用于将特定模板为渲染为相应输出格式的一组方法
type ViewEngine interface {
	Render(data map[string]interface{}, path string, out io.Writer) error
	//ResolveView(name string, locale string) View
}

//*************** OTPL模板引擎(默认) ***************

type OtplViewEngine struct {
	runtime *runtime.Runtime
}

func (engine *OtplViewEngine) Render(data map[string]interface{}, path string, out io.Writer) error {
	return engine.runtime.Render(data, path, out)
}

// Otpl 视图引擎
func Otpl(ilpath string) *OtplViewEngine {
	engine := &OtplViewEngine{
		runtime: otpl.New(ilpath),
	}
	return engine
}

//*************** 原生模板引擎 ***************

type GoTemplateViewEngine struct {
	runtime *runtime.Runtime
}

func (engine *GoTemplateViewEngine) Render(data map[string]interface{}, path string, out io.Writer) error {
	return engine.runtime.Render(data, path, out)
}

// GoTemplate 将使用 GO 原生模板技术作为视图引擎
func GoTemplate(ilpath string) *GoTemplateViewEngine {
	engine := &GoTemplateViewEngine{}
	return engine
}
