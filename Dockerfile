FROM golang:1.17.0-bullseye

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go build -o tg-sota-feedback ./main.go

CMD ["./tg-srach"]