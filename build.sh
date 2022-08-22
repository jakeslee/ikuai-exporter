#/bin/sh

VERSION=`git describe --long --dirty --abbrev=6 --tags`

docker build -t ccr.ccs.tencentyun.com/imoe-tech/go-playground:ikuai-exporter-v$VERSION . &&\
docker push ccr.ccs.tencentyun.com/imoe-tech/go-playground:ikuai-exporter-v$VERSION
