FROM golang:1.8

RUN ls
ADD ./ /go/
RUN go run main.go
