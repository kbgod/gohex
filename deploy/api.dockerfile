# Step 1: Modules
FROM golang:1.25.4-alpine3.22 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM modules AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o svc ./cmd/api/main.go

# Step 3: Final
FROM alpine:3.22

EXPOSE 8080

COPY --from=builder /build/svc .
CMD [ "./svc"]