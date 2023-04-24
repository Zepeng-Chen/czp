FROM golang:1.19-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o taurus cmd/main.go

EXPOSE 8080

CMD [ "/app/taurus" ]