FROM node:22-alpine AS builder

WORKDIR /app

# Copy app files
COPY . .

# Install http-server for serving static files
RUN npm install -g http-server

# Use a lightweight web server
FROM nginx:alpine

# Copy files from builder stage
COPY --from=builder /app /usr/share/nginx/html

# Copy custom nginx config if needed
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Default command uses nginx's default command
