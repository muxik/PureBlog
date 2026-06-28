# ---- build ----
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pureblog ./cmd/pureblog
# Skeleton upload dir owned by the nonroot uid. When a fresh named volume is
# mounted at /data, Docker seeds it from the image — inheriting this ownership —
# so the nonroot process can write uploads without a chown at runtime (the
# distroless image has no shell to chown in).
RUN mkdir -p /data/uploads && chown -R 65532:65532 /data

# ---- run ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/pureblog /pureblog
COPY --from=build --chown=65532:65532 /data /data
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/pureblog"]
