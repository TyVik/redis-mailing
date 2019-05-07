# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12.
ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache git ca-certificates

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user \
    && echo '1000:x:65534:65534:1000:/:' > /user/passwd \
    && echo '1000:x:65534:' > /user/group

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build  .

FROM scratch

# Import the user and group files from the build stage.
COPY --from=build /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /src/redis-mailing .

USER 1000:1000

ENTRYPOINT ["/redis-mailing"]
