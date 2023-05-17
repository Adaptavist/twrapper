FROM gitpod/workspace-go

ARG AWS_DEST="awscliv2.zip"

USER root

RUN apt-get install -yq \
    curl \
    unzip \
    less

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -L -o "${AWS_DEST}" \
    && unzip ${AWS_DEST} \
    && ./aws/install \
    && rm -rf ${AWS_DEST} ./amazonaws

USER gitpod

# We want a newer version of go that in the image, so we're going to overwrite it
RUN curl "https://go.dev/dl/go1.20.4.linux-amd64.tar.gz" -L -o "go1.20.4.linux-amd64.tar.gz" \
    && rm -rf ~/go \
    && tar -C ~/ -xzf go1.20.4.linux-amd64.tar.gz

RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv \
    && echo 'export PATH="$HOME/.tfenv/bin:$PATH"' >> ~/.bashrc.d/300-tfenv \
    && echo 'export AWS_PROFILE=main' >> ~/.bashrc.d/300-aws
