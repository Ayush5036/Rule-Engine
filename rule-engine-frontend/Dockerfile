# Dockerfile for React frontend

# Use an official Node.js image to build the React app
FROM node:16-alpine AS build

# Set working directory inside the container
WORKDIR /app

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React application for production
RUN npm run build

# Step 2: Serve the React build with a lightweight web server
FROM nginx:alpine

# Copy the built React app to the Nginx web server
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 3000 for the frontend
EXPOSE 3000

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
