FROM golang:1.15
RUN mkdir -p /agones-golang-wasm/internal/server
WORKDIR /agones-golang-wasm/internal/server
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /agones-golang-wasm/internal/server/server .
CMD ["./server"]