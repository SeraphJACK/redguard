FROM git.s8k.top/library/golang-gcc AS build
COPY . /app/
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o app

FROM alpine
COPY --from=build /app/app /usr/local/bin/
VOLUME /log
WORKDIR /log/
ENTRYPOINT [ "/usr/local/bin/app" ]
