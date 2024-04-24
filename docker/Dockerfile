FROM golang:1.22.2 AS go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY api /app/api/
COPY cmd /app/cmd/
COPY internal /app/internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o kinko ./cmd/main/main.go

FROM node:20.12.2-alpine AS svelte-builder
WORKDIR /app

COPY ./frontend /app

RUN mv .env.example .env

RUN npm install
RUN npm run build

RUN cp ./frontend/package.json ./frontend/package-lock.json /app/build/
RUN cd /app/frontend.build && npm ci --omit=dev

FROM node:20.12.2-alpine
WORKDIR /app

RUN apk --no-cache add ca-certificates wget \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk \
    && apk add --force-overwrite glibc-2.34-r0.apk

COPY --from=go-builder /app/kinko /app/api/
COPY --from=svelte-builder /app/build /app/frontend/

COPY ./docker/entrypoint.sh /app

CMD ["/bin/sh", "entrypoint.sh"]
