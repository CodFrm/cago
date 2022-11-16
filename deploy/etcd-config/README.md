# etcd-config
k8s helm 部署资源,使用etcd+etcdkeeper作为cago应用的配置中心

单节点的etcd

etcd设置密码
```bash
kubectl exec -it etcd-config-7694d65cbb-xx2zp -c etcd -n core bash

etcdctl user add root
```
