package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type stu struct {
	Name string
	Age int
	Gender string
}

func main() {

	mux:= http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			_, err := w.Write([]byte("HTTP Form Parse False"))
			if err != nil {
				log.Println(err)
				return
			}
		}
		tmpl, err := template.New("test").Parse(`
{{$one:= 1}} {{$two:=2}} {{if eq $one $two }} // 首先，“//” 这个注释不能在此使用
// golang template support `le ge lt gt eq en`几种比较的方法，
// 而且这集中方法是函数，比较的对象添加在后面
垃圾o
{{else}}
草泥马o
{{end}}
// range 根对象中的元素
// i为key, k is value
{{range $i, $k := .}}
{{$i}}:{{$k}}
{{end}}
// 可以更具参数的true or false 判断
{{if .IsZero}}
	<h1>不能为0</h1>
{{ else }}
	<h2>{{ .Result }}</h2>
{{ end }}

// with 可以将根中的对象包含起来，方便调用对象的元素，减少模板冗余
STU
{{- with .stu}}
name: {{.Name}}
age:  {{.Age}}
gender: {{.Gender}}
{{- end}}
`)
		if err != nil {
			fmt.Println(err)
			return
		}
		//parseint --> 将字符串转换为数字类型, 10,8, 16进制, 加上8,32,64位宽
		x, _:= strconv.ParseInt(r.URL.Query().Get("x"),10,8)
		y, _:= strconv.ParseInt(r.URL.Query().Get("y"),10,8)

		IsZero := y == 0
		result := 0.0
		if !IsZero {
			result = float64(x) / float64(y)
		}
		// map[string]interface{} 可以作为模板的传入参数，模板绘解析键值对
		err = tmpl.Execute(w, map[string]interface{}{
			"IsZero":IsZero,
			"Result":result,
			"Range":[]string{
				"caonima",
				"wocaonima",
				"caonimab",
				"rinima",
			},
			"stu":&stu{
				Name:   "jack",
				Age:    19,
				Gender: "Male",
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	err := http.ListenAndServe(":80",mux)
	if err != nil {
		fmt.Println(err)
	}

}
