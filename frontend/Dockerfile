# File: financial-comparison-app/frontend/Dockerfile

# Build stage
FROM --platform=$BUILDPLATFORM node:20.15.0-alpine as build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Set the API URL environemnt variable
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL

RUN npm run build

# Production stage
FROM --platform=$TARGETPLATFORM nginx:stable-alpine

COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]