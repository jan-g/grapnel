FROM fnproject/go:dev as build-stage
WORKDIR /function
ADD . /go/src/func/
ENV GOPATH=
ENV GOFLAGS=-mod=vendor
RUN cd /go/src/func/ && go mod vendor -v
RUN cd /go/src/func/ && go build -o func
FROM fnproject/go
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/

# Make the image a bit more useful
RUN sed -i -e '/^root:/s|/root|/tmp|' /etc/passwd && \
    rmdir /run && ln -s /run /tmp/run && \
    ln -s /tmp/etc/sshd /etc/sshd && \
    ln -s /tmp/etc/dropbear /etc/dropbear && \
    apk --no-cache add dropbear openssh-client && \
    apk --no-cache add bind-tools curl iproute2
COPY start /start
ENTRYPOINT ["/start"]
