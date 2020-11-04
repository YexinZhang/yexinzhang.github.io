## Pod Topology Spread Constraints

### <font face='consolas' face=3> Spread constraints for Pods</font>

``` yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  topologySpreadConstraints:
  - maxSkew: <int>
    topologyKey: <string>
    whenUnsatisfiable: <string>
    labelSelector: <object>
```

<font face='微软雅黑' face=2 >
您可以定义一个或多个topologySpreadConstraint，以指示kube调度程序如何将每个传入的Pod相对于整个集群中的现有Pod进行放置。
</font>

<font face='华文中宋' size=3>字段:  </font>
```
maxSkew: 
    两个给定topo之间的pod数量的最大差值, 值必须大于零.
    他还和 whenUnsatisfiable 的值有关联:
        DoNotSchedule: maxSkew是目标拓扑中的匹配Pod数与全局最小值之间的最大允许差值。
        ScheduleAnyway: scheduler gives higher precedence(优先权) to topologies that would help reduce the skew.
topologyKey:
    key of the node labels, 如果节点拥有相同的key和值，这些个节点会被视做在一个topology
    The scheduler tries to place a balanced number of pods into each topology domain
whenUnsatisfiable:
    indicates how to deal with a Pod if it doesn't satisfy the spread constraint:
    DoNotSchedule(default):  tells the scheduler not to schedule it.
    ScheduleAnyway: 告诉调度程序在优先级最小化偏斜的节点时仍要调度它。
labelSelector:
    used to find matching Pods.
```
example:
node-1 zoneA
node-2 zoneA
node-4 zoneA
node-3 zoneB

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 50
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      topologySpreadConstraints:
      - maxSkew: 2 //最大的偏斜
        topologyKey: zone // 被调度到的node 的label
        whenUnsatisfiable: DoNotSchedule
        labelSelector:
          matchLabels:
            app: nginx
      containers:
      - name: nginx
        image: qa-harbor.leihuo.netease.com/library/nginx:v2
        ports:
        - containerPort: 80
```
```
这个例子中的topo按照zone划分，zoneA中有node-1,2,4三台机器， zoneB有node-4一台机器
node-1 node-2 node-3 总共会被分配25个label为app=nginx的Pod， 25个平均分配到3台机器上.
node-4 将会被分到25个

集群中的node都标注了zone的label， pod创建的结果就是，会在zone节点上生成，未标注zone的节点不会被调度到.

topoSpreadConstraints 可以出现多个， 他们之间是“与”的关系
```
约定:
  - 只有与传入Pod具有相同名称空间的Pod才可以匹配候选。
  - 没有topologyKey存在的节点将会被绕过
  - topologySpreadConstraints[*].labelSelector 一般要与workerload自己的标签相互匹配
  - 如果传入的Pod定义了spec.nodeSelector或spec.affinity.nodeAffinity，则不匹配它们的节点将被绕过。