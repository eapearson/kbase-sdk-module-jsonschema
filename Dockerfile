# ------------------------------
# The build image
# ------------------------------
FROM golang:1.16.3-alpine3.13 as builder

# add deps for building it
RUN apk upgrade --update-cache --available && \
    apk add --update --no-cache bash make && \
    mkdir -p /build

COPY . /build
WORKDIR /build
RUN go build -o server .
LABEL stage=intermediate

# ------------------------------
# The product image
# ------------------------------
FROM alpine:3.13

# These ARGs values are passed in via the docker build command
ARG BUILD_DATE
ARG VCS_REF
ARG BRANCH=develop
ARG TAG

RUN mkdir -p /kb/module
COPY --from=builder /build/server /kb/module
COPY --from=builder /build/schemas /kb/module/schemas

# The BUILD_DATE value seem to bust the docker cache when the timestamp changes, move to
# the end
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.vcs-url="https://github.com/kbase/kbase-sdk-module-JSONSchema.git" \
    org.label-schema.vcs-ref=$COMMIT \
    org.label-schema.schema-version="1.0.0-rc1" \
    us.kbase.vcs-branch=$BRANCH  \
    us.kbase.vcs-tag=$TAG \
    maintainer="Steve Chan sychan@lbl.gov"

# Run as a regular user, not root.
RUN addgroup --system kbmodule && \
    adduser --system --ingroup kbmodule kbmodule && \
    chown -R kbmodule:kbmodule /kb

WORKDIR /kb/module
ENV SCHEMA_ROOT=/kb/module/schemas
ENV PORT=5000

CMD [ "./server" ]
