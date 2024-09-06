ARG BASE_IMAGE=gcr.io/distroless/static-debian12:nonroot

FROM golang:1.23-bullseye AS builder

WORKDIR /app

RUN --mount=type=ssh git config --global url.git@github.com:.insteadOf https://github.com/ \
    && mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .

RUN --mount=type=ssh make build

FROM $BASE_IMAGE

ENV USER=10001
ENV HOME=/app
ENV BIN=$HOME/bin
ENV PATH=$PATH:$BIN

USER $USER

COPY --from=builder --chown=$USER:$USER app/bin/app $BIN/app
ENTRYPOINT ["app"]
CMD ["run"]
