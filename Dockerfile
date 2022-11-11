FROM golang:1.19-alpine AS buildStage
WORKDIR /app
COPY . .
RUN apk add curl
RUN go build -o main index.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine
WORKDIR /app
COPY --from=buildStage /app/main .
COPY --from=buildStage /app/migrate ./migrate
COPY db/migration ./migration
COPY start.sh .
COPY wait-for.sh .


COPY .env .
ENV GIN_MODE=release

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]