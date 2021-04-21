预定义Values:

Values通过模板中.Values对象可访问的values.yaml文件(或者--set 参数)提供，
但可以模板中访问其他预定义的数据片段。
- Release.Name: 版本名称(非chart的)
- Release.Service: 组织版本的服务
- Release.IsUpgrade: 如果当前操作是升级或回滚，设置为true
- Release.IsInstall: 如果当前操作是安装，设置为true
- Chart: Chart.yaml的内容。因此，chart的版本可以从 Chart.Version 获得， 并且维护者在Chart.Maintainers里。
- Files: chart中的包含了非特殊文件的类图对象。这将不允许您访问模板， 但是可以访问现有的其他文件（除非被.helmignore排除在外）。 使用{{ index .Files "file.name" }}可以访问文件或者使用{{.Files.Get name }}功能。 您也可以使用{{ .Files.GetBytes }}作为[]byte方位文件内容。
- Capabilities: 包含了Kubernetes版本信息的类图对象。({{ .Capabilities.KubeVersion }}) 和支持的Kubernetes API 版本({{ .Capabilities.APIVersions.Has "batch/v1" }})



# values.yaml文件

```yaml
title: "My WordPress Site" # Sent to the WordPress template

global:
  app: MyWordPress

mysql:
  max_connections: 100 # Sent to MySQL
  password: "secret"

apache:
  port: 8080 # Passed to Apache
```

-----------------------
更高阶的chart可以访问下面定义的所有变量。因此WordPress chart可以用.Values.mysql.password访问MySQL密码。 
但是低阶的chart不能访问父级chart，所以MySQL无法访问title属性。同样也无法方位apache.port

Values 被限制在命名空间中，但是命名空间被删减了。因此对于WordPress chart， 它可以用.Values.mysql.password访问MySQL的密码字段。
但是对于MySQL chart，值的范围被缩减了且命名空间前缀被移除了， 因此它把密码字段简单地看作.Values.password。

全局Values: [global]
    这个值以.Values.global.app在 所有 chart中有效。
    mysql模板可以以{{.Values.global.app}}访问app，同样apachechart也可以访问。 实际上，上面的values文件会重新生成为这样：
子chart的global无法影响父chart的值
并且，父chart的全局变量优先于子chart中的全局变量。
```yaml
title: "My WordPress Site" # Sent to the WordPress template

global:
  app: MyWordPress

mysql:
  global:
    app: MyWordPress
  max_connections: 100 # Sent to MySQL
  password: "secret"

apache:
  global:
    app: MyWordPress
  port: 8080 # Passed to Apache
```
--------------------
# 架构文件

```json
{
  "$schema": "https://json-schema.org/draft-07/schema#",
  "properties": {
    "image": {
      "description": "Container Image",
      "properties": {
        "repo": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "name": {
      "description": "Service name",
      "type": "string"
    },
    "port": {
      "description": "Port",
      "minimum": 0,
      "type": "integer"
    },
    "protocol": {
      "type": "string"
    }
  },
  "required": [
    "protocol",
    "port"
  ],
  "title": "Values",
  "type": "object"
}
```

# 用户自定义资源(CRD)

CRD YAML文件应被放置在chart的crds/目录中。 多个CRD(用YAML的开始和结束符分隔)可以被放置在同一个文件中。Helm会尝试加载CRD目录中 所有的 文件到Kubernetes。

CRD 文件 无法模板化，必须是普通的YAML文档。

当Helm安装新chart时，会上传CRD，暂停安装直到CRD可以被API服务使用，然后启动模板引擎， 渲染chart其他部分，并上传到Kubernetes。
```yaml
# crontab.yaml文件必须包含没有模板指令的CRD:

kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  group: stable.example.com
  versions:
    - name: v1
      served: true
      storage: true
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
```
Helm在安装templates/内容之前会保证CronTab类型安装成功并对Kubernetes API可用。

您也可以使用helm 帮您找到chart的格式或信息的问题：
```cmd
helm lint chart
```


helm pull bitnami/wordpress --untar



