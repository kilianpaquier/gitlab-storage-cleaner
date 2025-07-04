# Code generated by craft; DO NOT EDIT.

#############################
#        STAGE BUILD        #
#############################
FROM golang:1.24.4 AS build

ARG CGO_ENABLED=0
ARG GIT_REF_NAME
ARG GIT_COMMIT
ARG VERSION=v0.0.0

WORKDIR /app

COPY . .

RUN go mod download && \
    CGO_ENABLED=$CGO_ENABLED go build \
        -ldflags "\
            -X 'github.com/kilianpaquier/gitlab-storage-cleaner/internal/build.branch=$GIT_REF_NAME' \
            -X 'github.com/kilianpaquier/gitlab-storage-cleaner/internal/build.commit=$GIT_COMMIT' \
            -X 'github.com/kilianpaquier/gitlab-storage-cleaner/internal/build.date=$(TZ="UTC" date '+%Y-%m-%dT%TZ')' \
            -X 'github.com/kilianpaquier/gitlab-storage-cleaner/internal/build.version=$VERSION' \
        " \
        -o gitlab-storage-cleaner ./cmd/gitlab-storage-cleaner

#############################
#         STAGE RUN         #
#############################
FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.authors="kilianpaquier"
LABEL org.opencontainers.image.vendor="kilianpaquier"

LABEL org.opencontainers.image.title="gitlab-storage-cleaner"
LABEL org.opencontainers.image.description="Easily clean gitlab maintained repositories storage (jobs artifacts only) with a simple command"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="github.com/kilianpaquier/gitlab-storage-cleaner"
LABEL org.opencontainers.image.source="github.com/kilianpaquier/gitlab-storage-cleaner"
LABEL org.opencontainers.image.documentation="github.com/kilianpaquier/gitlab-storage-cleaner"

WORKDIR /app

COPY --from=build \
    /app/gitlab-storage-cleaner \
    ./

EXPOSE 3000

ENTRYPOINT [ "/app/gitlab-storage-cleaner" ]
