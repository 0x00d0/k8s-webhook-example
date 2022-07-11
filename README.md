# 参考文档

`https://github.com/kubernetes/kubernetes/tree/release-1.21/test/images/agnhost/webhook`
`https://kubernetes.io/zh/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/`
`https://erosb.github.io/post/json-patch-vs-merge-patch/
`


# 生成证书

```bash
mkdir -p /opt/kubernetes/{bin,ssl}
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
chmod +x cfssl-certinfo_linux-amd64 cfssljson_linux-amd64 cfssl_linux-amd64
mv cfssl_linux-amd64 /opt/kubernetes/bin/cfssl
mv cfssljson_linux-amd64 /opt/kubernetes/bin/cfssljson
mv cfssl-certinfo_linux-amd64 /opt/kubernetes/bin/cfssl-certinfo
ln -s /opt/kubernetes/bin/cfssl /usr/local/bin/cfssl
ln -s /opt/kubernetes/bin/cfssljson /usr/local/bin/cfssljson
ln -s /opt/kubernetes/bin/cfssl-certinfo /usr/local/bin/cfssl-certinfo
```

```bash
vi ca-config.json  
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "usages": ["signing"],
        "expiry": "8760h"
      }
    }
  }
}

vi ca-csr.json  
{
  "CN": "Kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "zh",
      "L": "bj",
      "O": "bj",
      "OU": "CA"
   }
  ]
}

cfssl gencert -initca ca-csr.json | cfssljson -bare ca

```

```bash
vi server-csr.json
{
  "CN": "admission",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "zh",
      "L": "bj",
      "O": "bj",
      "OU": "bj"
    }
  ]
}

# 签发证书
# service name 设置为自己的 这里使用pod-example

cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=pod-example.kube-system.svc \
  -profile=server \
  server-csr.json | cfssljson -bare server

```

# 生成webhook配置

```bash
admissionregistration_config.yaml
caBundle的内容这么取
  cat ca.pem | base64

创建密文
kubectl create secret tls pod-example-tls --cert=server.pem --key=server-key.pem  -n kube-system

```


# WebHook部署到集群

```bash
kubectl apply -f deploy.yaml
kubectl apply -f admissionregistration_config.yaml

```





