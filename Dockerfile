FROM alpine:latest
MAINTAINER MessageBird <support@messagebird.com>

# Copy over the binary in the container.
COPY bin/pushprom-*.linux-amd64 /usr/bin/pushprom

EXPOSE 9090 9091

# Run
CMD ["/usr/bin/pushprom"]


