FROM czarsimon/godep:1.11.2-alpine3.8 as build

# Copy source
WORKDIR /go/src/httptrigger
COPY . .

# Install dependencies
RUN dep ensure

# Build application
WORKDIR /go/src/httptrigger
RUN go build

FROM alpine:3.8 as run
RUN mkdir /etc/mimir /etc/mimir/httptrigger

WORKDIR /opt/app
COPY --from=build /go/src/httptrigger/httptrigger httptrigger
CMD ["./httptrigger"]