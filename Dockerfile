FROM golang:alpine
WORKDIR /build
COPY . .
RUN go build -o backend ./cmd/shortener
CMD ["./backend"]
