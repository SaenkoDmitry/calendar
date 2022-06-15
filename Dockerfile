FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

EXPOSE 8080

RUN go build -o bin/calendar ./cmd

CMD [ "/app/bin/calendar" ]
