FROM golang:1.23.1-bookworm AS builder 

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/main.go

FROM alpine AS production 
WORKDIR /app
COPY --from=builder /app/server .


EXPOSE 8080
CMD [ "./server" ]