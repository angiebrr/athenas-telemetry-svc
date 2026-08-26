# =================================================================================================
# BUILD STAGE
# =================================================================================================
FROM golang:1.26.7-alpine AS builder

WORKDIR /app

# Copy the manifests on their own first so `go mod download` stays cached until dependencies
# actually change, rather than on every source edit.
#
# NOTE: go.sum comes along with go.mod so the download is verified against the lock file. Copying
# go.mod alone silently skips that check.
COPY ./athenas-telemetry-svc/go.mod ./athenas-telemetry-svc/go.sum ./
RUN go mod download

COPY ./athenas-telemetry-svc .

# CGO_ENABLED=0 produces a statically linked binary, which is what lets the runtime stage below
# drop the go toolchain and libc entirely.
RUN CGO_ENABLED=0 go build -o svr ./cmd/athenas

# =================================================================================================
# RUNTIME STAGE
# =================================================================================================
FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/svr ./svr

# Run as a non-root user -- nothing the service does needs root, and the acceptance tests bind a
# port well above 1024.
RUN adduser --disabled-password --no-create-home --uid 10001 athenas
USER athenas

EXPOSE 8080

CMD [ "./svr" ]
