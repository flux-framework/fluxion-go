FROM fluxrm/flux-sched:bookworm-amd64

# Basic container to provide quick developer environment with everything ready to go

LABEL maintainer="Vanessasaurus <@vsoch>"
ENV GO_VERSION=1.21.9

USER root
RUN apt-get update && apt-get install -y less

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz  && tar -xvf go${GO_VERSION}.linux-amd64.tar.gz && \
         mv go /usr/local && rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/home/vscode/go/bin

RUN git clone -b grow-api https://github.com/milroy/flux-sched /opt/flux-sched

# Assuming installing to /usr/local
ENV LD_LIBRARY_PATH=/usr/lib:/usr/lib/flux:/usr/local/lib
WORKDIR /workspace/fluxion-go
COPY . .
