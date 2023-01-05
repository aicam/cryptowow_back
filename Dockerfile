FROM golang:1.16-alpine

WORKDIR /go/src/github.com/aicam/cryptowow_back

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY ./ .

RUN go build -o ./build github.com/aicam/cryptowow_back

RUN ls

EXPOSE 4300

CMD [ "/go/src/github.com/aicam/cryptowow_back/build" ]

