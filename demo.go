package web

import (
	"log"
	"os"

	"github.com/one-go/logs"
	"github.com/one-go/otpl"
)

// import (
//     "log"
//     "github.com/junhwong/mano/crypto"
//     _ "github.com/junhwong/mano/crypto/bcrypt"
// )

type item struct {
	current bool
	name    string
	url     string
}

func (i *item) Fn() string {
	return "hello word"
}

// func (i *item) Hel2() string  {
// 	return "hello word"
// }

func test() {
	log.Print("")
}

func main() {
	// typ:=reflect.TypeOf(test)
	// panic(typ.Kind())

	logs.SetLevel(logs.LDEBUG)

	logs.Debug(logs.NewError(nil, "hello logs"))

	//logs.Debug("print red color text :{red:%s}", "i'm red!")

	//logs.Error("hello {red:%s}!","word")

	return

	data := make(map[string]interface{}, 0)
	data["title"] = "OTPL!"
	data["header"] = "Hello OTPL!"
	items := make([]item, 0)

	items = append(items, item{
		current: true,
		name:    "James",
		url:     "http://example.com",
	})
	items = append(items, item{
		name: "Foo1",
		url:  "http://example.com",
	})
	items = append(items, item{
		name: "Foo2",
		url:  "http://example.com",
	})
	items = append(items, item{
		name: "Foo3",
		url:  "http://example.com",
	})
	data["items"] = items
	data["item"] = &item{}
	// inst:=&item{}
	// elem:=reflect.ValueOf(&inst).Elem()
	//panic(fmt.Sprintf("%v",elem.CanAddr()))
	//panic(elem.MethodByName("Hello"))
	//panic(elem.NumMethod())

	//vv:=reflect.TypeOf(elem)

	//vv.Method(0).Func.Type().In(0).Name()

	//panic(elem.Method(0).Type().In(0).Name())

	rt := otpl.New("F:\\workspace\\otpl-node\\test\\.otpl")
	err := rt.Render(data, "/case", os.Stdout) //develop
	if err != nil {
		//log.Debug(err)
		panic(err)
	}

}
