FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rps-middleware-gitlab .

FROM scratch

COPY --from=builder /app/rps-middleware-gitlab /rps-middleware-gitlab

ENV REMOTE_URL "https://gitlab.com"
EXPOSE 8080

ENTRYPOINT ["/rps-middleware-gitlab"]
