FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
# Prefer a clean CI install; if it fails (lockfile out of sync) fall back to an install that
# relaxes peer deps so the build can continue in CI/docker environments.
RUN npm ci || npm install --legacy-peer-deps
COPY . .
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
