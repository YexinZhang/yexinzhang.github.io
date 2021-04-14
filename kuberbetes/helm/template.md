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