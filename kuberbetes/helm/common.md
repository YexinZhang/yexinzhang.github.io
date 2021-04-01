### 基础知识
```
chart: helm package
repository: charts can be collected and shared
release: is an instance of a chart running in a Kubernetes cluster.

Helm installs charts into Kubernetes, creating a new release for each installation. 
And to find new charts, you can search Helm chart repositories.
```

###  Command 
helm search
    helm search hub
         searches the Artifact Hub, which lists helm charts from dozens of different repositories.
    helm search repo
        searches the repositories that you have added to your local helm client (with helm repo add).

helm repo add
    helm  repo add   brigade https://brigadecore.github.io/charts
    helm search repo shti
    NAME            CHART VERSION   APP VERSION     DESCRIPTION
    brigade/kashti  0.5.0           v0.4.0          A Helm chart for Kubernetes 

helm search repo xxx

helm install my-release bitnami/wordpress
helm uninstall happy-panda

helm status RELEASE_NAME

自定义Chart
helm show values bitnami/wordpress > values.yaml
helm install -f values.yaml bitnami/wordpress --generate-name

    --values | -f file.yaml
    --set 直接修改值, --set can be cleared by running helm upgrade with --reset-values specified.
        --set name=value
            name: value
        --set a=b,c=d
            a: b
            c: d
        --set outer.inner=value
            outer:
                inner: value
        --set name={a, b, c}
            - a
            - b
            - c
        --set servers[0].port=80
            servers:
                - port: 80
        --set servers[0].port=80,servers[0].host=example
            servers:
                - port: 80
                  host: example
        --set name=value1\,value2
            name: "value1,value2"
        --set nodeSelector."kubernetes\.io/role"=master
            nodeSelector:
                kubernetes.io/role: master

When a new version of a chart is released, or when you want to change the configuration of your release,
you can use the helm upgrade command.

    helm upgrade -f panda.yaml happy-panda bitnami/wordpress

pull到本地自己修改

    helm pull bitnami/wordpress

回滚
Now, if something does not go as planned during a release,
 it is easy to roll back to a previous release using helm rollback [RELEASE] [REVISION].

    helm rollback happy-panda 1

Every time an install, upgrade, or rollback happens, the revision number is incremented by 1. 
The first revision number is always 1. 

    helm history [RELEASE]
    查看历史 REVISION


删除
    helm uninstall release

    helm uninstall --keep-history. 
    If you wish to keep a deletion release record
    Using helm list --uninstalled will only show releases that were uninstalled with the --keep-history flag.

repo

    helm repo list 
    helm repo add dev https://example.com/dev-charts
    helm repo update 更新

Create my own Chart.
        
    helm create deis-workflow

    As you edit your chart, you can validate that it is well-formed by running helm lint.

        helm lint

    将自己的Chart打包
        helm package deis-workflow
    安装自己的chart
        helm install deis-workflow ./deis-workflow-0.1.0.tgz

other

     It has also covered useful utility commands like helm status, helm get, and helm repo


Tips and tricks


    包含另外一个template
        value: {{ include "mytpl" . | lower | quote }}
    
    下面的必需函数示例为.Values.who声明了一个必需项，并且在缺少该项时将显示错误消息：
        value: {{ required "A valid .Values.who entry required!" .Values.who }} 
    
    Quote
        name: {{ .Values.MyName | quote }}
        port: {{ .Values.Port }} // 整型就不要加引号了

    Using the 'include' Function
        Go提供了一种使用内置模板指令将一个模板包含在另一个模板中的方法。 但是，内置函数不能在Go模板管道中使用。
        为了能够包含一个模板，然后对该模板的输出执行操作，Helm具有一个特殊的include函数：
        {{ include "toYaml" $value | indent 2 }}
        将其传递给$ value，然后将该模板的输出传递给indent函数。