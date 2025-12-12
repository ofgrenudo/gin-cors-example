FROM node:20-alpine
LABEL authors="jwintersbro"

WORKDIR /opt/gin-cors-example/frontend

COPY package*.json ./
RUN npm ci

COPY . .

EXPOSE 3000
CMD ["npm", "start"]
