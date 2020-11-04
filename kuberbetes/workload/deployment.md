```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx-deploy
spec:
  selector:
    matchLabels:
      app: nginx // 这个label是让 svc 可以搜索到的, 而且要与template的label相同
  strategy:
    type: RollingUpdate //滚动发布
    rollingUpdate:
      maxSurge: 8 //在yaml定义的replicas基础上可以多出来的pod
      maxUnavailable: 2 // 最大的可以关闭的pod的数量
  replicas: 10
  template:
    metadata:
      labels:
        app: nginx
    spec:
      topologySpreadConstraints: // 根据node的label进行的topo分区
      - maxSkew: 1
        topologyKey: zone //根据哪一个label分区
        whenUnsatisfiable: DoNotSchedule
        labelSelector:
          matchLabels:
            app: nginx
      containers:
      - name: nginx
        image: qa-harbor.leihuo.netease.com/library/nginx:v2
        ports:
        - containerPort: 80
          name: http-port
```

### rollout 命令
```
rollout

history     View rollout history
  spec.revisionHistoryLimit 可以控制ReplicaSet 的保留数量, 默认为10,设置为0，deployment将无法回滚
pause       Mark the provided resource as paused
restart     Restart a resource
resume      Resume a paused resource
status      Show the status of the rollout
undo        Undo a previous rollout

Example:
kubectl rollout history deployment nginx-deployment //可以用来查看deployment的更新历史
deployment.apps/nginx-deployment
REVISION  CHANGE-CAUSE
2         <none>
9         kubectl set image deployment/nginx-deployment nginx=nginx:1.14.2 --record=true
10        kubectl set image deployment/nginx-deployment nginx=nginx:1.18.0 --record=true
11        kubectl set image deployment/nginx-deployment nginx=nginx:1.17.1 --record=true

如果 kubectl set的时候没有使用 --record 选项， 则history选项记录的行的名称会重复, 但是 内容是新的内容
查看版本细节
  kubectl rollout history deployment nginx-deployment --revision=#
如果我们想回退到某一个版本 我们可以使用undo（上一个版本） 或者undo --to-revision=REVISION指定版本号
  kubectl rollout undo deployment/nginx-deployment --revision=10
直接使用undo的话，就是在最近的两个版本之间来回切假腿

暂停版本的发布, 其实会记录新版本，只是新版本不会应用到底层的 ReplicaSet
kubectl rollout pause deployment nginx-deployment
解禁版本发布, 在暂停期间做的操作都会应用
kubectl rollout resume deployment nginx-deployment

```

### Scaling a deployment
```
kubectl scale deployment.v1.apps/nginx-deployment --replicas=10

```


```
Pod Template:
  Deployment need .apiVersion, .metadata, .kind and alse need a .spec section
  .sepc.template and .sepc.selector 是2必需的字段
  .spec.template 是一个pod template，相当于kind=Pod 所需要的字段
  pod template 还必须要有适合的labels，和适合的restart policy(必须为Always[默认值]).

Replicas:
  optional field that specials the number of desired Pods. It default to 1.

Selector:
  .spec.selector 是一个需要的字段, 
  .spec.selector must match .spec.template.metadata.labels
  
Strategy（升级战略）:
  .spec.strategy Specifies the strategy used to replace old pod by new pods
  .spec.strategy.type can be Recreate and RollingUpdate(default)
  Recreate:
    所有的现有pod都将被杀死，然后重新建立pod
  Rolling Update Deployment:
    .spec.strategy.type == RollingUpdate. 
    set "maxUnavailable" and "maxSurge" to control the rolling update processing
    maxUnavailable:
      Optional field that specifies the maximum number of pods that can be unavailable during the update process
      就是一次更新的pod数量，比如replicas=10, maxunavailable=3, 这个时候就会有三个pod去升级 无法被访问，留下7个旧的
    maxSurge:
      用于指定可以在所需数量的Pod上创建的最大Pod数
      当此值设置为30％时，可以在滚动更新开始时立即按比例放大新的ReplicaSet，以使旧Pod和新Pod的总数不超过所需Pod的130％。
      一旦旧Pod被杀死，新的ReplicaSet可以进一步扩展，以确保更新期间随时运行的Pod总数最多为所需Pod的130％

Progress Deadline Seconds:
  .spec.progressDeadlineSeconds default == 600
  If specified, this field needs to be greater than .spec.minReadySeconds.

Min Ready Second:
  新创建的Pod应该准备就绪的最短秒数，而其任何容器都不会崩溃，这才被视为可用。
  就是创建多久之后 认为它是成功的, 默认为0

Revision History Limit:
  就是 kubectl get rs 中要保留以允许回滚的旧ReplicaSet的数量
  回滚使用 kubectl rollout undo sourcestype/source
  一旦replicaSets 丢失了，就无法回滚了

Paused:
  .spec.paused
  是一个可选的布尔字段，用于暂停和恢复部署。
  是一个可选的布尔字段，用于暂停和恢复部署。暂停的Deployment和未暂停的Deployment之间的唯一区别是，
  只要暂停，对暂停的Deployment的PodTemplateSpec所做的任何更改都不会触发新的推出。创建展开时，默认情况下不会暂停该展开。
```