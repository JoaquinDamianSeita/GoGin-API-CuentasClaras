FROM golang:1.20-alpine
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /app/main .
EXPOSE 8080
ENTRYPOINT [ "/app/main" ]