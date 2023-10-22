FROM golang:1.21.3-bookworm as builder
COPY . /build
WORKDIR /build
RUN go build -o components/producers/typescript-eslint/eslint-wrapper/eslint-wrapper ./components/producers/typescript-eslint/eslint-wrapper/main.go

# TODO: replace w/ digest and auto create PRs to update.
FROM node:lts

WORKDIR /home/node/workspace
COPY --from=builder /build/components/producers/typescript-eslint/eslint-wrapper/eslintrc.js /home/node/workspace
COPY --from=builder /build/components/producers/typescript-eslint/eslint-wrapper/package.json /home/node/workspace
COPY --from=builder /build/components/producers/typescript-eslint/eslint-wrapper/eslint-wrapper /home/node/workspace
RUN npm uninstall --save bcrypt &&\
    npm install --save-dev \
    eslint \
    eslint-plugin-security@latest \
    eslint-plugin-xss \
    eslint-plugin-no-unsanitized \
    eslint-plugin-security-node

ENTRYPOINT [ "/home/node/workspace/eslint-wrapper"]
