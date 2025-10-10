# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror/blob/main/domain.txt)中查看

# 快速上手
## 安装
```shell
kubectl create -k yaml
```
## 卸载
```shell
kubectl delete -k yaml
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
 - 支持amd64和arm64架构的镜像替换



# License

Apache-2.0

# 特别感谢

- [DaoCloud](https://github.com/DaoCloud)免费提供的镜像代理服务