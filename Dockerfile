### First stage ###
FROM golang:alpine AS builder
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . ./
RUN go build -o /amzn-scraper .

### Second stage ###
FROM alpine
WORKDIR /app
COPY --from=builder /amzn-scraper /amzn-scraper
COPY --from=builder /app /app
ENTRYPOINT ["/amzn-scraper"]