# Chart

Chart 文件结构

``` t
wordpress/
  Chart.yaml          # 包含了chart信息的YAML文件
  LICENSE             # 可选: 包含chart许可证的纯文本文件
  README.md           # 可选: 可读的README文件
  values.yaml         # chart 默认的配置值
  values.schema.json  # 可选: 一个使用JSON结构的values.yaml文件
  charts/             # 包含chart依赖的其他chart
  crds/               # 自定义资源的定义
  templates/          # 模板目录， 当和values 结合时，可生成有效的Kubernetes manifest文件
  templates/NOTES.txt # 可选: 包含简要使用说明的纯文本文件
```
## Chart的配置文件
### Chart.yaml
```yaml
apiVersion: chart API 版本(必需)
name: chart名称 （必需）
version: 语义化2 版本（必需）
kubeVersion: 兼容Kubernetes版本的语义化版本（可选）
# >= 1.13.0 < 1.15.0 
# >= 1.13.0 < 1.14.0 || >= 1.14.1 < 1.15.0
description: 一句话对这个项目的描述（可选）
type: chart类型 （可选）
keywords:
  - 关于项目的一组关键字（可选）
home: 项目home页面的URL （可选）
sources:
  - 项目源码的URL列表（可选)
  - name: chart名称 (nginx)
    version: chart版本 ("1.2.3")
    repository: （可选）仓库URL ("https://example.com/charts") 或别名 ("@repo-name")
    condition: （可选） 解析为布尔值的yaml路径，用于启用/禁用chart (e.g. subchart1.enabled )
    tags: # （可选）
      - 用于一次启用/禁用 一组chart的tag
    import-values: # （可选）
      - ImportValue 保存源值到导入父键的映射。每项可以是字符串或者一对子/父列表项
    alias: （可选） chart中使用的别名。当你要多次添加相同的chart时会很有用
maintainers: # （可选）
  - name: 维护者名字 （每个维护者都需要）
    email: 维护者邮箱 （每个维护者可选）
    url: 维护者URL （每个维护者可选）
icon: 用做icon的SVG或PNG图片URL （可选）
appVersion: 包含的应用版本（可选）。不需要是语义化，建议使用引号
deprecated: 不被推荐的chart （可选，布尔值）
annotations:
  example: 按名称输入的批注列表 （可选）.
```
### dependencies
```
一旦你定义好了依赖，运行 helm dependency update 就会使用你的依赖文件下载所有你指定的chart到你的charts/目录。
Condition - 条件字段field 包含一个或多个YAML路径（用逗号分隔）。 如果这个路径在上层values中已存在并解析为布尔值，chart会基于布尔值启用或禁用chart。
只会使用列表中找到的第一个有效路径，如果路径为未找到则条件无效。

```

### Tags
```yaml
# Tags - tag字段是与chart关联的YAML格式的标签列表。在顶层value中，通过指定tag和布尔值，可以启用或禁用所有的带tag的chart。

# parentchart/Chart.yaml

dependencies:
  - name: subchart1
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart1.enabled, global.subchart1.enabled
    tags:
      - front-end
      - subchart1
  - name: subchart2
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart2.enabled,global.subchart2.enabled
    tags:
      - back-end
      - subchart2


# parentchart/values.yaml

subchart1:
  enabled: true
tags:
  front-end: false
  back-end: true

在上面的例子中，
所有带 front-endtag的chart都会被禁用，但只要上层的value中 subchart1.enabled 路径被设置为 'true'，该条件会覆盖 front-end标签且 subchart1 会被启用。
一旦 subchart2使用了back-end标签并被设置为了 true，subchart2就会被启用。 
也要注意尽管subchart2 指定了一个条件字段， 但是上层value没有相应的路径和value，因此这个条件不会生效。
--set 参数一如既往可以用来设置标签和条件值。

helm install --set tags.front-end=true --set subchart2.enabled=false

#[导入子数据] 
下面示例中的import-values 指示Helm去拿到能再child:路径中找到的任何值，并拷贝到parent:的指定路径。

# parent's Chart.yaml file

dependencies:
  - name: subchart1
    repository: http://localhost:10191
    version: 0.1.0
    ...
    import-values:
      - child: default.data
        parent: myimports

# parent's values.yaml file

myimports:
  myint: 0
  mybool: false
  mystring: "helm rocks!"

# subchart1's values.yaml file

default:
  data:
    myint: 999
    mybool: true

父chart的结果值将会是这样：
# parent's final values

myimports:
  myint: 999
  mybool: true
  mystring: "helm rocks!"
```