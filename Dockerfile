FROM golang:1.12 as builder

ARG GOPROXY
ENV GORPOXY ${GOPROXY}

ADD . /builder

WORKDIR /builder

RUN go build main.go && go build api.go

FROM golang:1.12

COPY --from=builder /builder/main /app/site-monitor

COPY --from=builder /builder/api /app/site-monitor-api

WORKDIR /app
