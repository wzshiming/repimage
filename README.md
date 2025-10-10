# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror/blob/main/domain.txt)中查看

# 快速上手
## 安装
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

# 配置选项

repimage 现在支持多种方式配置域名替换列表（allowlist）：

## 环境变量配置

- `ALLOWLIST_FILE`: 指定包含域名映射的文件路径（默认：`/etc/repimage/allowlist.txt`）
- `ALLOWLIST_URL`: 指定获取域名映射的URL（默认：DaoCloud公共镜像仓库）
- `ALLOWLIST_UPDATE_INTERVAL`: 设置更新间隔，如 `30m`, `1h`, `24h`（默认：`1h`，设为`0`禁用更新）

详细配置说明请参考：[Allowlist配置指南](docs/ALLOWLIST_CONFIG.md)

## 特性

- ✅ **默认内置映射列表**：应用包含常用镜像仓库的默认映射
- ✅ **构建时获取**：Docker镜像构建时自动下载最新映射列表
- ✅ **参数化配置**：支持通过环境变量和配置文件自定义映射
- ✅ **定期更新**：支持定期从URL更新映射列表
- ✅ **多级回退**：文件 → URL → 默认内置列表的加载优先级

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