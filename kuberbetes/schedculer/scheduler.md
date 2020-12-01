kube-scheduler考虑的因素
- 单个或者集体的资源需求
- 硬件，软件，policy的约束
- 亲和性和反亲和性
- 数据局部性
- 进程之间的干扰
....

步骤
- Filtering
    PodFitsResources filter检查候选节点资源是否符合Pod配置的资源需求;
    After PodFitsResources， 通常会得到一些节点，如果为空，pod将不会被调度
- Scoring
    打分
最终， 会挑选一个分数最高的进行调度


### 有两种支持的方式来配置调度程序的筛选和评分行为：
1. 通过调度Sceduling Policies，您可以配置谓词进行过滤，配置优先级进行评分
2. 调度概要文件使您可以配置实现不同调度阶段的插件，
   包括：队列排序，过滤器，得分，绑定，保留，许可等。您还可以将kube-scheduler配置为运行不同的配置文件。

### 污点/容忍度:
taint
```
$ kubectl taint NODE NAME KEY_1=VAL_1:TAINT_EFFECT_1 ... KEY_N=VAL_N:TAINT_EFFECT_N
```
打污点:
    kubectl taint node node01 key1=value1:NoSchedule
    kubectl taint nodes foo bar:NoSchedule
    kubectl taint node -l myLabel=X  dedicated=foo:PreferNoSchedule

    kubectl taint node node01 key1:value1:NoSchedule- //取消污点
    kubectl taint nodes foo dedicated- //删除全部dedicated
    kubectl taint nodes foo dedicated:NoSchedule- //删除全部的dedicated,并且Taint_EFFECT为Noscheduler

容忍度:
    可以在podSpec中定义容忍度，可以使pod调度到被打了对应污点的node上
```yaml
tolerations:
- key: "key1"
  operator: "Equal"
  value: "value1"
  effect: "NoSchedule"

tolerations:
- key: "key1"
  operator: "Exists" //只要有key1， 不考虑value的值
  effect: "NoSchedule"

tolerations:
- key: "key1"
  operator: "Equal"
  value: "value1"
  effect: "NoExecute" 
  tolerationSeconds: > 
    6000 //如果没有配置，那么容忍此污点的pod将永久绑定在该node， 
    如果指定了秒数, 则表示还可以在node上存活的时间长度， 超时仍将被驱逐
```

如果key为空字符串，并且operator为Exists, 则表示容忍一切
effect为空，表示匹配所有的effect

effect:
```
PreferNoSchedule: 尽可能的不调度到该节点
NoExecute: 不仅不会调度，还会驱逐Node上已有的Pod。
    tolerationSeconds: 6000
```

You can put multiple taints on the same node and multiple tolerations on the same pod.