FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

EXPOSE 7540

CMD ["/my_app"]
