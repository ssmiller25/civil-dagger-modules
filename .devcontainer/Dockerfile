FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

ARG DAGGER_VERSION=0.13.3

# Install dagger
RUN curl -fsSL https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION=0.13.0 BIN_DIR=/usr/bin sudo -E sh
