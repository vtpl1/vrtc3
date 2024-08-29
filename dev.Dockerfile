FROM ubuntu:22.04

RUN apt update \
    && DEBIAN_FRONTEND=noninteractive apt install -y \
    ca-certificates \
    software-properties-common \
    build-essential \
    wget \
    curl \
    git \
    lsb-release && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN sed -i 's/# deb-src/deb-src/' /etc/apt/sources.list \
    && apt update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y build-dep python3 \
    && DEBIAN_FRONTEND=noninteractive apt install -y build-essential \
    gdb \
    lcov \
    pkg-config \
    libbz2-dev \
    libffi-dev \
    libgdbm-dev \
    libgdbm-compat-dev \
    liblzma-dev \
    libncurses5-dev \
    libreadline6-dev \
    libsqlite3-dev \
    libssl-dev \
    lzma \
    ninja-build \
    ccache \    
    zip unzip \
    autoconf autoconf-archive \
    lzma-dev \
    tk-dev \
    uuid-dev \
    zlib1g-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && wget -O Python-3.12.3.tar.xz https://www.python.org/ftp/python/3.12.3/Python-3.12.3.tar.xz \
    && tar xvf Python-3.12.3.tar.xz  \
    && cd Python-3.12.3 \
    && sh configure --enable-optimizations --enable-shared \
    && make install \
    && update-alternatives --install /usr/bin/python python /usr/local/bin/python3.12 0 \
    && cd .. \
    && rm -rf Python-3.12.3 \
    && rm Python-3.12.3.tar.xz

RUN mkdir debs && \
    dpkg --purge --force-remove-reinstreq mongodb-database-tools && \
    wget -q https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2204-x86_64-100.9.4.deb -P ./debs && \
    apt-get install -y -q --no-install-recommends ./debs/*.deb && \
    rm -rf debs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSL "https://repo.mysql.com/apt/ubuntu/pool/mysql-apt-config/m/mysql-apt-config/mysql-apt-config_0.8.29-1_all.deb" -o "mysql-apt-config.deb" && \
    export DEBIAN_FRONTEND=noninteractive && \
    dpkg -i "mysql-apt-config.deb" && \
    apt-get update && \
    apt-get install -y mysql-shell && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN wget -O cmake.sh https://github.com/Kitware/CMake/releases/download/v3.28.1/cmake-3.28.1-linux-x86_64.sh \
    && sh cmake.sh --prefix=/usr/local/ --exclude-subdir && rm -rf cmake.sh

ARG LLVM_VERSION=18
RUN wget https://apt.llvm.org/llvm.sh && chmod +x llvm.sh && ./llvm.sh ${LLVM_VERSION} all

ENV PATH="${PATH}:/usr/lib/llvm-${LLVM_VERSION}/bin"

ARG GO_VERSION=1.23.0

RUN mkdir gotmp && \
    wget -O go.linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go.linux-amd64.tar.gz && \
    rm -rf gotmp

ENV PATH="${PATH}:/usr/local/go/bin"

# RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

ENV SHELL /bin/bash
ENV LANG=en_US.utf-8
ENV LC_ALL=en_US.utf-8

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

RUN groupmod --gid $USER_GID $USERNAME \
    && usermod --uid $USER_UID --gid $USER_GID $USERNAME \
    && chown -R $USER_UID:$USER_GID /home/$USERNAME

USER ${USERNAME}

ARG NODE_VERSION=18

RUN wget -q -O - https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
    && export NVM_DIR="$HOME/.nvm" \
    && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" \
    && nvm install ${NODE_VERSION}

RUN python -m pip install -U pip
RUN python -m pip install poetry
# RUN python -m poetry self add poetry-bumpversion