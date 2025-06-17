FROM alpine:3 AS alpine-upgrader
RUN apk upgrade --no-cache

FROM scratch AS alpine-upgraded
COPY --from=alpine-upgrader / /
CMD ["/bin/sh"]

FROM alpine-upgraded AS pkg-builder

RUN apk -U add \
    sudo \
    alpine-sdk \
    apkbuild-pypi

RUN mkdir -p /var/cache/distfiles && \
    adduser -D packager && \
    addgroup packager abuild && \
    chgrp abuild /var/cache/distfiles && \
    chmod g+w /var/cache/distfiles && \
    echo "packager ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

WORKDIR /work
USER packager

RUN abuild-keygen -a -i -n

COPY --chown=packager:packager packages/ ./

RUN cd py3-cssselect2 && \
    abuild -r && \
    cd ../py3-weasyprint && \
    abuild -r

FROM alpine-upgraded

RUN --mount=from=pkg-builder,source=/home/packager/packages/work,target=/packages \
    --mount=from=pkg-builder,source=/etc/apk/keys,target=/etc/apk/keys \
    apk add --no-cache --repository /packages \
    font-liberation \
    font-liberation-sans-narrow \
    ttf-linux-libertine \
    python3 \
    py3-weasyprint

ENV PYTHONUNBUFFERED=1
WORKDIR /app
# Copy your application code
COPY rest-weasyprint /app
EXPOSE 8080


# Set up health checks.
HEALTHCHECK CMD wget -qO- http://localhost:8080/ || exit 1

ENTRYPOINT ["/app/rest-weasyprint"]


