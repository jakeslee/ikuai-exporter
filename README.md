# iKuai Exporter
![GitHub Tag](https://img.shields.io/github/v/tag/jakeslee/ikuai-exporter?logo=github&label=release)

一个用于获取采集爱快路由的统计数据，并导出为 Prometheus 格式的 Exporter。

### 部署

部署下面的容器，配置容器环境变量，设置爱快地址和登录密码。

```shell
docker pull ghcr.io/jakeslee/ikuai-exporter:latest
# or
docker pull docker.io/jakes/ikuai-exporter:latest
```

### 使用

登录的帐号密码建议创建一个只读用户使用。

```bash
# ikuai-exporter server -h
Run metrics endpoint

Usage:
  ikuai-exporter server [flags]

Flags:
  -h, --help              help for server
      --insecure-skip     Skip iKuai certificate verification (default true)
  -l, --level string      Log level (default "info")
  -p, --password string   The password for the user on iKuai (default "test123")
      --url string        iKuai URL (default "http://10.0.1.253")
  -u, --username string   iKuai username (default "test")
```

从 v0.2.1 开始，可以使用环境变量来设置上面的参数，格式为 `IKUAI_XXX`，如 `IKUAI_URL=http://10.0.1.253` 或 `IKUAI_USERNAME=test`。

下面的方式依然支持，但**将在以后版本中弃用**。

| 变量名     | 说明     |
|:------- |:------ |
| IK_URL  | 爱快地址   |
| IK_USER | 爱快登录用户 |
| IK_PWD  | 爱快登录密码 |
