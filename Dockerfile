FROM golang:latest as builder

WORKDIR /boa
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./boa ./cmd/boa

FROM golang:latest
LABEL maintainer="Dominik Roos <domi.roos@gmail.com>"

ARG WORKDIR=/workdir
RUN mkdir -p ${WORKDIR}
WORKDIR /workdir
COPY --from=builder /boa/boa /usr/local/bin

ENTRYPOINT ["boa"]
