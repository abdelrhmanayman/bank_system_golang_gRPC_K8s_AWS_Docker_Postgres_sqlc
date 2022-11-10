FROM golang:1.19-alpine AS buildStage
WORKDIR /app
COPY . .
RUN go build -o main index.go

FROM alpine
WORKDIR /app
COPY --from=buildStage /app/main .
COPY .env .
ENV GIN_MODE=release

EXPOSE 8080
CMD ["/app/main"]