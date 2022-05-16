FROM golang:1.16-alpine

WORKDIR /go/src/github.com/aicam/cryptowow_back

ENV MYSQLCONNECTION root:021021ali@tcp(mysqlserver:3306)/server?charset=utf8mb4&parseTime=True
ENV REDISPASS 021021ali
ENV ARENAFILEPATH /test

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY ./ .

RUN go build -o ./build github.com/aicam/cryptowow_back

RUN ls

EXPOSE 4300

CMD [ "/go/src/github.com/aicam/cryptowow_back/build" ]

