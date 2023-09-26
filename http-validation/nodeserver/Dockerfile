FROM node:20-slim

WORKDIR /app
COPY package.json .
COPY nodeserver/index.js .
RUN npm i --production

ENV HOST_ADDR="0.0.0.0"
CMD ["node", "index.js"]
