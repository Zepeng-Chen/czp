FROM golang:1.19-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/app/*.go ./

RUN go build -o taurus

EXPOSE 8080

CMD [ "taurus" ]