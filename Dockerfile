FROM golang:1.22.1 AS builder
RUN mkdir /tmp/build
COPY . /tmp/build
RUN cd /tmp/build && CGO_ENABLED=0 go build ./cmd/fastroutes

FROM scratch
COPY --from=builder /tmp/build/fastroutes /bin/fastroutes
CMD ["/bin/fastroutes"]
