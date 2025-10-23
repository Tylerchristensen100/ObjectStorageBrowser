FROM node:22.20.0-alpine AS builder

ARG VITE_BASE_URL
ENV VITE_BASE_URL=$VITE_BASE_URL

WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

RUN echo "VITE_BASE_URL is $VITE_BASE_URL"



FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
