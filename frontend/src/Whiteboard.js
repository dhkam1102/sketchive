import React, { useRef, useEffect, useState } from 'react';
import './Whiteboard.css';

function Whiteboard() {
  const canvasRef = useRef(null);
  const [isDrawing, setIsDrawing] = useState(false);

  useEffect(() => {
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
      // Ensure this is a mouse event
      if (event.type === 'mousedown') {
        const { offsetX, offsetY } = event;  // Use event directly for mouse events
        ctx.beginPath();
        ctx.moveTo(offsetX, offsetY);
        setIsDrawing(true);
      }
    };

    // Mouse move event to draw
    const draw = (event) => {
      if (!isDrawing) return;
      if (event.type === 'mousemove') {
        const { offsetX, offsetY } = event;  // Use event directly for mouse events
        ctx.lineTo(offsetX, offsetY);
        ctx.stroke();
      }
    };

    // Mouse up event to stop drawing
    const stopDrawing = () => {
      setIsDrawing(false);
      ctx.closePath();
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
  }, [isDrawing]);

  // Function to clear the canvas
  const clearCanvas = () => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);  // Clear the entire canvas
  };

  return (
    <div className="whiteboard-container">
      <canvas ref={canvasRef} className="whiteboard-canvas"></canvas>
      <button onClick={clearCanvas} className="clear-btn">Clear Canvas</button>
    </div>
  );
}

export default Whiteboard;
