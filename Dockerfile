FROM golang:1.19-alpine as build

RUN apk update && \
    apk upgrade && \
    apk add --update libc-dev gcc make rsync grep ca-certificates openssl wget

WORKDIR /go_wkspc/src/github.com/FuTuL

COPY . /go_wkspc/src/github.com/FuTuL/jobber

RUN cd jobber && \
    make check && \
    make install DESTDIR=/jobber-dist/

FROM alpine:3.16

RUN mkdir /jobber
COPY --from=build /jobber-dist/usr/local/libexec/jobberrunner /jobber/jobberrunner
COPY --from=build /jobber-dist/usr/local/bin/jobber /jobber/jobber
COPY --from=build /jobber-dist/etc/jobber.conf /etc/jobber.conf
ENV PATH /jobber:${PATH}

ENV USERID 1000
RUN addgroup jobberuser && \
    adduser -S -u "${USERID}" -G jobberuser jobberuser && \
    mkdir -p "/var/jobber/${USERID}" && \
    chown -R jobberuser:jobberuser "/var/jobber/${USERID}"

COPY --chown=jobberuser:jobberuser .jobber /home/jobberuser/.jobber

RUN chmod 0600 /home/jobberuser/.jobber
RUN mkdir -p /usr/local/var/jobber/1000 && chmod 777 -R /usr/local/var/jobber

USER jobberuser

CMD ["jobberrunner", "-u", "/usr/local/var/jobber/1000/cmd.sock", "/home/jobberuser/.jobber"]
