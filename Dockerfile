FROM golang as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main .
ENV CONSUL_TEMPLATE_VERSION 0.25.0
RUN curl -sL https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz -o consul-template.tgz && tar xfvz consul-template.tgz

FROM alpine
RUN adduser -S -D -H -h /app appuser
RUN apk add --no-cache curl
USER appuser
ENV PORT=8080
EXPOSE ${PORT}
COPY --from=builder /build/consul-template /bin/consul-template
COPY --from=builder /build/main /app/main
WORKDIR /app
CMD ["/app/main"]
