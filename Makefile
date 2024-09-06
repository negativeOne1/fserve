###############################################################################
#
#               ❇❇❇ Welcome to the project Makefile ❇❇❇
#
# This Makefiles holds all related adaptations like installation of external
# libraries.
#
###############################################################################
IMG=harbor3.piplanning.io/stable/fserve:latest

include ./deployment/Makefile

run: $(BIN)
	$(Q)./bin/app run


.PHONY: build-docker-tilt
build-docker-tilt: $(BINS)## Build the docker image
	$(Q)$(ECHO) $(call UPPER, $@)
ifndef IMG
	$(error IMG variable is not set)
endif
	$(Q) DOCKER_BUILDKIT=1 docker build --ssh default . -t ${IMG} --build-arg="BASE_IMAGE=${BASE_IMAGE}"| $(FORMAT)
