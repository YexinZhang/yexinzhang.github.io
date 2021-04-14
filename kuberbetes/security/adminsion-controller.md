LimitRange 资源可以在名称空间上限制单个容器，Pod，PVC相关的资源用量。

ResourceQuota:
    定义名称空间级别的资源配额，从而在名称空间上限制聚合资源消耗的边界，
    它支持以资源类型来限制用户可在本地名称空间中创建的相关资源的对象数量，以及这些对象可消耗的计算资源总量等。    
    而同名的ResourceQuota准入控制器负责观察传入的请求, 属于验证类型的准入控制器

PodSecurityPolicy:
    PSP资源就是集群全局范围内定义的Pod资源可用的安全上下文策略。
    同名的PodSecurityPolicy准入控制器负责观察集群范围内的Pod资源的运行属性，并确保它没有违反PodSecurityPolicy资源定义的约束条件。


