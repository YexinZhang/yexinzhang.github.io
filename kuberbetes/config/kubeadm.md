[Basic]

A kubeadm config file could contain multiple configuration types separated using "---"

kubeadm support following configuration types
``` yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration

apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration

apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration

apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration

apiVersion: kubeadm.k8s.io/v1beta2
kind: JoinConfiguration
```
To print the defaults for "init" and "join" actions using the following commands:
```
kubeadm config print init-defaults
kubeadm config print join-defaults
```

initConfiguration Type:
    should be used to configure runtime settings.
    - nodeRegistration: hold fields that relate  registering the new node to the cluster;
        use it to customize the node name, the CRI socket and all setting which are specific to the node
        where kubeadm is executed, including.(e.g. the node IP)
        包含与将新节点注册到集群有关的字段；使用它来自定义节点名称，要使用的CRI套接字或仅应应用于该节点的任何其他设置（例如，节点ip）
    - localAPIEndpoint:
        
