FROM node
WORKDIR /e2e

COPY package*.json ./
RUN npm ci

COPY . .
CMD ["npx", "jest"]
