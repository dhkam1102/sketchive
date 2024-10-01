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

export async function clearWhiteboard(whiteboardId) {
  const response = await fetch(`${API_BASE_URL}/whiteboards/clear?id=${whiteboardId}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Failed to clear whiteboard");
  }

  return await response.json();
}


// -----------------------------------
// Stroke API calls
// -----------------------------------

// Add a stroke to the whiteboard
export async function addStroke(strokeData) {
  // Example strokeData structure:
  // {
  //   whiteboardID: 1,
  //   ownerID: 1,
  //   path: [
  //     { x: 10, y: 20 },
  //     { x: 15, y: 25 },
  //     { x: 25, y: 35 }
  //   ],
  //   color: "#000000",
  //   width: 5,
  //   createdAt: "2024-09-26T15:04:05Z" // Optional
  // }

  // Ensure `path` is formatted correctly before sending
  console.log("Adding stroke with data:", strokeData);
  const response = await fetch(`${API_BASE_URL}/strokes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(strokeData), // Send complete stroke data in the request body
  });

  if (!response.ok) {
    const errorMessage = await response.text();
    throw new Error(`Failed to add stroke: ${response.status} - ${errorMessage}`);
  }

  return await response.json();
}

// Get stroke history for a whiteboard by its ID
export async function getStrokesHistoryByWhiteboard(whiteboardId) {
  const response = await fetch(`${API_BASE_URL}/strokes?id=${whiteboardId}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to get stroke history");
  }

  return await response.json();
}

// Delete strokes based on the bounding box (eraser action)
export async function deleteStrokesByBoundingBox(whiteboardID, boundingBox) {
  // Example boundingBox structure:
  // {
  //   minX: 10,
  //   maxX: 50,
  //   minY: 20,
  //   maxY: 60,
  // }

  const response = await fetch(`${API_BASE_URL}/strokes/delete`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ whiteboardID, ...boundingBox }),
  });

  if (!response.ok) {
    const errorMessage = await response.text();
    throw new Error(`Failed to delete strokes: ${response.status} - ${errorMessage}`);
  }

  return await response.json();
}