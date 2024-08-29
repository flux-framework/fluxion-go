# Fluxion Go

Fluxion Go provides the Go bindings for the flux scheduler [flux-framework.org/flux-sched](https://github.com/flux-framework.org/flux-sched) that we call "fluxion." You can read more about the project there. In short, you might want to use these bindings if you're interested in integrating the flux scheduler logic (graph-based, hierarchical scheduling) into your Go applications.

## Usage

Currently, you can build in the [Development Container](.devcontainer) in VSCode and then flux-sched will be installed to `/usr` and header files (.h) available at `/opt/flux-sched`. You can build the test suite as follows:

```bash
make
```
```console
mkdir -p /workspaces/fluxion-go/bin
GOOS=linux CGO_CFLAGS="-I/opt/flux-sched -I/opt/flux-sched/resource/reapi/bindings/c" CGO_LDFLAGS="-L/usr/lib -L/usr/lib/flux -L/opt/flux-sched/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lflux-hostlist -lboost_graph -lyaml-cpp" go build -ldflags '-w' -o /workspaces/fluxion-go/bin/test cmd/test/test.go
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

## Docker 

In addition to the developer environment, we provide an automated build with the [Dockerfile](Dockerfile) here
that will give you a containerized environment with Go, the bindings, and flux-sched. You can pull the repository
from our package registry, or build on your own:

```bash
docker build -t ghcr.io/flux-framework/fluxion-go .
docker run -it ghcr.io/flux-framework/fluxion-go
```

Then you can build, and test.

```bash
make
make test
```

Have fun! 🧞‍♀️️

#### License

SPDX-License-Identifier: LGPL-3.0

LLNL-CODE-764420
