FROM golang:alpine3.18 AS builder

COPY ${PWD} /app
WORKDIR /app

# Toggle CGO on your app requirement
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/app *.go
# Use below if using vendor
# RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-extldflags "-static"' -o /app/appbin *.go

FROM golang:alpine3.18
LABEL MAINTAINER=LowK

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# Add new user 'staff'. App should be run without root privileges as a security measure
RUN adduser --home "/staff" --disabled-password staff --gecos "staff,-,-,-"
USER staff

COPY --from=builder /app/app /home/staff/app

WORKDIR /home/staff/app

CMD ["./app"]
