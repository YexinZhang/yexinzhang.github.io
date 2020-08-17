#NetworkPolicy

Mandatory Fields:
apiVersion, metadata, kind
  namespace: default
```yaml
spec:
podSelector:
  matchLabels: 
  // 搜索到对应的pod之后，会对其网络流量进行隔离
policyTypes:
  - Ingress
  - Egress
ingress:
```
```
每个NetworkPolicy都可以包括允许的ingress规则的列表，
每一个规则都会放行匹配到from和port规则的流量
```
``` yaml
- from:
  - ipBlock
  - namespaceSelector
  - podSelector:
      matchLabels:
        xx: xx
  ports:
  - protocol: TCP
    port: ##
egress
#  egress 主要匹配前往的目的地
- to
  - ipBlock...
  ...
  ports:
  - ...
```
###Example：
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-network-policy
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: db
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - ipBlock:
        cidr: 172.17.0.0/16
        except:
        - 172.17.1.0/24
    - namespaceSelector:
        matchLabels:
          project: myproject
    - podSelector:
        matchLabels:
          role: frontend
    ports:
    - protocol: TCP
      port: 6379
  egress:
  - to:
    - ipBlock:
        cidr: 10.0.0.0/24
    ports:
    - protocol: TCP
      port: 5978d
```

1. 隔离 role=db 的 namespace为default的pod的网络流量。(if they weren't already isolated)
2. (Ingress rule)
   allow connections to all pod in the "default" namespace with the label 'role: db' on tcp 6379 port:
   - any pod in namespace default with the label role: db
   - any pod in a namespace with label "myproject"
   - ip cird and except cidr
3. 允许从“默认”名称空间中带有标签“ role = db”的任何Pod连接到TCP端口5978上的CIDR 10.0.0.0/24允许从“默认”名称空间中带有标签“ role = db”的任何Pod连接到TCP端口5978上的CIDR 10.0.0.0/24
    
----
ingress 和 egress 和from | to有以下几个筛选字段
podSelector
namespaceSelector
podSelector and namespaceSelector
```yaml
ingress:
- from:
  - namespaceSelector:
      matchLabels:
        user: Tom
    podSelector:
      matchLabels:
        role: db
这个就要求两者同时满足
```
```yaml
ingress:
- from:
  - namespaceSelector:
      matchLabels:
        user: Tom
  - podSelector:
      matchLabels:
        role: db

这个就是一种或的关系，两者满足其中一种就可以了
```
---
###ipBlock

should be cluster-external IP.


###Default deny all ingress policy
This policy does not change the default egress isolation behavior
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Ingress

```

###Default allow all ingress policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  ingress:
  - {}
```

###Default deny all egress policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Egress
```

###Default allow all egress policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Egress
  egress:
  - {}
```


###Default allow all ingress egress policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Egress
  - Ingress
  egress:
  - {}
  ingress:
  - {}
```

###Default deny all ingress egress policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy1
spec:
  podSelector: {}
  policyTypes:
  - Egress
  - Ingress
```