# syntax=docker/dockerfile:1

FROM node:20-alpine AS web-build
WORKDIR /app
RUN apk add --no-cache git
COPY web/package.json web/pnpm-lock.yaml web/package-lock.json* web/ ./web/
RUN cd web && \
    if [ -f pnpm-lock.yaml ]; then npm i -g pnpm && pnpm install; \
    else npm install; fi
COPY web/ ./web/
RUN cd web && npm run build

FROM golang:1.24.2-alpine AS go-build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN rm -rf static/out && mv web/out static/
RUN go build -o bestsub ./cmd/bestsub

FROM alpine:3.20
WORKDIR /app
COPY --from=go-build /app/bestsub ./bestsub
VOLUME ["/data"]
ENV BESTSUB_SERVER_HOST=0.0.0.0
ENV BESTSUB_SERVER_PORT=8080
EXPOSE 8080
CMD ["./bestsub", "-c", "/data/config.json"]
