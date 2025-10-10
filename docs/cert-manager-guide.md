# Cert-Manager Integration Guide

本文档详细说明了如何使用 Cert-Manager 为 repimage webhook 自动生成和管理 TLS 证书。

## 证书管理方式对比

### 传统手动方式 vs Cert-Manager

```
┌─────────────────────────────────────────────────────────────────┐
│                     传统手动证书方式                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. 运行 gencerts.sh 脚本生成证书                                │
│     └─> 生成 CA 和服务器证书                                     │
│                                                                   │
│  2. 手动 base64 编码 CA 证书                                     │
│     └─> 粘贴到 admission.yaml 的 caBundle 字段                  │
│                                                                   │
│  3. 证书到期后需手动重复以上步骤                                 │
│                                                                   │
│  ❌ 需要手动操作                                                 │
│  ❌ 容易出错                                                      │
│  ❌ 证书更新复杂                                                  │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                   Cert-Manager 自动化方式                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. 部署 Cert-Manager 资源                                       │
│     └─> kubectl apply -f yaml/cert-manager.yaml                 │
│                                                                   │
│  2. Cert-Manager 自动完成所有工作                                │
│     ├─> 自动生成 CA 证书                                         │
│     ├─> 自动生成服务器证书                                       │
│     ├─> 自动注入 CA 到 MutatingWebhookConfiguration            │
│     └─> 自动续期证书（到期前 30 天）                             │
│                                                                   │
│  ✅ 完全自动化                                                    │
│  ✅ 零手动操作                                                    │
│  ✅ 自动续期                                                      │
│  ✅ 生产级可靠性                                                  │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## 架构说明

### Cert-Manager 证书链

```
┌─────────────────────────────────────────────────────────────┐
│                    Cert-Manager 证书架构                      │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Self-Signed ClusterIssuer (repimage-selfsigned-issuer)   │
│     │                                                         │
│     ├──> 2. CA Certificate (repimage-ca)                     │
│          │   - Secret: repimage-ca-secret                    │
│          │                                                    │
│          ├──> 3. CA Issuer (repimage-ca-issuer)              │
│               │                                               │
│               ├──> 4. Server Certificate (repimage-webhook-cert) │
│                    - Secret: repimage-webhook-tls            │
│                    - DNS: repimage.kube-system.svc           │
│                    - Validity: 1 year                        │
│                    - Auto-renew: 30 days before expiry       │
│                                                               │
│  5. CA Bundle Auto-Injection                                 │
│     - Cert-Manager injects CA to MutatingWebhookConfiguration │
│     - Annotation: cert-manager.io/inject-ca-from             │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

1. **Self-Signed ClusterIssuer** (`repimage-selfsigned-issuer`)
   - 用于生成根 CA 证书的自签名颁发者
   - 作用域：集群级别

2. **CA Certificate** (`repimage-ca`)
   - 自签名的根 CA 证书
   - 存储在 Secret: `repimage-ca-secret`
   - 用于签发服务器证书

3. **CA Issuer** (`repimage-ca-issuer`)
   - 使用 CA 证书签发其他证书的颁发者
   - 作用域：kube-system 命名空间

4. **Server Certificate** (`repimage-webhook-cert`)
   - webhook 服务器的 TLS 证书
   - 存储在 Secret: `repimage-webhook-tls`
   - DNS 名称：
     - `repimage.kube-system.svc`
     - `repimage.kube-system.svc.cluster.local`
   - 有效期：1年
   - 自动续期：到期前 30 天

### CA Bundle 自动注入

Cert-Manager 通过以下注解自动将 CA 证书注入到 MutatingWebhookConfiguration 中：

```yaml
metadata:
  annotations:
    cert-manager.io/inject-ca-from: kube-system/repimage-webhook-cert
```

这消除了手动编码和配置 caBundle 的需要。

## 部署步骤

### 1. 安装 Cert-Manager

