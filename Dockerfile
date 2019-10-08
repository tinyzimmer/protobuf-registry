###############
# API Builder #
###############

FROM golang:1.13-alpine as apibuilder

# Build deps
RUN apk add --update curl git upx curl unzip autoconf automake libtool make g++ file

WORKDIR /workspace

# Do go deps first
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go code
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
  go build \
    -a \
    -o app \
    -ldflags "-X 'main.CompileDate=`date -u`'" \
    cmd/main.go
RUN upx app

##############
# UI Builder #
##############

FROM node as uibuilder

WORKDIR /workspace

# Do deps first
COPY ui/package.json /workspace/ui/package.json
COPY ui/package-lock.json /workspace/ui/package-lock.json
RUN cd ui && npm install

# Do actual asset build
COPY ui/ /workspace/ui
RUN cd ui && npm run build

###############
# Final image #
###############

FROM alpine

# Add protobuf utilities
RUN apk add --update protobuf protobuf-dev

# setup a user and data directories
RUN adduser -u 1000 -h /opt/proto-registry -s /bin/false -D protoregistry && \
  mkdir -p /data && \
  chown -R protoregistry: /data

USER protoregistry

WORKDIR /opt/proto-registry

# Copy API assets
COPY --from=apibuilder /workspace/app /opt/proto-registry/app

# Copy UI assets
COPY --from=uibuilder /workspace/ui/build /opt/proto-registry/static

CMD ["/opt/proto-registry/app"]
