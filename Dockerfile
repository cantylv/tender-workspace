FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl

ARG MODULE_NAME=backend

COPY . /home/${MODULE_NAME}/

WORKDIR /home/${MODULE_NAME}/

RUN go build -o main cmd/main/main.go

# Service
FROM alpine:latest as production
ARG BUILDER_MODULE_NAME=backend
WORKDIR /root/

COPY --from=builder /home/${BUILDER_MODULE_NAME}/config/config.yaml config/config.yaml
COPY --from=builder /home/${BUILDER_MODULE_NAME}/main .

RUN chown root:root main
CMD ["./main"]