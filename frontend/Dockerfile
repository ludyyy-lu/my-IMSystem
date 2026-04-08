# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
# Point API and WS at the deployed services
ARG VITE_API_URL=http://api-gateway:8090
ARG VITE_WS_URL=ws://ws-gateway:8888
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_WS_URL=${VITE_WS_URL}
RUN npm run build

# Serve stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# nginx config: serve SPA and proxy /api and /ws to backend services
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
