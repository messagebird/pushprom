FROM alpine:latest
LABEL maintainer="support@messagebird.com"

# Copy over the binary in the container.
ADD bin/pushprom /usr/bin/

EXPOSE 9090 9091

# Run
CMD ["/usr/bin/pushprom"]


