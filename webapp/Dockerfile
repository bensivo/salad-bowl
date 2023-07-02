from node:16-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .
RUN npm run build

EXPOSE 3000

# TODO: create a multi-stage dockerfile with a separate build image, to reduce final image size
CMD ["npx", "serve@latest", "out"]