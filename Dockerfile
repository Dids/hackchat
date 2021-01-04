# Compiler image
FROM didstopia/base:go-alpine-3.5 AS go-builder

# Copy the project 
COPY . /tmp/hackchat/
WORKDIR /tmp/hackchat/

# Install dependencies
RUN make deps

# Build the binary
#RUN make build && ls /tmp/hackchat
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/hackchat



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

# Expose volumes
VOLUME [ "/.db" ]

# Run the binary
ENTRYPOINT ["/go/bin/hackchat"]