# 模板语法
函数:
  quote 
  repeat N str 重复str N次
  upper 大写
  default 定义默认的值
  lookup:
    lookup 函数可以用于在运行的集群中 查找 资源。
    命令	                                    Lookup 函数
    kubectl get pod mypod -n mynamespace	    lookup "v1" "Pod" "mynamespace" "mypod"  // v1 --> apiVersion
    kubectl get pods -n mynamespace	          lookup "v1" "Pod" "mynamespace" ""
    kubectl get pods --all-namespaces	        lookup "v1" "Pod" "" ""
    kubectl get namespace mynamespace	        lookup "v1" "Namespace" "" "mynamespace"
    kubectl get namespaces	                  lookup "v1" "Namespace" "" ""    

  运算符:
    eq
    ne 
    lt
    gt
    and
    or...

Logic and Flow Control Functions:
and 
  and .Arg1 .Arg2
  返回bool

or
  ...
not
  bool求反
  not .Arg
eq
  bool等式
  eq .Arg1 .Arg2
ne -> !=
lt -> <
le -> <=
gt
ge
default  
  default "foo" .Bar
  整形: 0
  字符串: ""
  列表: []
  字典: {}
  布尔: false
  以及所有的nil (或 null)
empty
  如果给定的值被认为是空的，则empty函数返回true，否则返回false。
  empty .Arg
fail
  无条件地返回带有指定文本的空 string 或者 error。这在其他条件已经确定而模板渲染应该失败的情况下很有用。
  fail "Please accept the end user license agreement"
coalesce
  coalesce函数获取一个列表并返回第一个非空值。
  coalesce 0 1 2
  返回 1

  coalesce .name .parent.name "Matt"
  依次检查 .name .parent.name是否为空

ternary
  ternary函数获取两个值和一个test值。如果test值是true，则返回第一个值。如果test值是空，则返回第二个值。
  ternary "foo" "bar" true |   true | ternary "foo" "bar"
  返回 foo


printf
  printf "%s has %d dogs." .Name .NumberDogs

trim
  trim行数移除字符串两边的空格
  trim "   hello    "

trimAll
  移除给定的字符
  trimAll "$" "$5.00"
  --> 5.00

trimPrefix
  移除前缀
  trimPrefix "-" "-hello"
  --> hello

trimSuffix
  移除后缀
  trimSuffix "-" "hello-"
  --> hello

lower

upper

title
  首字母大写

untitle

repeat

substr
  获取字符串的字串
  start (int)
  end (int)
  string (string)
  substr 0 5 "hello world"

nospace
  删除字符串所有空格

trunc
  截断字符串
  trunc 5 "hello world"
  --> "hello "
  trunc -5 "hello world"
  --> world

contains
  测试字符串是否包含在另一个字符串中：
  contains "cat" "catch"
- --> true 因为 catch 包含了 cat.

hasPrefix

hasSuffix

quote

squote
  单引号

cat 
  将多个字符串合并成一个，用空格分隔：

indent
   以指定长度缩进给定字符串所在行
   indent 4 $lots_of_text 
   缩进四个空格

nindent
  和indent一样，并且在开头添加新行

replace
  "I Am Henry VIII" | replace " " "-"
   I-Am-Henry-VIII

snakecase
  将驼峰写法转换成蛇形写法。
  snakecase "FirstName"
  --> first_name

atoi
  字符串转换成整形。

toString

toJson

toPrettyJson

regexMatch
  regexMatch "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$" "test@acme.com"
  如果输入字符串包含可匹配正则表达式任意字符串，则返回true。

date
  date函数格式化日期

  日期格式化为YEAR-MONTH-DAY：

  now | date "2006-01-02"



键值存储数据类型

dict

  $myDict := dict "name1" "value1" "name2" "value2" "name3" "value 3"

  get:
    get $myDict "key1"

  set:
    $_ := set $myDict "name4" "value4"

  unset
    $_ := unset $myDict "name4"

  haskey:


list
  
  整形列表:
    $myList := list 1 2 3 4 5




seq
  原理和bash的 seq 命令类似
  seq 5       => 1 2 3 4 5
  seq -3      => 1 0 -1 -2 -3
  seq 0 2     => 0 1 2
  seq 2 -2    => 2 1 0 -1 -2
  seq 0 2 10  => 0 2 4 6 8 10
  seq 0 -2 -5 => 0 -2 -4

add
  add 1 2 3

add1


sub 
  减法

div
  整除

mod
  取模

mul
  相乘
  mul 1 2 3

max


min


len 
  以整数返回参数的长度

uuidv4

  生成一个uuid