FROM golang AS build1
WORKDIR /build/app
COPY ./ .
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLE=0
RUN go build -ldflags "-w -s"

FROM debian AS build2
WORKDIR /app
COPY --from=build1 /build/app/imageCreator /build/app/config.yaml /app/
RUN apt update
RUN apt install -y upx-ucl
RUN upx imageCreator

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=build2 /app/imageCreator /app/config.yaml /app/
EXPOSE 99
ENTRYPOINT ["./imageCreator"]