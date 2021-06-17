DIR = $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

PACKAGE = fibonacci

COMPANY = reserve_trust

define SET_LDFLAGS
	-X fibonacci/models.DB_URL=$(1)  \
	-X fibonacci/models.DB_PORT=$(2) \
	-X fibonacci/models.DB_USER=$(3) \
	-X fibonacci/models.DB_PW=$(4)   \
	-X fibonacci/models.DB_NAME=$(5) \
	-X main.PORT=$(6)
endef

BUILD = ${DIR}/build
BUILD_DIRS = ${BUILD}

APP_BUILD = ${BUILD}/${PACKAGE}_build

MAKEFILE = $(abspath $(lastword $(MAKEFILE_LIST)))
DOCKERFILE = ${DIR}/Dockerfile
SRC_DIR = ${DIR}/src
SRC = $(shell find ${SRC_DIR} -type f -name *.go) ${DOCKERFILE} ${MAKEFILE}

SCRIPTS = ${DIR}/scripts

.PHONY: all
all: ${BUILD_DIRS} dev test app compose

.PHONY: dev
dev:
	$(eval LD_FLAGS = $(call SET_LDFLAGS,postgres,5432,reserve_trust,foo,fibonacci,8080))

${BUILD_DIRS}:
	@mkdir -p $@

${APP_BUILD}: ${BUILD_DIRS} ${SRC}
	@docker build                           \
		--build-arg PACKAGE=${PACKAGE}      \
		--build-arg LD_FLAGS='${LD_FLAGS}'  \
		-t ${COMPANY}                       \
		.
	@touch $@

.PHONY: test
test:
	@docker build                           \
		--build-arg PACKAGE=${PACKAGE}      \
		--build-arg LD_FLAGS='${LD_FLAGS}'  \
		-t ${COMPANY}_test                  \
		--target ${COMPANY}_test            \
		.
	
.PHONY: app
app: test ${APP_BUILD}
	@docker build                           \
		--build-arg PACKAGE=${PACKAGE}      \
		--build-arg LD_FLAGS='${LD_FLAGS}'  \
		-t ${COMPANY}_application           \
		--target ${COMPANY}_application     \
		.

.PHONY: compose
compose: app
	@docker-compose up -d

.PHONY: clean
clean:
	@echo cleaning ${BUILD_DIRS} and removing docker images
	@rm -rf ${BUILD_DIRS}
	@${SCRIPTS}/kill_docker

