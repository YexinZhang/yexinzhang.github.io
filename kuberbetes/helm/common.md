### 基础知识
```
chart: helm package
repository: charts can be collected and shared
release: is an instance of a chart running in a Kubernetes cluster.

Helm installs charts into Kubernetes, creating a new release for each installation. 
And to find new charts, you can search Helm chart repositories.
```

###  Command 
helm search
    helm search hub
         searches the Artifact Hub, which lists helm charts from dozens of different repositories.
    helm search repo
        searches the repositories that you have added to your local helm client (with helm repo add).

helm repo add
    helm  repo add   brigade https://brigadecore.github.io/charts
    helm search repo shti
    NAME            CHART VERSION   APP VERSION     DESCRIPTION
    brigade/kashti  0.5.0           v0.4.0          A Helm chart for Kubernetes 

