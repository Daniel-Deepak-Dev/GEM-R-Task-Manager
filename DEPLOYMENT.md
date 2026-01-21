# 🚀 Going to Production

This guide outlines the steps to take your GEM-R Stack application from development to a production environment.

## 1. Database (MongoDB Atlas)
For production, use a managed cloud database instead of a local instance.

1.  **Sign Up**: Create an account on [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
2.  **Create Cluster**: Deploy a free shared cluster (M0 Sandbox).
3.  **Get Connection String**:
    -   Click "Connect" -> "Drivers" -> "Go".
    -   Copy the connection string (e.g., `mongodb+srv://<username>:<password>@cluster.mongodb.net/...`).
4.  **Update Config**:
    -   In your production environment variables, set `MONGO_URI` to this new string.

## 2. Backend (Go & Docker)
Containerizing your Go application ensures it runs the same way everywhere.

### Dockerfile
We have included a `Dockerfile` in the `backend/` directory.

### Deploying
You can deploy this Docker container to any cloud provider:
*   **Render / Railway**: Connect your GitHub repo, point to the `backend` directory, and it will auto-detect the Dockerfile.
*   **AWS App Runner / Google Cloud Run**: Push your image to a registry and deploy.

**Environment Variables**:
Ensure you set `MONGO_URI` and `PORT` in your cloud provider's dashboard.

## 3. Frontend (Expo & React Native)

### A. Mobile App (iOS & Android)
Use **EAS Build** (Expo Application Services) to build your app for the App Store and Play Store.

1.  **Install EAS CLI**:
    ```bash
    npm install -g eas-cli
    ```
2.  **Login**:
    ```bash
    eas login
    ```
3.  **Configure**:
    ```bash
    eas build:configure
    ```
4.  **Build**:
    ```bash
    eas build --platform all
    ```

### B. Web App
You can host the web version of your app on Vercel or Netlify.

1.  **Export for Web**:
    ```bash
    npx expo export
    ```
    This creates a `dist` folder with static files.
2.  **Deploy**:
    -   Drag and drop the `dist` folder into Netlify Drop.
    -   OR connect Vercel to your GitHub and configure the build command (`npx expo export`) and output directory (`dist`).

**Environment Variables**:
Set `EXPO_PUBLIC_API_URL` to your **production backend URL** (e.g., `https://my-go-backend.onrender.com/api/tasks`).
