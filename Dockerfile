FROM golang:1.18-alpine as builder
WORKDIR /go/src/github.com/cecobask/apiserver
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/apiserver cmd/apiserver/main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go/bin/apiserver .
EXPOSE 8080 8080
CMD ["/apiserver"]