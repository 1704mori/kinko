FROM golang:1.22.2 AS go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY api /app/api/
COPY cmd /app/cmd/
COPY internal /app/internal/

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o kinko ./cmd/main/main.go

FROM node:20.12.2-alpine AS svelte-builder
WORKDIR /app

COPY ./frontend /app

RUN mv .env.example .env

RUN npm install
RUN npm run build

FROM node:20.12.2-alpine
WORKDIR /app

COPY --from=go-builder /app/kinko /app/api/
COPY --from=svelte-builder /app/build /app/frontend/

COPY ./docker/entrypoint.sh /app

CMD ["/bin/sh", "entrypoint.sh"]
