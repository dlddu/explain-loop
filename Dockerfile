# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24
ARG DEBIAN_CODENAME=bookworm

FROM golang:${GO_VERSION}-${DEBIAN_CODENAME} AS builder
WORKDIR /app

ENV CGO_ENABLED=0 GOTOOLCHAIN=local

# 외부 의존성이 생기면 go.sum 도 함께 복사한다.
COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /out/explain-loop .

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

ENV LISTEN_ADDR=0.0.0.0:8080
EXPOSE 8080
# runAsNonRoot 는 숫자 UID 여야 검증된다(distroless nonroot = 65532).
USER 65532:65532
COPY --from=builder /out/explain-loop /usr/local/bin/explain-loop
ENTRYPOINT ["/usr/local/bin/explain-loop"]
