ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG BUILD_REVISION
ARG BUILD_USER
ARG BUILD_DATE
ARG BUILD_BRANCH
ARG VERSION

LABEL BUILD_REVISION="$BUILD_REVISION" \
      BUILD_USER="$BUILD_USER" \
      BUILD_DATE="$BUILD_DATE" \
      BUILD_BRANCH="$BUILD_BRANCH" \
      VERSION="$VERSION"

WORKDIR /
ADD cmd/operator/bin/operator .

# On busybox 'nobody' has uid `65534'
USER 65534

ENTRYPOINT ["/operator"]
