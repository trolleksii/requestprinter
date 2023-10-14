FROM golang:1.21-alpine as builder
ADD go.mod /app/
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
COPY --from=builder /app/requestprinter /requestprinter/server
USER nobody
ENTRYPOINT ["/requestprinter/server"]
CMD ["-url", "-headers", "-body", "-method"]
