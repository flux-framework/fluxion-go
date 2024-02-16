HERE ?= $(shell pwd)
LOCALBIN ?= $(shell pwd)/bin
JGF ?= $(HERE)/cmd/test/data/tiny.json
JOBSPECS ?= $(HERE)/cmd/test/data/jobspecs

# This assumes a build in the .devcontainer Dockerfile environment
FLUX_SCHED_ROOT ?= /opt/flux-sched
INSTALL_PREFIX ?= /usr
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
LD_LIBRARY_PATH=/usr/lib:/usr/lib/flux:/usr/local/lib:/usr/local/lib/flux

BUILDENVVAR=CGO_CFLAGS="-I${FLUX_SCHED_ROOT} -I${FLUX_SCHED_ROOT}/resource/reapi/bindings/c" CGO_LDFLAGS="-L${INSTALL_PREFIX}/lib -L${INSTALL_PREFIX}/lib/flux -L${FLUX_SCHED_ROOT}/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp"
# BUILDENVAR=CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS='${CGO_LIBRARY_FLAGS}' go build -ldflags '-w'


.PHONY: all
all: build

.PHONY: test
test: test-binary test-modules

.PHONY: test-modules
test-modules: 
	go test -v ./pkg/types

.PHONY: test-binary
test-binary: 	
	LD_LIBRARY_PATH=$(LD_LIBRARY_PATH) $(LOCALBIN)/test --jgf=$(JGF) --jobspec=$(JOBSPECS)/test001.yaml

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
