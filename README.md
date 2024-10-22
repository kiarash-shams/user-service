# OAuth2 Server Microservice

## Overview
The **OAuth2 Server Microservice** is an independent authentication service that implements OAuth2 for secure user authentication and authorization. It handles user registration via email and mobile phone, manages user information securely using Vault, and supports KYC (Know Your Customer) processes.

## Features
- **User Registration**: Allows users to register with email or mobile phone.
- **Vault Integration**: Secures user data by storing sensitive information using Vault encryption.
- **KYC Labeling**: Manages KYC verification stages, including email, phone, and profile completion.
- **Mobile Number Registration**: Separate service to handle phone number verification.
- **OAuth2 Authentication**: Provides OAuth2-compliant authentication for external applications.
- **Performance & Scalability**: Optimized for handling large numbers of users and requests with low latency.

## Architecture
The microservice is built using **Go** and the **Gin** framework. It interacts with various components:
- **PostgreSQL**: For storing user data and related information.
- **Vault**: To store sensitive information securely (e.g., user credentials, personal information).
- **Redis**: Used for caching and managing user sessions.
- **Elasticsearch**: Integrated for efficient logging and search capabilities.
- **Prometheus**: For monitoring application metrics and performance.