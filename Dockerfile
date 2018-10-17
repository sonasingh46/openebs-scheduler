FROM golang:1.10-alpine as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG VERSION=0.0.1

# build
WORKDIR /go/src/openebs-scheduler
COPY . .
RUN go install -ldflags "-s -w -X main.version=$VERSION" openebs-scheduler

# runtime image
FROM gcr.io/google_containers/ubuntu-slim:0.14
COPY --from=builder /go/bin/openebs-scheduler /usr/bin/openebs-scheduler
ENTRYPOINT ["openebs-scheduler"]
