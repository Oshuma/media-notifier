BIN_DIR=build
BIN=$(BIN_DIR)/media-notifier
BUILD_SHA=$(shell git describe --always --long --dirty)

BUILD_OPTS=-v -buildmode exe -trimpath -o ${BIN}
LDFLAGS=-ldflags "-X main.BuildTag=${BUILD_SHA}"

.PHONY: all
all:
	[ -d ${BIN_DIR} ] || mkdir -p ${BIN_DIR}
	go build ${BUILD_OPTS} ${LDFLAGS}

.PHONY: clean
clean:
	if [ -f ${BIN} ] ; then rm -v ${BIN} ; fi
