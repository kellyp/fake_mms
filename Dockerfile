FROM golang:alpine as golang
WORKDIR /go/src/app
COPY . .
# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o fake_mms
RUN pwd
RUN ls

FROM scratch

COPY --from=golang /go/src/app/fake_mms /fake_mms

EXPOSE 8080/tcp

ENTRYPOINT [ "/fake_mms" ]
