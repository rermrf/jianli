FROM jenkins/jenkins:lts-jdk21
USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends bash ca-certificates curl docker.io docker-compose-plugin git rsync \
    && rm -rf /var/lib/apt/lists/*
USER jenkins
