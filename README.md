# Fluxion Go

> 🚧️ This is a work in progress and should not be used yet! 🚧️

Fluxion Go are the Go bindings for the flux scheduler [flux-framework.org/flux-sched](https://github.com/flux-framework.org/flux-sched) that we call "fluxion." You can read more about the project there. In short, you might want to use these bindings if you're interested in integrating the flux scheduler logic (graph-based, hierarchical scheduling) into your Go applications.

## Usage

Currently, you can build in the [Development Container](.devcontainer) in VSCode and then flux-sched will be installed to `/usr` and header files (.h) available at `/opt/flux-sched`. You can build the test suite as follows:

```bash
make
```
```console
mkdir -p /workspaces/fluxion-go/bin
GOOS=linux CGO_CFLAGS="-I/opt/flux-sched -I/opt/flux-sched/resource/reapi/bindings/c" CGO_LDFLAGS="-L/usr/lib -L/usr/lib/flux -L/opt/flux-sched/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp" go build -ldflags '-w' -o /workspaces/fluxion-go/bin/test cmd/test/test.go
```

If you need to customize the flux install prefix or the location (root) of the flux-sched repository (with header files):

```bash
FLUX_SCHED_ROOT=/home/path/flux-sched INSTALL_PREFIX=/usr/local make
```

Here is how to run tests:

```bash
# run all tests
make test

# run binary test (e.g., build main and run)
make test-binary

# run go native tests
make test-modules
```

More work and updates will be coming soon.

#### License

SPDX-License-Identifier: LGPL-3.0

LLNL-CODE-764420
