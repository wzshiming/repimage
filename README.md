# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror/blob/main/domain.txt)中查看

# 快速上手

## 前置要求

本项目支持两种证书管理方式：
1. **使用 Cert-Manager（推荐）** - 自动管理证书生命周期
2. **手动生成证书** - 使用 `certs/gencerts.sh` 脚本手动生成证书

### 使用 Cert-Manager（推荐）

#### 安装 Cert-Manager

如果你的集群中还没有安装 cert-manager，请先安装：

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

等待 cert-manager 的所有 Pod 都处于 Running 状态：

```shell
kubectl wait --for=condition=Ready pods --all -n cert-manager --timeout=300s
```

#### 安装 repimage（使用 Cert-Manager）

```shell
git clone https://github.com/shixinghong/repimage.git
cd repimage 
kubectl apply -f yaml/cert-manager.yaml  ## 创建证书资源
kubectl apply -f yaml/webhook.yaml       ## 安装webhook的服务
kubectl apply -f yaml/admission.yaml     ## 安装admission配置
```

cert-manager 会自动生成和管理 TLS 证书，并注入到 MutatingWebhookConfiguration 中。

### 使用手动生成的证书

如果你不想使用 cert-manager，可以使用手动生成的证书：

```shell
git clone https://github.com/shixinghong/repimage.git
cd repimage 
# 生成证书
cd certs && ./gencerts.sh && cd ..
# 需要手动将生成的 CA 证书编码后填入 yaml/admission-manual-cert.yaml 的 caBundle 字段
kubectl apply -f yaml/webhook-manual-cert.yaml  ## 安装webhook的服务
kubectl apply -f yaml/admission-manual-cert.yaml ## 安装admission配置
```

**注意：** 使用手动证书时，需要使用 `webhook-manual-cert.yaml` 和 `admission-manual-cert.yaml` 文件。

## 卸载

### 卸载 repimage（使用 Cert-Manager）

```shell
kubectl delete -f yaml/admission.yaml
kubectl delete -f yaml/webhook.yaml
kubectl delete -f yaml/cert-manager.yaml
```

### 卸载 repimage（使用手动证书）

```shell
kubectl delete -f yaml/admission-manual-cert.yaml -f yaml/webhook-manual-cert.yaml
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

# 配置选项

## 环境变量

- `TLS_CERT_PATH`: TLS 证书文件路径（默认：`/etc/webhook/certs/tls.crt`，cert-manager 生成的路径）
- `TLS_KEY_PATH`: TLS 私钥文件路径（默认：`/etc/webhook/certs/tls.key`，cert-manager 生成的路径）

如果 cert-manager 路径下的证书不存在，系统会自动回退到 `./certs/serverCert.pem` 和 `./certs/serverKey.pem`（手动生成的证书路径）。



# License

Apache-2.0

# 特别感谢

- [DaoCloud](https://github.com/DaoCloud)免费提供的镜像代理服务