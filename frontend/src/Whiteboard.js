import React, { useRef, useEffect, useState } from 'react';
import './Whiteboard.css';
import { getWhiteboard, createWhiteboard, updateWhiteboard } from './api'; // Import necessary API functions

function Whiteboard() {
  const canvasRef = useRef(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [whiteboard, setWhiteboard] = useState(null); // State to store whiteboard data

  useEffect(() => {
    const loadOrCreateWhiteboard = async () => {
      try {
        // Try to load whiteboard with id = 1
        const savedWhiteboard = await getWhiteboard(1);
        setWhiteboard(savedWhiteboard);
        drawSavedWhiteboard(savedWhiteboard.currentState); // Render the saved strokes
      } catch (error) {
        console.error("Whiteboard not found. Creating a new one...");
        // If no whiteboard is found, create a new one
        const newBoard = await createWhiteboard();
        setWhiteboard(newBoard);
      }
    };

    loadOrCreateWhiteboard(); // Load or create whiteboard on component load

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');

    // Set canvas size and drawing styles
    canvas.width = window.innerWidth * 0.8;
    canvas.height = window.innerHeight * 0.6;
    ctx.lineWidth = 5;
    ctx.lineCap = 'round';
    ctx.strokeStyle = '#000';

    // Mouse down event to start drawing
    const startDrawing = (event) => {
      if (event.type === 'mousedown') {
        const { offsetX, offsetY } = event;
        ctx.beginPath();
        ctx.moveTo(offsetX, offsetY);
        setIsDrawing(true);
      }
    };

    // Mouse move event to draw
    const draw = (event) => {
      if (!isDrawing) return;
      if (event.type === 'mousemove') {
        const { offsetX, offsetY } = event;
        ctx.lineTo(offsetX, offsetY);
        ctx.stroke();
      }
    };

    // Mouse up event to stop drawing
    const stopDrawing = async () => {
      setIsDrawing(false);
      ctx.closePath();

      // After drawing is done, save the whiteboard state to the backend (whiteboard id = 1)
      if (whiteboard) {
        const canvasData = canvas.toDataURL(); // Convert canvas to an image
        const updatedWhiteboard = {
          ...whiteboard,
          currentState: canvasData, // Store the image data as the current state
        };
        await updateWhiteboard(1, updatedWhiteboard); // Save whiteboard with id = 1
      }
    };

    // Event listeners
    canvas.addEventListener('mousedown', startDrawing);
    canvas.addEventListener('mousemove', draw);
    canvas.addEventListener('mouseup', stopDrawing);
    canvas.addEventListener('mouseleave', stopDrawing);

    // Cleanup event listeners on component unmount
    return () => {
      canvas.removeEventListener('mousedown', startDrawing);
      canvas.removeEventListener('mousemove', draw);
      canvas.removeEventListener('mouseup', stopDrawing);
      canvas.removeEventListener('mouseleave', stopDrawing);
    };
  }, [isDrawing, whiteboard]); // Dependencies include whiteboard and isDrawing

  // Function to clear the canvas
  const clearCanvas = () => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);  // Clear the entire canvas
  };

  // Function to render the saved whiteboard (if available)
  const drawSavedWhiteboard = (savedState) => {
    if (!savedState) return;
    const img = new Image();
    img.src = savedState; // Load the image data
    img.onload = () => {
      const canvas = canvasRef.current;
      const ctx = canvas.getContext('2d');
      ctx.drawImage(img, 0, 0); // Draw the saved whiteboard on the canvas
    };
  };

  return (
    <div className="whiteboard-container">
      <canvas ref={canvasRef} className="whiteboard-canvas"></canvas>
      <button onClick={clearCanvas} className="clear-btn">Clear Canvas</button>
    </div>
  );
}

export default Whiteboard;
