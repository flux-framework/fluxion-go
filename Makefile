HERE ?= $(shell pwd)
LOCALBIN ?= $(shell pwd)/bin
JGF ?= $(HERE)/cmd/test/data/tiny.json
JOBSPECS ?= $(HERE)/cmd/test/data/jobspecs
CANCELDATA ?= $(HERE)/cmd/test/data/cancel

# This assumes a build in the .devcontainer Dockerfile environment
FLUX_SCHED_ROOT ?= /opt/flux-sched
INSTALL_PREFIX ?= /usr

# Needed to distinguish /usr/lib and /usr/lib54
LIB_PREFIX ?= /usr/lib
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)

# Note that el8 and derivatives are in /usr/lib64
LD_LIBRARY_PATH=$(LIB_PREFIX):$(LIB_PREFIX)/flux

BUILDENVVAR=CGO_CFLAGS="-I${FLUX_SCHED_ROOT} -I${FLUX_SCHED_ROOT}/resource/reapi/bindings/c" CGO_LDFLAGS="-L${LIB_PREFIX} -L${LIB_PREFIX}/flux -L${FLUX_SCHED_ROOT}/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -ljansson -lhwloc -lflux-hostlist -lboost_graph -lyaml-cpp"

.PHONY: all
all: test

.PHONY: test
test: 	
#	$(COMMONENVVAR) $(BUILDENVVAR) LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) go test -count 1 -run TestCancel -ldflags '-w' ./pkg/fluxcli ./pkg/types
	$(COMMONENVVAR) $(BUILDENVVAR) LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) go test -ldflags '-w' ./pkg/fluxcli ./pkg/types

.PHONY: test-v
test-v: 	
#	$(COMMONENVVAR) $(BUILDENVVAR) LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) go test -count 1 -run TestCancel -v -ldflags '-w' ./pkg/fluxcli ./pkg/types
	$(COMMONENVVAR) $(BUILDENVVAR) LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) go test -v -ldflags '-w' ./pkg/fluxcli ./pkg/types

.PHONY: $(LOCALBIN)
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: clean
clean:
	rm -rf $(LOCALBIN)/test
