FROM flynn/busybox
MAINTAINER Jonathan Rudenberg <jonathan@titanous.com>

ADD ./build/scheduler-example /bin/scheduler-example

ENTRYPOINT ["/bin/scheduler-example"]
