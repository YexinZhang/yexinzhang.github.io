# ConfigMap

#### Create a ConfigMap

You can use either kubectl create configmap or ConfigMap generator in kustomizaton.yaml to create a configmap

```
kubectl create configmap <map-name> <data-source>
<data-source> --> file, directory, literal value draw the data from.
当我们以文件为data-source时，<data-source>中的键默认为文件的基本名称，而值默认为文件内容。
```

Create a configmap from directories

```
kubectl create configmap game-config --from-file=/path/to/dir

mkdir -p configure-pod-container/configmap/

# Download the sample files into `configure-pod-container/configmap/` directory
wget https://kubernetes.io/examples/configmap/game.properties -O configure-pod-container/configmap/game.properties
wget https://kubernetes.io/examples/configmap/ui.properties -O configure-pod-container/configmap/ui.properties

kubectl create configmap game-config --from-file=configure-pod-container/configmap/


```
The above command packages each files,
可以查看生成的configmap-->
```
kubectl describe configmaps game-config

Name:         game-config // configmap的名称
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data //configmap的内容
====
game.properties: // 文件名
---- // 从这里开始是文件的内容
enemies=aliens
lives=3
enemies.cheat=true
enemies.cheat.level=noGoodRotten
secret.code.passphrase=UUDDLRLRBABAS
secret.code.allowed=true
secret.code.lives=30
ui.properties: // 文件名
----
color.good=purple
color.bad=yellow
allow.textmode=true
how.nice.to.look=fairlyNice
```

kubectl 命令行创建的configmap 等于下面 -->
``` yaml
apiVersion: v1
kind: ConfigMap
metadata:
  creationTimestamp: 2016-02-18T18:52:05Z
  name: game-config
  namespace: default
  resourceVersion: "516"
  uid: b4952dc3-d670-11e5-8cd0-68f728db1985
data:
  game.properties: |
    enemies=aliens
    lives=3
    enemies.cheat=true
    enemies.cheat.level=noGoodRotten
    secret.code.passphrase=UUDDLRLRBABAS
    secret.code.allowed=true
    secret.code.lives=30
  ui.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true
    how.nice.to.look=fairlyNice
```


#### Create ConfigMaps from files

```
kubectl create configmap file-config --from-file=configure-pod-container/configmap/game.properties
```

从文件中创建，只会含有一个文件和内容
不过可以执行很多次，并且--from-file如果指向别的文件的话， 送有的内容都会追加进去


#### Use the option --from-env-file to create a ConfigMap from an env-file, for example:
```
Env-file 包含了很多的键值对，开头为"#"的行，空白的行将会被忽略
examples:

cat configure-pod-container/configmap/game-env-file.properties
enemies=aliens
lives=3
allowed="true"

kubectl create configmap game-config-env-file \
       --from-env-file=configure-pod-container/configmap/game-env-file.properties
```

我们查看configmap 的详细信息
```yaml
data:
  color: purple
  how: fairlyNice
  textmode: "true"
```

#### Define the key to use when creating a ConfigMap from a file 

我们可以在创建configmap的时候修改data字段列表的名称
```
kubectl create configmap config-map-1 --from-file=<key-name>=</path/to/file>
```
example: -->
```
kubectl create configmap my-config-map --from-file=my-key=/config/xxx.yaml

kubectl get configmaps my-config-map -o yaml -->
...
data:
  my-key: |
    enemies=aliens
    lives=3
    enemies.cheat=true
    enemies.cheat.level=noGoodRotten
    secret.code.passphrase=UUDDLRLRBABAS
    secret.code.allowed=true
    secret.code.lives=30
```

#### Create configmap from literal values.

We can use --from-literal define the literal values from command line:
```
kubectl create configmap my-config --from-literal=<key>=<value>

kubectl create configmap my-config --from-literal=special.value=camellia --from-literal=special.color=pink

kubectl get ... -o yaml
...
data:
  special.value: camellia
  special.color: pink
```

#### Create a ConfigMap from Generator

    kubectl support kustomization.yaml since 1.14. We can alse create a configmap from generators and then apply
it to create object on the Apiserver.
    kustomization.yaml should be specified in a directory with many files.
example: -->
```text
cat server/config.json
{
    "name":"camellia"
}

cat game/blue
color=blue

cat game/pink
color=pink

cat game/white
color=white
--------
利用kubectl apply -k . 创建configmap:
# cat kustomization.yaml
configMapGenerator:
- name: wocao-haha
  files:
  - server/config.json
- name: dingni-xixi
  files:
  - game/blue
  - game/pink
  - game/white

Then:
kubectl apply -k .

output --> /path/to/config-file created
```
