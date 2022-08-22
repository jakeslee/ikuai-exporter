# IKuai Prometheus Exporter

### 部署

部署下面的容器，配置容器环境变量，设置爱快地址和登录密码。

```shell
docker pull jakes/ikuai-exporter:latest
```

登录的帐号密码建议创建一个只读用户使用。

| 变量名     | 说明     |
|:------- |:------ |
| IK_URL  | 爱快地址   |
| IK_USER | 爱快登录用户 |
| IK_PWD  | 爱快登录密码 |
