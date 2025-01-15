FROM fluxrm/flux-sched:bookworm-amd64

LABEL maintainer="Vanessasaurus <@vsoch>"

# Match the default user id for a single system so we aren't root
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000
ENV USERNAME=${USERNAME}
ENV USER_UID=${USER_UID}
ENV USER_GID=${USER_GID}
ENV GO_VERSION=1.21.9

USER root
RUN apt-get update && apt-get install -y less

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz  && tar -xvf go${GO_VERSION}.linux-amd64.tar.gz && \
         mv go /usr/local && rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/home/vscode/go/bin

# Testing grow/shrink from custom branch
# RUN git clone -b debug-resource-error-messages https://github.com/researchapps/flux-sched /opt/flux-sched
RUN git clone -b grow-api https://github.com/milroy/flux-sched /opt/flux-sched
# RUN git clone https://github.com/flux-framework/flux-sched /opt/flux-sched

# We also need to rebuild into the system install
# TODO remove this when it is added to production container
ENV FLUX_SCHED_VERSION=0.39.0
RUN cd /opt/flux-sched && \
    cmake -B build && \
    make -C build && \
    make -C build install

# Add the group and user that match our ids
RUN groupadd -g ${USER_GID} ${USERNAME} && \
    adduser --disabled-password --uid ${USER_UID} --gid ${USER_GID} --gecos "" ${USERNAME} && \
    echo "${USERNAME} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers

USER $USERNAME
# Assuming installing to /usr/local
ENV LD_LIBRARY_PATH=/usr/lib:/usr/lib/flux:/usr/local/lib