```bash
# 安装 cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# 验证安装
kubectl wait --for=condition=Ready pods --all -n cert-manager --timeout=300s
kubectl get pods -n cert-manager
```

### 2. 部署 Repimage

```bash
# 创建证书资源
kubectl apply -f yaml/cert-manager.yaml

# 等待证书生成（通常需要几秒钟）
kubectl wait --for=condition=Ready certificate/repimage-ca -n kube-system --timeout=60s
kubectl wait --for=condition=Ready certificate/repimage-webhook-cert -n kube-system --timeout=60s

# 部署 webhook 服务
kubectl apply -f yaml/webhook.yaml

# 部署 admission 配置
kubectl apply -f yaml/admission.yaml
```

### 3. 验证部署

```bash
# 检查证书状态
kubectl get certificate -n kube-system
kubectl describe certificate repimage-webhook-cert -n kube-system

# 检查 Secret
kubectl get secret repimage-webhook-tls -n kube-system
kubectl get secret repimage-ca-secret -n kube-system

# 检查 webhook 配置中的 CA Bundle
kubectl get mutatingwebhookconfiguration repimage -o yaml | grep -A 5 caBundle

# 检查 webhook pod 日志
kubectl logs -n kube-system -l app=repimage
```

## 证书续期

Cert-Manager 会自动处理证书续期：

- **续期时间**：证书到期前 30 天开始续期
- **证书有效期**：1 年
- **自动更新**：无需人工干预

查看证书续期状态：

```bash
kubectl describe certificate repimage-webhook-cert -n kube-system
```

## 故障排查

### 证书未生成

```bash
# 查看 Certificate 资源状态
kubectl describe certificate repimage-webhook-cert -n kube-system

# 查看 cert-manager 日志
kubectl logs -n cert-manager -l app=cert-manager
```

### CA Bundle 未注入

确保：
1. Cert-Manager CA Injector 正在运行
2. MutatingWebhookConfiguration 中有正确的注解
3. Certificate 资源已成功创建

```bash
# 检查 CA Injector
kubectl get pods -n cert-manager -l app=cainjector

# 检查注解
kubectl get mutatingwebhookconfiguration repimage -o yaml | grep inject-ca-from
```

### Webhook 无法启动

检查 Pod 日志：

```bash
kubectl logs -n kube-system -l app=repimage
```

常见问题：
- 证书 Secret 未挂载：检查 Deployment 的 volumeMounts
- 证书路径错误：查看日志中的 "Using TLS cert" 消息

## 手动证书方式（不推荐）

如果不想使用 Cert-Manager，可以使用手动证书：

```bash
# 生成证书
cd certs && ./gencerts.sh && cd ..

# 使用手动证书配置文件
kubectl apply -f yaml/webhook-manual-cert.yaml
kubectl apply -f yaml/admission-manual-cert.yaml
```

## 环境变量配置

可以通过环境变量自定义证书路径：

```yaml
env:
  - name: TLS_CERT_PATH
    value: /custom/path/tls.crt
  - name: TLS_KEY_PATH
    value: /custom/path/tls.key
```

默认路径：
- Cert-Manager: `/etc/webhook/certs/tls.crt` 和 `/etc/webhook/certs/tls.key`
- 手动证书: `./certs/serverCert.pem` 和 `./certs/serverKey.pem`

## 安全建议

1. **定期轮换 CA 证书**：虽然自动续期，但建议定期更新 CA
2. **限制 CA Secret 访问**：使用 RBAC 限制对 `repimage-ca-secret` 的访问
3. **监控证书到期**：设置告警监控证书续期状态
4. **使用命名空间隔离**：CA Issuer 仅在 kube-system 命名空间有效

## 参考资源

- [Cert-Manager 官方文档](https://cert-manager.io/docs/)
- [Kubernetes Admission Webhooks](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- [Cert-Manager CA Injector](https://cert-manager.io/docs/concepts/ca-injector/)
