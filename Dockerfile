FROM flynn/busybox
MAINTAINER Jeff Lindsay <progrium@gmail.com>

ADD ./build/scheduler-example /bin/scheduler-example

ENTRYPOINT ["/bin/scheduler-example"]
