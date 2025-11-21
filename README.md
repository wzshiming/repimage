# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror)中查看

# 快速上手
## 安装
```shell
kubectl create -f https://files.m.daocloud.io/github.com/wzshiming/repimage/releases/download/latest/repimage.yaml
kubectl rollout status deployment/repimage -n kube-system
```

# 使用后效果
自动替换yaml文件中的镜像地址，例如: 
```
k8s.gcr.io/coredns/coredns => m.daocloud.io/k8s.gcr.io/coredns/coredns

nginx => m.daocloud.io/docker.io/library/nginx
```

# 配置选项
## 忽略指定域名
如果你有一些私有镜像仓库或者不需要加速的域名，可以通过 `--ignore-domains` 参数来忽略这些域名。

例如，在 deployment.yaml 中添加参数：
```yaml
containers:
- command:
  - /repimage
  - --ignore-domains=myregistry.example.com,private.registry.local
```

这样，来自 `myregistry.example.com` 和 `private.registry.local` 的镜像将不会被替换。

## 自定义镜像前缀
默认使用 `m.daocloud.io` 作为镜像前缀，可以通过 `--prefix` 参数自定义：
```yaml
containers:
- command:
  - /repimage
  - --prefix=mirror.example.com
```

建议内网再部署一级缓存, 可以使用 `--prefix=你的内网地址/mirror.example.com`

# License

Apache-2.0

# 特别感谢

- [@shixinghong](https://github.com/shixinghong) 感谢原作者提供的灵感
- [DaoCloud](https://github.com/DaoCloud) 免费提供的镜像代理服务
