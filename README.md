# Project to Learn GEM-R Stack 

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go](https://img.shields.io/badge/Backend-Go_1.21+-00ADD8.svg?logo=go&logoColor=white)
![Expo](https://img.shields.io/badge/Mobile-Expo_50+-000020.svg?logo=expo&logoColor=white)
![MongoDB](https://img.shields.io/badge/Database-MongoDB-47A248.svg?logo=mongodb&logoColor=white)
![React Native](https://img.shields.io/badge/UI-React_Native-61DAFB.svg?logo=react&logoColor=black)

Welcome! This repository serves as a practical demonstration of my full-stack development skills, specifically showcasing the **GEM-R Stack**. It was built to understand and implement fundamental CRUD (Create, Read, Update, Delete) operations in a modern, distributed application architecture.

## 🚀 The GEM-R Stack

This project leverages a powerful combination of technologies:

*   **G**o (Golang): High-performance backend API handling business logic and request processing.
*   **E**xpo: React Native framework for building a universal, native-feeling mobile application.
*   **M**ongoDB: Flexible NoSQL database for efficient data storage and retrieval.
*   **R**eact Native: Component-based UI library for crafting the user interface.

## 🎯 Project Goals

*   **Demonstrate Proficiency**: Showcasing the ability to integrate four distinct technologies into a cohesive application.
*   **Backend Mastery**: Implementing RESTful services with Go.
*   **Mobile Development**: Building responsive mobile interfaces with Expo and React Native.
*   **Database Management**: Handling data persistence with MongoDB.

## 🛠️ Installation & Walkthrough

If you'd like to run this project locally to see it in action, follow these steps.

### Prerequisites

*   [Go](https://go.dev/dl/) (v1.21+)
*   [Node.js](https://nodejs.org/) (LTS)
*   [MongoDB](https://www.mongodb.com/try/download/community) (Local or Atlas)

### 1. Environment & Database
Ensure your MongoDB instance is running.
```powershell
# Verify MongoDB service (Windows)
get-service MongoDB
```

### 2. Backend Setup (Go)
The backend service connects to MongoDB and exposes API endpoints.

```bash
cd backend
# 1. Create .env file for configuration
cp .env.example .env
# (Windows Command Prompt: copy .env.example .env)

# 2. Install dependencies
go mod tidy

# 3. Start the server
go run main.go
```

### 3. Frontend Setup (Expo)
The frontend application connects to the Go backend.

```bash
cd frontend
# 1. Configure API URL
cp .env.example .env
# (Windows Command Prompt: copy .env.example .env)

# 2. Install Node dependencies
npm install

# 3. Start the Expo development server
npx expo start
```

**To Run on Your Device:**
*   Download the **Expo Go** app on iOS or Android.
*   Scan the QR code displayed in the terminal.

## 🚀 Going to Production
Ready to deploy? Check out our detailed [Deployment Guide](DEPLOYMENT.md) for instructions on:
*   ☁️ **Database**: Setting up MongoDB Atlas.
*   🐳 **Backend**: Dockerizing the Go API.
*   📱 **Frontend**: Building for App Stores (EAS) and Web (Vercel).

## 🔮 Future Improvements

*   authentication (JWT)
*   Enhanced error handling and validation
*   Unit and Integration tests

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2026 Deepak Thomas


