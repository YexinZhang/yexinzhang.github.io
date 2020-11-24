先创建pv， 在创建引用pv的pvc，
pvc可以是静态的，他必须匹配相同大小的pv，
pvc也可以是动态的，但是要使用storageclass，直接申请pvc，会让storageclass自动创建出指定大小的pv,

using
pod在挂载时，挂载的是pvc

storage object in Use Protection
删除已经挂载到正在运行的pod中的pvc不会被立即删除

[_Reclaiming]
PV回收策略:
    Retain
        当pvc被删除时，pv依旧存在，该卷被认为为released， 但是不能被后面的pvc引用，
        删除:
            1. 删除PersistentVolume。删除PV之后，外部基础架构中的关联存储资产（例如AWS EBS，GCE PD，Azure Disk或Cinder卷）仍然存在。
            2. 相应地手动清理关联存储资产上的数据。
            3. 手动删除关联的存储资产，或者如果要重复使用相同的存储资产，请使用存储资产定义创建一个新的PersistentVolume。
    Delete
        对于支持``删除回收''策略的卷插件，删除操作会同时从Kubernetes中删除PersistentVolume对象以及外部基础架构中的关联存储资产，
        例如AWS EBS，GCE PD，Azure Disk或Cinder卷。
        动态预配置的卷继承其StorageClass的回收策略，该策略默认为Delete。

    Recycle
        不建议使用“回收回收”策略。相反，推荐的方法是使用动态配置。


[binding]
pvc绑定指定的pv
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: foo-pvc
  namespace: foo
spec:
  storageClassName: "" # Empty string must be explicitly set otherwise default StorageClass will be set
  volumeName: foo-pv
```


[pv]
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-new-01
spec:
  capacity:
    storage: 10Gi # G是基于10的乘法， Gi是基于二的乘法  1G == 1000M ; 1Gi == 1024Mix`
  volumeMode: Filesystem  # Filesystem(default) Block
  accessModes:
    - ReadWriteOnce
    # RWO ReadWriteOnce -- can be mounted as read-write by a single node
    # RWX ReadWriteMany -- can be mounted read-only by many nodes
    # ROX ReadOnlyMany  -- can be mounted as read-write by many nodes
  persistentVolumeReclaimPolicy: Retain # 回收
  storageClassName: slow # 特定类别的PV只能绑定到请求该类别的PVC
  mountOptions:
    - hard
    - nfsvers=4.1
  nfs:
    path: /tmp
    server: 172.17.0.2
```

[Node_Affinity]
    You need to explicitly set this for local volume.
    PV可以指定node affinity来定义约束，这些约束限制了此volume可以被访问的node.
    使用PV的Pod仅会安排到由节点亲缘关系选择的节点上

[Phase]
A volume will be in one of the following phases.
- Available -- 尚未绑定到claim
- Bound  --  绑定到了claim
- Released  --  claim被删除了，但是居群资源尚未回收
- Failed  --  自动回收失败



[pvc]
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myclaim
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 8Gi
  storageClassName: slow
  selector:
    matchLabels:
      releasae: "stable"
    matchExpressions:
      -- {key: environment, operator: In, values: [dev]}
```