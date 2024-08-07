# Start with the NVIDIA CUDA base image
FROM nvidia/cuda:12.5.0-devel-ubuntu22.04

ARG EFA_INSTALLER_VERSION=latest
# 1.7.4+ is required, to enforce proper EFA function with OFI_NCCL_DISABLE_GDR_REQUIRED_CHECK=0
ARG AWS_OFI_NCCL_VERSION=1.9.1
ARG NCCL_TESTS_VERSION=master

# Install necessary dependencies
RUN apt-get update -y
RUN apt-get remove -y --allow-change-held-packages \
    libmlx5-1 \
    ibverbs-utils \
    libibverbs-dev \
    libibverbs1 \
    libnccl2 \
    libnccl-dev

RUN rm -rf /opt/hpcx \
    && rm -rf /usr/local/mpi \
    && rm -rf /usr/local/ucx \
    && rm -f /etc/ld.so.conf.d/hpcx.conf \
    && ldconfig

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-unauthenticated \
    sudo \
    git \
    gcc \
    vim \
    kmod \
    openssh-client \
    openssh-server \
    build-essential \
    wget curl \
    autoconf \
    libtool \
    gdb \
    automake \
    python3-distutils \
    cmake \
    apt-utils \
    devscripts \
    debhelper \
    libsubunit-dev \
    check \
    pkg-config \
    libhwloc-dev \
    datacenter-gpu-manager \
    cloud-utils \ 
    cuda-demo-suite-12-5

RUN mkdir -p /var/run/sshd
RUN sed -i 's/[ #]\(.*StrictHostKeyChecking \).*/ \1no/g' /etc/ssh/ssh_config && \
    echo "    UserKnownHostsFile /dev/null" >> /etc/ssh/ssh_config && \
    sed -i 's/#\(StrictModes \).*/\1no/g' /etc/ssh/sshd_config
ENV LD_LIBRARY_PATH /opt/amazon/openmpi/lib64:/opt/amazon/openmpi/lib:/opt/amazon/efa/lib64:/opt/aws-ofi-nccl/install/lib:/usr/local/cuda/lib64:/usr/local/cuda/extras/CUPTI/lib64:/usr/local/lib/:/usr/lib64:/usr/lib/x86_64-linux-gnu/:$LD_LIBRARY_PATH
ENV PATH /usr/local/cuda/bin:/opt/amazon/openmpi/bin:/opt/amazon/efa/bin:/usr/sbin:/usr/bin:/usr/local/bin:$PATH

# Install EFA
RUN cd $HOME \
    && curl -O https://efa-installer.amazonaws.com/aws-efa-installer-${EFA_INSTALLER_VERSION}.tar.gz \
    && tar -xf $HOME/aws-efa-installer-${EFA_INSTALLER_VERSION}.tar.gz \
    && cd aws-efa-installer \
    && ./efa_installer.sh -y -g -d --skip-kmod --skip-limit-conf --no-verify \
    && rm -rf $HOME/aws-efa-installer

# Install NCCL
RUN apt-key del 7fa2af80 \
    && curl -L -O https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/cuda-keyring_1.0-1_all.deb \
    && dpkg -i cuda-keyring_1.0-1_all.deb \
    && sudo apt install libnccl2=2.18.5-1+cuda12.2 libnccl-dev=2.18.5-1+cuda12.2

## Install AWS-OFI-NCCL plugin
RUN export OPAL_PREFIX="" \
    && git clone https://github.com/aws/aws-ofi-nccl.git /opt/aws-ofi-nccl \
    && cd /opt/aws-ofi-nccl \
    && git checkout v${AWS_OFI_NCCL_VERSION}-aws \
    && ./autogen.sh \
    && ./configure --prefix=/opt/aws-ofi-nccl/install \
       --with-libfabric=/opt/amazon/efa/ \
       --with-cuda=/usr/local/cuda \
       --with-mpi=/opt/amazon/openmpi/ \
    && make && make install
    
# Install NCCL Tests
RUN git clone https://github.com/NVIDIA/nccl-tests.git /opt/nccl-tests \
    && cd /opt/nccl-tests \
    && git checkout ${NCCL_TESTS_VERSION} \
    && make MPI=1 \
       MPI_HOME=/opt/amazon/openmpi/ \
       CUDA_HOME=/usr/local/cuda 


# Set a default command for debugging or modify as per requirements
ENV NCCL_PROTO simple
RUN rm -rf /var/lib/apt/lists/*
ENV LD_PRELOAD /usr/lib/x86_64-linux-gnu/libnccl.so:$LD_PRELOAD

COPY e2e2/test/images/nvidia/gpu_unit_tests ./gpu_unit_tests
RUN chmod +x ./gpu_unit_tests/unit_test