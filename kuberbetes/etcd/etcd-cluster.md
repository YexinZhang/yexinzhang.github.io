部署一个 TLS 的etcd集群

一旦etcd集群开始运作，我们可以通过 runtime reconfiguration。

Static 方式的部署
``` yaml
Machine:
  - Name: daffodil
    IP: 192.168.65.103
    Hostname: daffodil.zyxasr.io
  - Name: pink
    IP: 192.168.65.102
    Hostname: pink.zyxasr.io
  - Name: orange
    IP: 192.168.65.101
    Hostname: orange.zyxasr.io
```

我们可以使用以下环境变量或者参数来指定集群的成员
ETCD_INITIAL_CLUSTER="daffodil=http://192.168.65.103:2380,pink=http://192.168.65.102:2380,orange=http://192.168.65.101:2380"
ETCD_INITIAL_CLUSTER_STATE=new

COMMAND:
--initial-cluster daffodil=http://192.168.65.103:2380,pink=http://192.168.65.102:2380,orange=http://192.168.65.101:2380
--initial-cluster-state new

initial-cluster 用来集群内部的etcd实例相互之间通信，成为advertised peer URLS. 

--initial-advertise-peer-urls 节点绑定的peer通信url

--initial-cluster-token 集群令牌， 用来区分不同的集群，相当于一个集群名称,通过这样做，
                    etcd可以为集群生成唯一的集群ID和成员ID，即使它们的配置完全相同

--listen-client-urls  接受客户端流量的地址
--advertise-client-urls  暴漏给客户端的地址

``` 简单的非TLS集群配置
etcd 
    --name daffodil \
    --initial-cluster-token etcd-cluster-1 \

    --initial-advertise-peer-urls http://192.168.65.103:2380 \  暴露我们的peer通信地址
    --listen-peer-urls http://192.168.65.103:2380 \  指定我们要监听的peer地址
    --initial-cluster http://192.168.65.103:2380,pink=http://192.168.65.102:2380,orange=http://192.168.65.101:2380 \
    // 指定我们的集群信息

    --listen-client-urls http://192.168.65.103:2379,http://127.0.0.1:2379 \
    // 指明对外服务的地址
    --advertise-client-urls http://192.168.65.103:2379 \
    // 暴露对外服务的地址

    --initial-cluster-state new
```

TLS

etcd supports encrypted communication through the TLS protocol.
tls channel can be used for internal cluster communitcaion between peers
as well as encrypted client traffic.
当我们使用自己签证书时
peer端通信的证书都由同一CA颁发
CN 可以使用泛域名  比如 *.zyxasr.io
SANS: 必须要包括所有的etcd节点，比如
    DNS:localhost, 
    DNS:orange.zyxasr.io, 
    DNS:pink.zyxasr.io, 
    DNS:daffodil.zyxasr.io, 
    DNS:camellia.zyxasr.io, 
    IP Address:127.0.0.1, 
    IP Address:192.168.65.101, 
    IP Address:192.168.65.102, 
    IP Address:192.168.65.103, 
    IP Address:192.168.65.100
api server的客户端证书也需要由同样的CA签发，

TLS 配置的相关选项
etcd --name camellia \
  --data-dir /data/etcd/  \
  --initial-advertise-peer-urls https://192.168.65.100:2380 \
  --listen-peer-urls https://192.168.65.100:2380 \
  --listen-client-urls https://192.168.65.100:2379,https://127.0.0.1:2379 \
  --advertise-client-urls https://192.168.65.100:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster camellia=https://192.168.65.100:2380,orange=https://192.168.65.101:2380,pink=https://192.168.65.102:2380 \
  --initial-cluster-state new \
  // 证书配置
  --client-cert-auth \
  --trusted-ca-file /opt/etcd/certs/ca.crt \
  --cert-file=/opt/etcd/certs/server.crt --key-file=/opt/etcd/certs/server.key \
  --peer-client-cert-auth --peer-trusted-ca-file=/opt/etcd/certs/ca.crt \
  --peer-cert-file=/opt/etcd/certs/server.crt --peer-key-file=/opt/etcd/certs/server.key
  // 用于启动server的证书和节点之间相互通信的证书可以使用一套，，并且每个节点上的证书可以使用同一套