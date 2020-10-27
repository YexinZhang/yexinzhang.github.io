## 如何使用docker的私有仓库

### 方法一： 利用docker login
docker login HUB 会在家目录产生一个配置文件 .docker/config.json
我们只需要将这个配置文件复制到所有工作节点的/var/lib/kubelet/ 目录中去，保持原名
就可以直接在yaml文件中使用私有仓库的镜像了
```
1. Run docker login [server]
2. Get a list of our nodes
    nodes=$( kubectl get nodes -o jsonpath='{range.items[*].metadata}{.name} {end}' )
3. Use scp CMD copy these config files to remote nodes.
```


### 方法二: 准备imagePullSecret
    This is the recommanded approach to run containers based on the private registries
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: qa-harbor
  namespace: default
type: docker-registry
stringData:
  docker-server: "qa-harbor.leihuo.netease.com"
  docker-username: "admin"
  docker-password: "Harbor12345"
```

### 方法三：更具docker login产生的config.yaml 准备secret
```
kubectl create secret generic regcred \
    --from-file=.dockerconfigjson=<path/to/.docker/config.json> \
    --type=kubernetes.io/dockerconfigjson
```

定制 secret
先要将config.json用base64编码
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myregistrykey
  namespace: awesomeapps
data:
  .dockerconfigjson: UmVhbGx5IHJlYWxseSByZWVlZWVlZWVlZWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGx5eXl5eXl5eXl5eXl5eXl5eXl5eSBsbGxsbGxsbGxsbGxsbG9vb29vb29vb29vb29vb29vb29vb29vb29vb25ubm5ubm5ubm5ubm5ubm5ubm5ubm5ubmdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2cgYXV0aCBrZXlzCg==
type: kubernetes.io/dockerconfigjson
```


使用时:
``` yaml
apiVersion: v1
kind: Pod
metadata:
  name: private-reg
spec:
  containers:
  - name: private-reg-container
    image: <your-private-image>
  imagePullSecrets:
  - name: regcred
```