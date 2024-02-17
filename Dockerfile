FROM fluxrm/flux-sched:bookworm-amd64

# Basic container to provide quick developer environment with everything ready to go

LABEL maintainer="Vanessasaurus <@vsoch>"

USER root
RUN apt-get update && apt-get install -y less

# Install Go 19 (we should update this)
RUN wget https://go.dev/dl/go1.19.10.linux-amd64.tar.gz  && tar -xvf go1.19.10.linux-amd64.tar.gz && \
         mv go /usr/local && rm go1.19.10.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/home/vscode/go/bin

RUN git clone https://github.com/flux-framework/flux-sched /opt/flux-sched

# Assuming installing to /usr/local
ENV LD_LIBRARY_PATH=/usr/lib:/usr/lib/flux:/usr/local/lib
WORKDIR /workspace/fluxion-go
COPY . .
