FROM node:lts

WORKDIR /home/node/workspace
COPY components/producers/typescript-eslint/eslint-wrapper/eslintrc.js /home/node/workspace
COPY components/producers/typescript-eslint/eslint-wrapper/package.json /home/node/workspace
COPY bin/components/producers/typescript-eslint/eslint-wrapper/eslint-wrapper-parser /home/node/workspace/

RUN npm uninstall --save bcrypt &&\
    npm install --save-dev \
    eslint \
    eslint-plugin-security@latest \
    eslint-plugin-xss \
    eslint-plugin-no-unsanitized \
    eslint-plugin-security-node

ENTRYPOINT [ "/home/node/workspace/eslint-wrapper-parser"]
