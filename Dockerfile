FROM golang:1.20.3

WORKDIR /app
ADD . /app

RUN go mod download
RUN go env -w GO111MODULE=on
RUN go build -o bin/server main.go
# EXPOSE 8080

CMD [ "./bin/server" ] 
