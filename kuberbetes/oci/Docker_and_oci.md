Docker and OCI runtime

What is The Open Container Initiative
    The Open Container Initiative (OCI) is a Linux Foundation project to design open standards for containers.
    OCI currently contains two specifications: The runtime specification and image specification
    OCI runtime spec defines how to run the OCI image bundle as a container
    OCI image spec defines how to create an OCI image, which include an image manifest,a filesystem(layer) serialization, and image configuration

Container runtimes
    [containerd]
        A CNCF project.
        它管理其主机系统的完整容器生命周期，包括映像管理，存储和容器生命周期，监督，执行和联网
    [lxc]
    [runc]
        runc是一个CLI工具，用于根据OCI规范生成和运行容器。
    [cri-o]
    [rkt]


What is docker then?
    在v1.11.0之前，docker engine 作为一个整体来管理容器，包括镜像管理，生命周期管理，产生容器，资源限制，网络等

    在v1.11.0之后, 
    见 OneNote

containerd
    containerd service
    依赖:
        runc: to run container
        ctr: a cli for container
        containerd-shim: to support daemonless containers
    
    containerd can do:
        1. manage images( like download images from dockerhub )
        2. Manage containers(create and run)
        3. Manage namespaces
    
    interact with containerd:
    ctr client:
        ctr images pull docker.io/library/redis:latest
        ctr images list
        ctr run -d docker.io/library/redis:latest NAME
        ctr containers list
        

