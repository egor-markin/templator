###############################################################################
# build stage
###############################################################################

FROM golang:alpine AS builder

WORKDIR /build
COPY . .

# Building the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

###############################################################################
# final stage
###############################################################################

FROM madnight/docker-alpine-wkhtmltopdf

ARG BUILD_DATE

LABEL maintainer="info@compasslabs.org"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.url="https://www.compasslabs.org/"

WORKDIR /app
COPY --from=builder /build/app /app
COPY templates /app/templates
COPY web /app/web

ENTRYPOINT ["/app/app"]