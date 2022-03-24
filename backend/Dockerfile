FROM golang:1.18-bullseye

EXPOSE 9080
ARG USER="server"
RUN useradd ${USER}  && \
    mkdir -p /home/${USER} && \ 
    chown ${USER} /home/${USER}
COPY --chown=${USER} . /home/${USER}/

USER ${USER}
WORKDIR /home/${USER}