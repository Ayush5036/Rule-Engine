# Rule Engine Application

This is a **Rule Engine** application built with **Golang** in the backend and **React** in the frontend. The application uses MongoDB for persistent storage and Docker for containerization. The Rule Engine allows you to create, combine, and evaluate complex conditional rules, which can be used to determine user eligibility based on different attributes.

## Table of Contents
- [Project Overview](#project-overview)
- [Technologies Used](#technologies-used)
- [Architecture](#architecture)
- [Setup and Installation](#setup-and-installation)
- [API Endpoints](#api-endpoints)
- [How to Run the Project](#how-to-run-the-project)
- [Frontend Features](#frontend-features)
- [Backend Features](#backend-features)
- [Usage](#usage)

## Project Overview
This Rule Engine provides a flexible system for evaluating business rules and determining user eligibility. It works by representing rules using an **Abstract Syntax Tree (AST)**, allowing for rule creation, combination, and evaluation. The rules are stored in MongoDB, and the backend is built with **Golang**, while the frontend is developed using **React**.

## Technologies Used

- **Backend**: Golang (Go)
- **Frontend**: React.js
- **Database**: MongoDB
- **Containerization**: Docker, Docker Compose
- **Storage**: MongoDB (NoSQL)
- **Other Tools**: AST for rule representation, Docker for environment setup

## Architecture

The project uses a **3-tier architecture**:

1. **Frontend**: Developed using React.js, allowing users to interact with the Rule Engine through a user-friendly interface. Users can create, view, combine, and evaluate rules.
2. **API**: A REST API created with Golang, responsible for handling requests related to rule creation, rule combination, and evaluation. It communicates with the MongoDB database and manages rule storage.
3. **Database**: MongoDB is used to persist rule data. Rules are stored in collections and queried as needed.

### Design Patterns

- **AST (Abstract Syntax Tree)**: The rules are represented as an AST, which is used for evaluating complex logical expressions.
- **Microservice-Oriented Design**: Each part of the application (frontend, backend) is independently containerized using Docker.
- **RESTful API**: The backend exposes a set of RESTful API endpoints for managing rules, combining them, and evaluating rules based on query data.

## Setup and Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo-url/rule-engine.git
   cd rule-engine
   ```

2. **Environment Setup**:
   Ensure that you have Docker and Docker Compose installed on your machine.

   - Docker: [Install Docker](https://docs.docker.com/get-docker/)
   - Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

3. **Build the Docker Images**:
   Inside the project directory, build the Docker images for the frontend, backend, and MongoDB using:

   ```bash
   docker-compose build
   ```

4. **Start the Services**:
   Start the services using Docker Compose:

   ```bash
   docker-compose up
   ```

5. **Access the Application**:
   - Frontend: The React frontend will be running on `http://localhost:3000`
   - Backend: The Go API will be running on `http://localhost:8080`
   - MongoDB: MongoDB will be accessible on `localhost:27017` (but handled internally by Docker)

## API Endpoints

The following API endpoints are available:

1. **Create a Rule** (POST `/api/rules`)
   - Request Body: `{ "id": "<rule_id>", "expression": "<rule_expression>" }`
   - Response: Success or error message with rule details.

2. **Combine Rules** (POST `/api/rules/combine`)
   - Request Body: `{ "rule_ids": "<comma_separated_rule_ids>" }`
   - Response: Combined rule expression.

3. **Evaluate Rules** (POST `/api/rules/evaluate`)
   - Request Body: `{ "query_data": { <key-value attributes> } }`
   - Response: Evaluation result based on the combined rule.

4. **Get All Rules** (GET `/api/rules`)
   - Fetch all stored rules from MongoDB.

5. **Clean Database** (DELETE `/api/rules/clean`)
   - Remove all rules from MongoDB.

## How to Run the Project

1. **Step 1**: Build the Docker images:

   ```bash
   docker-compose build
   ```

2. **Step 2**: Start the services:

   ```bash
   docker-compose up
   ```

   This will spin up three services:
   - **Frontend**: React application at `http://localhost:3000`
   - **Backend**: Golang-based API at `http://localhost:8080`
   - **MongoDB**: Database running on `localhost:27017` (internal by Docker)

3. **Step 3**: Open the frontend in your browser and start creating and evaluating rules!

## Frontend Features

- **Create Rule**: Add new rules with ID and expression.
- **View Rules**: Display all existing rules stored in MongoDB.
- **Combine Rules**: Combine multiple rules to create a complex expression.
- **Evaluate Rule**: Evaluate the combined rule or individual rules based on input query data.
- **Responsive UI**: A responsive interface designed for user-friendly interaction with the Rule Engine.

## Backend Features

- **Create Rule API**: Endpoint to create and store rules in MongoDB.
- **Combine Rule API**: Combines multiple rules into a single expression.
- **Evaluate Rule API**: Evaluates the combined rule using AST and returns the result.
- **MongoDB Storage**: All rules are persisted in MongoDB.

## Usage

### Creating and Evaluating Rules

1. **Create a Rule**:
   - Use the frontend or send a POST request to `/api/rules` with the rule expression.
   
2. **Combine Rules**:
   - Combine multiple rules via the frontend or the `/api/rules/combine` endpoint.
   
3. **Evaluate Rules**:
   - Input a query to evaluate the combined rule expression. Results will be displayed based on user attributes and rule conditions.
