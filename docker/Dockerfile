FROM golang:1.22.2 AS go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o kinko .

FROM node:20.12.2-alpine AS svelte-builder
WORKDIR /app

COPY ./frontend /app

RUN npm install
RUN npm run build

FROM node:20.12.2-alpine
WORKDIR /app

COPY --from=go-builder /app/kinko /app/api/
COPY --from=svelte-builder /app/build /app/frontend/

COPY ./entrypoint.sh /app

CMD ["/bin/sh", "entrypoint.sh"]
