FROM golang1.21-alpine3 as builder
ADD go.mod go.sum /app/
WORKDIR /app
RUN go mod download
ADD . /app
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build  \
    -o requestprinter \
    .

FROM scratch
WORKDIR /requestprinter
ADD --from=builder /app/requestprinter /requestprinter/server
USER nobody
ENTRYPOINT ["/requestprinter/server"]
CMD ["-url", "-headers", "-body", "-method"]
