# build-stage
FROM golang:1.17.2-bullseye as build-stage
USER root

ENV GO111MODULE="on"
ENV TZ=Asia/Taipei

WORKDIR /
RUN mkdir build_space
WORKDIR /build_space
COPY . .
WORKDIR /build_space/cmd/toc_trader
RUN go build

# production-stage
FROM debian:bullseye as production-stage
USER root

ENV DEPLOYMENT=docker
ENV TZ=Asia/Taipei

WORKDIR /
RUN apt update -y && \
    apt install -y tzdata && \
    apt autoremove -y && \
    apt clean && \
    mkdir toc_trader && \
    mkdir toc_trader/cmd && \
    mkdir toc_trader/cmd/toc_trader && \
    mkdir toc_trader/configs && \
    mkdir toc_trader/logs && \
    mkdir toc_trader/scripts && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /toc_trader

COPY --from=build-stage /build_space/cmd/toc_trader/toc_trader ./cmd/toc_trader/toc_trader
COPY --from=build-stage /build_space/scripts/docker-entrypoint.sh ./scripts/docker-entrypoint.sh

WORKDIR /toc_trader/cmd/toc_trader

ENTRYPOINT ["/toc_trader/scripts/docker-entrypoint.sh"]
