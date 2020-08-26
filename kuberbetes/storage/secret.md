# Secret

secret 本质上其实就是configmap, 只是将一些value加密了

create
```
kubectl create secret { generic | docker-registry | tls } <name> --from-file=<key>=<value> ...

```