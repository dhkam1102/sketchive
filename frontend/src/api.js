// src/api.js

// will have to replace localhost for deployment
const API_BASE_URL = "http://localhost:8080"; 

export async function createWhiteboard() {
  const response = await fetch(`${API_BASE_URL}/whiteboards`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error("Failed to create whiteboard");
  }

  return await response.json(); // Return the whiteboard data
}

export async function getWhiteboard(id) {
  const response = await fetch(`${API_BASE_URL}/whiteboards?id=${id}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to get whiteboard");
  }

  return await response.json();
}

export async function updateWhiteboard(id, updatedData) {
  const response = await fetch(`${API_BASE_URL}/whiteboards?id=${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updatedData),
  });

  if (!response.ok) {
    throw new Error("Failed to update whiteboard");
  }

  return await response.json();
}

export async function deleteWhiteboard(id) {
  const response = await fetch(`${API_BASE_URL}/whiteboards?id=${id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Failed to delete whiteboard");
  }

  return await response.json();
}
