FROM golang:1.21-bullseye AS dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "./cmd/access_token/main.go"]
