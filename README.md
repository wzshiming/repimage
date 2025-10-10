# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror)中查看

# 快速上手
## 安装

### 方式一: 使用 Cert-Manager（推荐）
如果你的集群中已经安装了 [cert-manager](https://cert-manager.io/)，推荐使用这种方式，cert-manager 会自动管理证书的生成和更新。

1. 确保 cert-manager 已经安装在集群中：
```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

2. 安装 repimage：
```shell
kubectl apply -k https://github.com/wzshiming/repimage/yaml
```

### 方式二: 使用预生成证书
如果你的集群中没有 cert-manager，可以使用预生成的证书：
```shell
kubectl apply -k https://github.com/wzshiming/repimage/yaml/static-certs
```

或者使用预编译的单文件（推荐用于快速测试）：
```shell
kubectl create -f https://files.m.daocloud.io/github.com/wzshiming/repimage/releases/download/latest/repimage.yaml
```

# 使用后效果
自动替换yaml文件中的镜像地址，例如: 
```
k8s.gcr.io/coredns/coredns => m.daocloud.io/k8s.gcr.io/coredns/coredns

nginx => m.daocloud.io/docker.io/library/nginx
```

# 文档

- [Cert-Manager 集成文档](docs/cert-manager.md) - 如何使用 cert-manager 自动管理 TLS 证书

# License

Apache-2.0

# 特别感谢

- [@shixinghong](https://github.com/shixinghong) 感谢原作者提供的灵感
- [DaoCloud](https://github.com/DaoCloud) 免费提供的镜像代理服务
