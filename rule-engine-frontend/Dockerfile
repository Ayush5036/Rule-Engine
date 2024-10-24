# Use Node.js image for building the frontend
FROM node:16 AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app
RUN npm run build

# Use NGINX to serve the static files
FROM nginx:alpine

# Copy built files to NGINX html directory
COPY --from=builder /app/build /usr/share/nginx/html

# Expose the port the app runs on
EXPOSE 3000

# Command to run NGINX
CMD ["nginx", "-g", "daemon off;"]
