FROM golang:1.25.3

WORKDIR /app

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/cmd ./cmd
COPY backend/internal ./internal

RUN go build -o server ./cmd/server

EXPOSE 8199

CMD ["/app/server"]
