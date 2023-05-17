FROM hocusdev/workspace

ARG GO_VERSION="go1.20.4.linux-amd64.tar.gz"
ARG TF_VERSION="1.3.9"

RUN echo $GO_SOURCE_URL
RUN sudo apt-get update -yq

RUN sudo apt-get install -yq \
    unzip \
    gzip

RUN curl https://go.dev/dl/${GO_VERSION} -L -s --output /tmp/${GO_VERSION} \
    && sudo tar -C /usr/local -xzf /tmp/${GO_VERSION} \
    && sudo rm /tmp/${GO_VERSION} \
    && sudo ln -s /usr/local/go/bin/* /usr/local/bin \
    && go install honnef.co/go/tools/cmd/staticcheck@latest 

# install tfenv and terraform version
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv \
    && sudo ln -s ~/.tfenv/bin/* /usr/local/bin \
    && tfenv install ${TF_VERSION} \
    && tfenv use ${TF_VERSION}

# test
RUN go version && terraform version