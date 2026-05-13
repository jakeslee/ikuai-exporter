# iKuai Exporter
![GitHub Release](https://img.shields.io/github/v/release/jakeslee/ikuai-exporter?include_prereleases)

一个用于获取采集爱快路由的统计数据，并导出为 Prometheus 格式的 Exporter。

### 部署

部署下面的容器，配置容器环境变量，设置爱快地址和登录密码。

```shell
docker pull ghcr.io/jakeslee/ikuai-exporter:latest
# or
docker pull docker.io/jakes/ikuai-exporter:latest
```

使用 docker-compose 部署：

```yaml
services:
    ikuai-exporter:
        image: ghcr.io/jakeslee/ikuai-exporter:latest
        restart: always
        environment:
            IKUAI_URL: "http://10.0.1.253"
            IKUAI_USERNAME: "test"
            IKUAI_PASSWORD: "test123"
        ports:
            - "9090:9090"
```

部署完成后，访问 `http://IP:9090/metrics` 验证运行情况。

接下来将 exporter 的采集地址 IP 配置到 Prometheus 的 `scrape_configs` 中就可开始使用。

详细配置和 Grafana 配置示例可以参考使用[样例](https://blog.imoe.tech/2022/12/25/48-use-ikuai-exporter-to-gather-metrics/)，最新的演示 Dashboard 在[这里](https://github.com/jakeslee/ikuai-exporter/raw/refs/heads/master/examples/grafana-dashboard.json)。

### 参数说明

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
