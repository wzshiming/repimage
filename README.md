# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror/blob/main/domain.txt)中查看

# 快速上手
## 安装

### 使用预构建的镜像
Docker 镜像已发布到 GitHub Container Registry (ghcr.io):
```shell
docker pull ghcr.io/wzshiming/repimage:latest
```

### 从源码安装
```shell
git clone https://github.com/shixinghong/repimage.git
cd repimage 
kubectl apply -f  yaml/webhook.yaml ## 一定要先安装webhook的服务 ready之后再安装admission
kubectl apply -f  yaml/admission.yaml
```
## 卸载
```shell
kubectl delete -f  yaml/webhook.yaml -f yaml/admission.yaml
```

# 使用后效果
自动替换yaml文件中的镜像地址，例如: 
```
k8s.gcr.io/coredns/coredns => m.daocloud.io/k8s.gcr.io/coredns/coredns

nginx => m.daocloud.io/docker.io/library/nginx
```
# 注意事项：
 - 只有在 [public-image-mirror
   ](https://github.com/DaoCloud/public-image-mirror/blob/main/domain.txt)中的地址才会被替换，否则不会被替换
 - 替换的方式是**增加前缀**方式，不是**替换**方式
 - 目前只支持在amd64架构下的镜像替换，如果需要可以自行编译打包是使用



# License

Apache-2.0

# 特别感谢

- [DaoCloud](https://github.com/DaoCloud)免费提供的镜像代理服务