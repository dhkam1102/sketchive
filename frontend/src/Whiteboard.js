import React, { useRef, useEffect, useState } from 'react';
import './Whiteboard.css';
import { getWhiteboard, addStroke, getStrokesHistoryByWhiteboard, clearWhiteboard } from './api';

function Whiteboard() {
  const canvasRef = useRef(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [whiteboard, setWhiteboard] = useState(null);
  const [canvasSize, setCanvasSize] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const loadWhiteboard = async () => {
      try {
        const savedWhiteboard = await getWhiteboard(1); // Assuming whiteboard with id=1
        setWhiteboard(savedWhiteboard);
  
        // Fetch the stroke history for the whiteboard
        const strokeHistory = await getStrokesHistoryByWhiteboard(1);
  
        // Log strokeHistory to check what the API returns
        console.log("Stroke History:", strokeHistory);
  
        // Ensure strokeHistory is valid before using it
        if (Array.isArray(strokeHistory) && strokeHistory.length > 0) {
          strokeHistory.forEach(stroke => {
            drawStroke(stroke);
          });
        } else {
          console.log("No strokes found or strokeHistory is not an array.");
        }
      } catch (error) {
        console.error("Failed to load whiteboard or strokes:", error);
      }
    };
  
    loadWhiteboard();
  }, []);
  

  useEffect(() => {
    const updateCanvasSize = () => {
      setCanvasSize({
        width: window.innerWidth,
        height: window.innerHeight * 0.7
      });
    };

    updateCanvasSize();
    window.addEventListener('resize', updateCanvasSize);

    return () => {
      window.removeEventListener('resize', updateCanvasSize);
    };
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    canvas.width = canvasSize.width;
    canvas.height = canvasSize.height;

    const ctx = canvas.getContext('2d');
    ctx.lineCap = 'round';

    // Redraw existing strokes when canvas size changes
    if (whiteboard) {
      getStrokesHistoryByWhiteboard(whiteboard.id).then(strokeHistory => {
        strokeHistory.forEach(stroke => {
          drawStroke(stroke);
        });
      });
    }
  }, [canvasSize, whiteboard]);

  const drawStroke = (stroke) => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (stroke.path && stroke.path.length > 0) {
      ctx.beginPath();
      ctx.moveTo(stroke.path[0].x, stroke.path[0].y);
      for (let i = 1; i < stroke.path.length; i++) {
        ctx.lineTo(stroke.path[i].x, stroke.path[i].y);
      }
      ctx.strokeStyle = stroke.color;
      ctx.lineWidth = stroke.width;
      ctx.stroke();
      ctx.closePath();
    }
  };

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
  
    const ctx = canvas.getContext('2d');
  
    let currentPath = []; // Array to store all points of the current stroke
  
    const startDrawing = (event) => {
      if (!whiteboard || !whiteboard.id) {
        console.error("Whiteboard is not loaded properly");
        return;
      }
      const { offsetX, offsetY } = event;
      currentPath = [{ x: offsetX, y: offsetY }]; // Start with the initial point
      ctx.beginPath();
      ctx.moveTo(offsetX, offsetY);
      setIsDrawing(true);
    };
    
  
    const draw = (event) => {
      if (!isDrawing) return;
      const { offsetX, offsetY } = event;
      currentPath.push({ x: offsetX, y: offsetY }); // Collect points as user draws
      ctx.lineTo(offsetX, offsetY);
      ctx.stroke();
    };
  
    const stopDrawing = async (event) => {
      if (!isDrawing) return;
      setIsDrawing(false);
      ctx.closePath();
    
      if (!whiteboard || !whiteboard.id) {
        console.error("Whiteboard is not set or does not have a valid ID");
        return;
      }
    
      // Create strokeData with the collected path
      const strokeData = {
        whiteboardID: whiteboard.id,
        ownerID: 1, // Replace with actual user ID if needed
        path: currentPath,
        color: ctx.strokeStyle,
        width: ctx.lineWidth,
      };
    
      console.log("Attempting to add stroke:", strokeData);
    
      try {
        await addStroke(strokeData);
        console.log("Stroke added:", strokeData);
      } catch (error) {
        console.error("Failed to add stroke:", error);
      }
    };
    
  
    // Add event listeners for drawing
    canvas.addEventListener('mousedown', startDrawing);
    canvas.addEventListener('mousemove', draw);
    canvas.addEventListener('mouseup', stopDrawing);
    canvas.addEventListener('mouseleave', stopDrawing);
  
    // Clean up event listeners
    return () => {
      canvas.removeEventListener('mousedown', startDrawing);
      canvas.removeEventListener('mousemove', draw);
      canvas.removeEventListener('mouseup', stopDrawing);
      canvas.removeEventListener('mouseleave', stopDrawing);
    };
  }, [isDrawing, whiteboard]);

  const handleClearBoard = async () => {
    if (!whiteboard || !whiteboard.id) {
      console.error("Whiteboard is not loaded properly");
      return;
    }

    try {
      await clearWhiteboard(whiteboard.id);
      console.log("Whiteboard cleared successfully");

      // Clear the canvas
      const canvas = canvasRef.current;
      const ctx = canvas.getContext('2d');
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    } catch (error) {
      console.error("Failed to clear whiteboard:", error);
    }
  };

  return (
    <div className="whiteboard-container">
      <canvas ref={canvasRef} className="whiteboard-canvas"></canvas>
      <button onClick={handleClearBoard} className="clear-board-button">Clear Board</button>
    </div>
  );
}

export default Whiteboard;