# Build Step
FROM golang:1.20-alpine AS builder

# Dependencies
RUN apk update && apk add --no-cache upx make git

# Source
WORKDIR $GOPATH/src/github.com/Depado/scarecrow
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp
RUN upx --best --lzma /tmp/scarecrow


# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/scarecrow /go/bin/scarecrow
ENTRYPOINT ["/go/bin/scarecrow"]
