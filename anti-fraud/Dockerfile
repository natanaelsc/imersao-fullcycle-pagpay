FROM node:23.11.0-slim

RUN apt update && \
    apt install openssl procps -y && \
    apt clean && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/* && \
    npm install -g @nestjs/cli@11.0.6

USER node

WORKDIR /home/node/app

CMD [ "tail", "-f", "/dev/null" ]
