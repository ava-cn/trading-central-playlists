FROM alpine:latest

LABEL maintainer="Curder <q.curder@gmail.com>" \
  org.label-schema.name="ava-cn" \
  org.label-schema.vendor="curder" \
  org.label-schema.schema-version="1.0"

RUN apk --no-cache add ca-certificates && \
  rm -rf /var/cache/apk/*

# add tzdata options, Fixed "unknown time zone Asia/Shanghai" Error
RUN apk add -U --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# Copy the Pre-built binary file from the previous stage
ADD release/linux/amd64/trading-central-playlists /bin/

# Expose port 8087 to the outside world
EXPOSE 8087

# Command to run the executable
CMD ["/bin/trading-central-playlists"]
