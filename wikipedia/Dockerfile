
FROM golang as step1

RUN mkdir /app
RUN mkdir /app/wikipedia
COPY . /app

ENV CGO_ENABLED=0
WORKDIR /app/cmd/web
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR ./
COPY --from=step1 /app /app/wikipedia

CMD ["/app/wikipedia/cmd/web/main"]
