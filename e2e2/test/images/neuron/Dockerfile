FROM public.ecr.aws/docker/library/ubuntu:20.04

# Neuron SDK components version numbers
ARG NEURONX_DISTRIBUTED_VERSION=0.8.0
ARG NEURONX_CC_VERSION=2.15.128.0
ARG NEURONX_FRAMEWORK_VERSION=2.1.2.2.3.0
ARG NEURONX_COLLECTIVES_LIB_VERSION=2.22.26.0-17a033bc8
ARG NEURONX_RUNTIME_LIB_VERSION=2.22.14.0-6e27b8d5b
ARG NEURONX_TOOLS_VERSION=2.19.0.0

ARG PYTHON=python3.10
ARG PYTHON_VERSION=3.10.12
ARG PIP=pip3
ARG OMPI_VERSION=4.1.5

# This arg required to stop docker build waiting for region configuration while installing tz data from ubuntu 20
ARG DEBIAN_FRONTEND=noninteractive

# Python won’t try to write .pyc or .pyo files on the import of source modules
# Force stdin, stdout and stderr to be totally unbuffered. Good for logging
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PYTHONIOENCODING=UTF-8
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/opt/aws/neuron/lib"
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/opt/amazon/efa/lib"
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/opt/amazon/efa/lib64"
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/opt/amazon/openmpi/lib64"
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/usr/local/lib"
ENV PATH /opt/aws/neuron/bin/:$PATH
# ENV SAGEMAKER_TRAINING_MODULE=sagemaker_pytorch_container.training:main
ENV DGLBACKEND=pytorch

