FROM golang

WORKDIR /app

COPY go.* .

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /imService

EXPOSE 8080

CMD ["go", "run", "main.go"]