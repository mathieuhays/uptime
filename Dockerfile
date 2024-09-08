FROM golang:1.23 AS build

WORKDIR /src

# cache dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy source code
COPY ./ ./

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

RUN go build -v -o /app ./cmd/server

FROM ubuntu

COPY --from=build /app /app

ENV PORT=80
ENV DATABASE_PATH=/data/uptime.db

EXPOSE 80

ENTRYPOINT ["/app"]