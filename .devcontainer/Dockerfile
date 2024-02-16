FROM fluxrm/flux-sched:bookworm

LABEL maintainer="Vanessasaurus <@vsoch>"

# Match the default user id for a single system so we aren't root
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000
ENV USERNAME=${USERNAME}
ENV USER_UID=${USER_UID}
ENV USER_GID=${USER_GID}
USER root
RUN apt-get update && apt-get install -y less

# Install Go 19 (we should update this)
RUN wget https://go.dev/dl/go1.19.10.linux-amd64.tar.gz  && tar -xvf go1.19.10.linux-amd64.tar.gz && \
         mv go /usr/local && rm go1.19.10.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/home/vscode/go/bin

RUN git clone https://github.com/flux-framework/flux-sched /opt/flux-sched && \
    cd /opt/flux-sched && \
    git fetch && \
    export FLUX_SCHED_VERSION=0.58.0 && \
    ./autogen.sh && \
    ./configure --prefix=/usr && \
    make && make install

# Add the group and user that match our ids
RUN groupadd -g ${USER_GID} ${USERNAME} && \
    adduser --disabled-password --uid ${USER_UID} --gid ${USER_GID} --gecos "" ${USERNAME} && \
    echo "${USERNAME} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers

USER $USERNAME
# Assuming installing to /usr/local
ENV LD_LIBRARY_PATH=/usr/lib:/usr/lib/flux:/usr/local/lib