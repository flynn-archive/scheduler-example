FROM flynn/busybox
MAINTAINER Jeff Lindsay <progrium@gmail.com>

ADD ./build/primesum /bin/primesum

ENTRYPOINT ["time", "/bin/primesum"]
