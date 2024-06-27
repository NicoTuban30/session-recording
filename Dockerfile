# syntax=docker/dockerfile:1

## UI
FROM node:19-alpine AS ui

RUN npm i -g pnpm

WORKDIR /src/ui
COPY ui/package.json ui/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

# Copy the rest of the UI source code
COPY ui/ .

RUN pnpm run build


## Build
FROM golang:1-alpine AS build

WORKDIR /src

RUN apk add --no-cache build-base

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o server cmd/server/main.go


## Deploy
FROM alpine

RUN apk add --no-cache tini ca-certificates

WORKDIR /
COPY --from=build /src/server server
COPY --from=ui /src/ui/dist ./public

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/server"]
