#!/bin/sh

version=$(git describe --long --dirty --abbrev=6 --tags)
flags="-X main.buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.version=$version"

export CGO_ENABLED=0 GOOS=linux

echo "start to build version:$version"

mkdir -p ./output/linux/

for i in "arm64" "amd64" ; do
  echo "building for $i..."
  GOARCH="$i" go build -ldflags "$flags" -o ./output/linux/$i/app main.go
  chmod +x ./output/linux/$i/app
done

image=ccr.ccs.tencentyun.com/imoe-tech/go-playground:ikuai-exporter-"$version"
echo "packaging docker multiplatform image: $image"

docker buildx build --push \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION="$version" \
  -t "$image" .

echo "finished: $image"

