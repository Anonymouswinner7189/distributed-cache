FROM golang:1.25-alpine 

WORKDIR /app

COPY . .

#Downloading dependencies
RUN go mod tidy

#building binaries
RUN go build -o node cmd/node/main.go
RUN go build -o router cmd/router/main.go

CMD ["./node"]