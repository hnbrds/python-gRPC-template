FROM golang:latest as go_builder
MAINTAINER jysdev@naver.com

FROM #python image
# multi-stage build for go environment
COPY --from=go_builder /usr/local/go /usr/local/go
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

# env / working directory settings
RUN apt-get update
RUN mkdir -p /workspace/service
COPY . /workspace/service
RUN ln -s /workspace/service $GOROOT/src/
WORKDIR /workspace/service

# setup python evnironment
RUN apt install -y python3.8-venv
RUN python -m venv venv && source ./venv/bin/activate
RUN python -m pip install --upgrade pip setuptools wheel
RUN python -m pip install -r requirements.txt

# clean tmp files
RUN rm -rf /tmp/*

EXPOSE 13270

ENTRYPOINT ["./run_server.sh"]
