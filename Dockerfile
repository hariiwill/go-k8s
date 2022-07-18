### Build application 
FROM golang:1.19rc2-alpine3.16 as builder
LABEL maintainer="Hariharan"
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

### Run application from scratch
FROM  alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8081
CMD ["./main"]