FROM golang:1.22.2-alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8000
CMD ["./main"]
