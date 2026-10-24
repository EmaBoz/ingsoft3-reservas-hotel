# ---- Etapa 1: build ----
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json .npmrc ./
RUN npm ci
COPY . .
RUN npm run build

# ---- Etapa 2: tests (parte del build: ya tiene Node y las dependencias) ----
FROM build AS test
ENTRYPOINT ["npm", "run", "test:ci"]

# ---- Etapa 3: nginx sirve los estáticos ----
FROM nginx:alpine
COPY default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html
# Valores por defecto: los del compose. Cada entorno los pisa con sus variables.
ENV USERS_API_URL=http://users-api:8080 \
    ROOMS_API_URL=http://rooms-api:8080 \
    RESERVATIONS_API_URL=http://reservations-api:8080 \
    SEARCH_API_URL=http://search-api:8083 \
    DNS_RESOLVER=127.0.0.11
EXPOSE 80
