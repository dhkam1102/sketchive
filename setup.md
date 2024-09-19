
# Sketchive: Real-Time Whiteboard Collaboration Platform - Setup Guide

This guide will help you set up the development environment for the **Sketchive** project on your local machine.

## Prerequisites

Make sure you have the following installed before proceeding:

- **Node.js**: v18.20.4 (or newer)
  - Install using [nvm](https://github.com/nvm-sh/nvm) (Node Version Manager):
    ```bash
    nvm install 18
    nvm use 18
    ```
- **npm**: v10.8.3 (automatically installed with Node.js)

## Step-by-Step Setup

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/your-username/sketchive.git
cd sketchive
```

### 2. Create the React App

We have already created the frontend React app, but in case you are setting it up from scratch:

```bash
npx create-react-app frontend
cd frontend
```

### 3. Install Project Dependencies

Install the necessary libraries for React, Konva, and Tailwind CSS:

```bash
npm install fabric konva react-konva tailwindcss
```

### 4. Initialize Tailwind CSS

Initialize Tailwind CSS for styling:

```bash
npx tailwindcss init
```

### 5. Start the Development Server

To run the React development server, use:

```bash
npm start
```

This will start the application in development mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

## Optional: Fix Vulnerabilities

After installing dependencies, you might notice some vulnerabilities. You can fix them using the following command:

```bash
npm audit fix --force
```

Note that this may introduce breaking changes, so it's recommended to review the changes after running this command.

## Versions of Tools

Here are the exact versions of the tools used to set up this project:

- **Node.js**: `v18.20.4`
- **npm**: `v10.8.3`
- **React**: `^18.x.x`
- **Fabric.js**: `^5.x.x`
- **Konva.js**: `^8.x.x`
- **Tailwind CSS**: `^3.x.x`

## Additional Resources

- [Node.js](https://nodejs.org/en/)
- [React.js](https://reactjs.org/)
- [Fabric.js](http://fabricjs.com/)
- [Konva.js](https://konvajs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
