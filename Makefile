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
all: build

.PHONY: test
test: test-binary test-modules

.PHONY: test-modules
test-modules: 
	go test -v ./pkg/types

.PHONY: test-binary
test-binary: 	
	LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) $(LOCALBIN)/test --jgf=$(JGF) --jobspecs=$(JOBSPECS) --cancel=$(CANCELDATA)

# test001_desc="match allocate 1 slot: 1 socket: 1 core (pol=default)"
# test_expect_success "${test001_desc}" '
#    ${main} --jgf=${jgf} --jobspec=${jobspec1} > 001.R.out &&
#    sed -i -E "s/, 0\.[0-9]+//g" 001.R.out &&
#    test_cmp 001.R.out ${exp_dir}/001.R.out
#'

#test002_desc="match allocate 2 slots: 2 sockets: 5 cores 1 gpu 6 memory"
#test_expect_success "${test002_desc}" '
#    ${main} --jgf=${jgf} --jobspec=${jobspec2} > 002.R.out &&
#    sed -i -E "s/, 0\.[0-9]+//g" 002.R.out &&
#    test_cmp 002.R.out ${exp_dir}/002.R.out
#'

.PHONY: $(LOCALBIN)
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# This serves as a single test file to build a dummy main to test
.PHONY: build $(LOCALBIN)
build: 
	mkdir -p $(LOCALBIN)
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-w' -o $(LOCALBIN)/test cmd/test/test.go

.PHONY: clean
clean:
	rm -rf $(LOCALBIN)/test
