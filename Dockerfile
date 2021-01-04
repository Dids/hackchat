## FIXME: The Go builder should also be able to support aarch64/Apple Silicon, so Go 1.1.6?

## FIXME: This breaks when running on ARM (git rev-list specifically fails to run correctly)
# Compiler image
FROM didstopia/base:go-alpine-3.5 AS go-builder

# Copy the project 
COPY . /tmp/hackchat/
WORKDIR /tmp/hackchat/

# Install dependencies
RUN make deps

# Build and test the binary
RUN make build && \
    cp hackchat /go/bin/hackchat && \
    /go/bin/hackchat version
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/hackchat



# Runtime image
FROM scratch

# Copy certificates
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the built binary
COPY --from=go-builder /go/bin/hackchat /go/bin/hackchat

# Expose environment variables
ENV HACKMUD_CHAT_API_TOKEN  ""
ENV HACKMUD_CHANNEL_ID      ""
ENV HACKMUD_OWNER_ID        ""
ENV DISCORD_API_TOKEN       ""

## FIXME: Can we not do this on the scratch image?
## OCI runtime create failed: container_linux.go:370: starting container process caused: exec: "/bin/sh": stat /bin/sh: no such file or directory: unknown
# Test binary
# RUN /go/bin/hackchat version

# Expose volumes
VOLUME [ "/.db" ]

# Run the binary
ENTRYPOINT ["/go/bin/hackchat"]
