模板

{{ Action }}

{{-  }} // 去掉左边的所有的空格
{{  -}} // 去掉右边

"{{ 12 -}} < {{- 13 }}" --> 12<13
"{{ 12 }} < {{ 13 }}"   --> 12 < 13

Actions:
```go
{{/* a comment */}}

{{ pipeline }}
	The default textual representation 
	// 默认文字表示

{{if pipeline}} T1 {{end}}
	// 如果管道的值为空，则不生成任何输出

{{if pipeline}} T1 {{else}} T0 {{end}}
	// If the value of the pipeline is empty, T0 is executed;
	// otherwise, T1 is executed. Dot is unaffected.

{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}


{{range pipeline}} T1 {{end}}
	// The value of the pipeline must be an array, slice, map, or channel.
	// If the value of the pipeline has length zero, nothing is output;

{{range pipeline}} T1 {{else}} T0 {{end}}


{{template "name"}}
	// 指定名称的模板将使用nil数据执行。

{{template "name" pipeline}}

{{with pipeline}} T1 {{end}}
	// 将dot设置为管道的值，并且T1为 被执行。

{{with pipeline}} T1 {{else}} T0 {{end}}
```

Arguments:

参数是一个简单的值，由以下之一表示。
```go

```

eq
	Returns the boolean truth of arg1 == arg2
ne
	Returns the boolean truth of arg1 != arg2
lt
	Returns the boolean truth of arg1 < arg2
le
	Returns the boolean truth of arg1 <= arg2
gt
	Returns the boolean truth of arg1 > arg2
ge
	Returns the boolean truth of arg1 >= arg2