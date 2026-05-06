FROM jakes/base-image:alpine-3.14.0-tz
LABEL maintainers="Jakes Lee"
LABEL description="iKuai exporter"
ARG TARGETPLATFORM

COPY $TARGETPLATFORM/ikuai-exporter /app

EXPOSE 9090
WORKDIR /data

CMD ["/app", "server"]
