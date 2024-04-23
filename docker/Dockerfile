# Build Golang application
FROM golang:1.22.2-alpine AS go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY api /app/api/
COPY cmd /app/cmd/
COPY internal /app/internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o kinko ./cmd/main/main.go

RUN apk add --no-cache nodejs npm

COPY ./frontend /app/frontend

RUN cd frontend && mv .env.example .env

RUN cd frontend && npm install
RUN cd frontend && npm run build

RUN cp package.json package-lock.json /app/frontend/build/
RUN cd /app/build && npm ci --omit=dev

RUN apk add --no-cache curl && \
    curl -L -s https://github.com/just-containers/s6-overlay/releases/download/v3.1.6.2/s6-overlay-amd64.tar.gz | tar xvzf - -C /

RUN apk --no-cache add ca-certificates wget \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk \
    && apk add --force-overwrite glibc-2.34-r0.apk

COPY ./docker/services.d /etc/services.d

ENTRYPOINT ["/init"]