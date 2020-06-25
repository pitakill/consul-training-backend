FROM golang as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -ldflags "-linkmode external -extldflags -static" -a main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
ENV PORT=8080
EXPOSE ${PORT}
COPY --from=builder /build/main /app/main
WORKDIR /app
CMD ["/app/main"]
