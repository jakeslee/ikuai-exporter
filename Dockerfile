

FROM jakes/base-image:alpine-3.14.0-tz
LABEL maintainers="Jakes Lee"
LABEL description="iKuai exporter"
ARG TARGETPLATFORM

ADD ./output /output
RUN cp /output/"$TARGETPLATFORM"/app /app

EXPOSE 9090
WORKDIR /data

RUN chmod +x /app
CMD ["/app"]