RUN apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    cmake \
    curl \
    emacs \
    git \
    jq \
    libopencv-dev \
    software-properties-common \
    wget \
    unzip \
    vim \
    zlib1g-dev \
    openssl \
    libssl-dev \
    libsqlite3-dev \
    libgdbm-dev \
    libc6-dev \
    libbz2-dev \
    libncurses-dev \
    tk-dev \
    libffi-dev \
    libcap-dev \
    gnupg2 \
    gpg-agent \
 && rm -rf /var/lib/apt/lists/* \
 && apt-get clean

RUN echo "deb https://apt.repos.neuron.amazonaws.com focal main" > /etc/apt/sources.list.d/neuron.list
RUN wget -qO - https://apt.repos.neuron.amazonaws.com/GPG-PUB-KEY-AMAZON-AWS-NEURON.PUB | apt-key add -

RUN apt-get update \
 && apt-get install -y \
    aws-neuronx-tools=$NEURONX_TOOLS_VERSION \
    aws-neuronx-collectives=$NEURONX_COLLECTIVES_LIB_VERSION \
    aws-neuronx-runtime-lib=$NEURONX_RUNTIME_LIB_VERSION \
 && rm -rf /var/lib/apt/lists/* \
 && rm -rf /tmp/tmp* \
 && apt-get clean

# Install Open MPI
RUN mkdir -p /tmp/openmpi \
 && cd /tmp/openmpi \
 && wget --quiet https://download.open-mpi.org/release/open-mpi/v4.1/openmpi-${OMPI_VERSION}.tar.gz \
 && tar zxf openmpi-${OMPI_VERSION}.tar.gz \
 && cd openmpi-${OMPI_VERSION} \
 && ./configure --enable-orterun-prefix-by-default \
 && make -j $(nproc) all \
 && make install \
 && ldconfig \
 && rm -rf /tmp/openmpi

# install Python
RUN wget -q https://www.python.org/ftp/python/$PYTHON_VERSION/Python-$PYTHON_VERSION.tgz \
 && tar -xzf Python-$PYTHON_VERSION.tgz \
 && cd Python-$PYTHON_VERSION \
 && ./configure --enable-shared --prefix=/usr/local \
 && make -j $(nproc) && make install \
 && cd .. && rm -rf ../Python-$PYTHON_VERSION* \
 && ln -s /usr/local/bin/pip3 /usr/bin/pip \
 && ln -s /usr/local/bin/$PYTHON /usr/local/bin/python \
 && ${PIP} --no-cache-dir install --upgrade \
    pip \
    setuptools

WORKDIR /

# The ENV variables declared below are changed in the previous section
# Grouping these ENV variables in the first section causes
# ompi_info to fail. This is only observed in CPU containers
ENV PATH="$PATH:/home/.openmpi/bin"
ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/home/.openmpi/lib/"
RUN ompi_info --parsable --all | grep mpi_built_with_cuda_support:value

RUN ${PIP} install --no-cache-dir -U \
    "bokeh>=2.3,<3" \
    "awscli<2" \
    scipy \
    click \
    "cryptography" \
    psutil==5.6.7 \
    dataset \
    transformers==4.36.2 \
    Pillow

RUN mkdir -p /etc/pki/tls/certs && cp /etc/ssl/certs/ca-certificates.crt /etc/pki/tls/certs/ca-bundle.crt
RUN ${PIP} config set global.extra-index-url https://pip.repos.neuron.amazonaws.com \
 && ${PIP} install --force-reinstall torch-neuronx==$NEURONX_FRAMEWORK_VERSION --extra-index-url https://pip.repos.neuron.amazonaws.com \
 && ${PIP} install --force-reinstall neuronx-cc==$NEURONX_CC_VERSION --extra-index-url https://pip.repos.neuron.amazonaws.com \
 && ${PIP} install --force-reinstall --no-deps neuronx_distributed==$NEURONX_DISTRIBUTED_VERSION --extra-index-url https://pip.repos.neuron.amazonaws.com

# attrs, neuronx-cc required: >=19.2.0, sagemaker <24,>=23.1.0
# protobuf neuronx-cc<4, sagemaker-training >=3.9.2,<=3.20.3
# awscli 1.25.47 has requirement docutils<0.17,>=0.10
# etcd for kubernetes installation
# awscli 1.27.127 has requirement rsa<4.8,>=3.1.2, but you have rsa 4.9.
# awscli 1.27.127 requires urllib3 < 1.27, python-etcd requires urllib3 >= 1.7, latest urllib3 release is 2.0.2
RUN ${PIP} install --no-cache-dir -U \
    "attrs<24,>=23.1.0" \
    "protobuf>=3.18.3,<=3.20.3" \
    "docutils>=0.10,<0.17" \
    "rsa<4.8,>=3.1.2" \
    "urllib3>=1.26.0,<1.27"

# EFA Installer does apt get. Make sure to run apt update before that
RUN apt-get update
RUN cd $HOME \
 && curl -O https://efa-installer.amazonaws.com/aws-efa-installer-latest.tar.gz \
 && wget https://efa-installer.amazonaws.com/aws-efa-installer.key && gpg --import aws-efa-installer.key \
 && cat aws-efa-installer.key | gpg --fingerprint \
 && wget https://efa-installer.amazonaws.com/aws-efa-installer-latest.tar.gz.sig && gpg --verify ./aws-efa-installer-latest.tar.gz.sig \
 && tar -xf aws-efa-installer-latest.tar.gz \
 && cd aws-efa-installer \
 && ./efa_installer.sh -y -g -d --skip-kmod --skip-limit-conf --no-verify \
 && cd $HOME


# Clean up after apt update
RUN rm -rf /var/lib/apt/lists/* \
 && rm -rf /tmp/tmp* \
 && apt-get clean

# Install some common packages used by training scripts
# torchvision needed for MLP. since it depends on torch and torch neuron/torch
# is already installed install it with nodeps
RUN pip3 install --no-cache-dir --no-deps -U \
    torchvision==0.16.*


RUN HOME_DIR=/root \
 && curl -o ${HOME_DIR}/oss_compliance.zip https://aws-dlinfra-utilities.s3.amazonaws.com/oss_compliance.zip \
 && unzip ${HOME_DIR}/oss_compliance.zip -d ${HOME_DIR}/ \
 && cp ${HOME_DIR}/oss_compliance/test/testOSSCompliance /usr/local/bin/testOSSCompliance \
 && chmod +x /usr/local/bin/testOSSCompliance \
 && chmod +x ${HOME_DIR}/oss_compliance/generate_oss_compliance.sh \
 && ${HOME_DIR}/oss_compliance/generate_oss_compliance.sh ${HOME_DIR} ${PYTHON} \
 && rm -rf ${HOME_DIR}/oss_compliance* \
 && rm -rf /tmp/tmp*

RUN curl -o /license.txt  https://aws-dlc-licenses.s3.amazonaws.com/pytorch-2.1/license.txt

COPY e2e2/test/images/neuron/tests ./tests