FROM flynn/busybox
MAINTAINER Jonathan Rudenberg <jonathan@titanous.com>

ADD ./build/primesum /bin/primesum

ENTRYPOINT ["time", "/bin/primesum"]
