import React, { useRef, useEffect, useState } from 'react';
import './Whiteboard.css';
// import { getWhiteboard, getStrokesHistoryByWhiteboard } from './api.js';

import { getWhiteboard, addStroke, getStrokesHistoryByWhiteboard, clearWhiteboard, deleteStrokesByBoundingBox} from './api.js';

function Whiteboard() {
  const canvasRef = useRef(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [whiteboard, setWhiteboard] = useState(null);
  const [canvasSize, setCanvasSize] = useState({ width: 0, height: 0 });
  const [currentTool, setCurrentTool] = useState("pen");
  const [strokeStyle, setStrokeStyle] = useState("#000000");
  const [lineWidth, setLineWidth] = useState(2);

  // loading the whiteboard: getting stroke history and drawing
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
        if (strokeHistory && Array.isArray(strokeHistory) && strokeHistory.length > 0) {
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
  
  // Setting the canvas size
  useEffect(() => {
    const updateCanvasSize = () => {
      setCanvasSize({
        width: window.innerWidth,
        height: window.innerHeight * 0.7
      });
    };

    updateCanvasSize();
    // setting up a event listener so whiteboard will be updated when changing the window size
    window.addEventListener('resize', updateCanvasSize);

    return () => {
      window.removeEventListener('resize', updateCanvasSize);
    };
  }, []);

  // Redraw the stokes: when whiteboard is updated or canvas is resized
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
        // Ensure strokeHistory is valid before using it
        if (strokeHistory && Array.isArray(strokeHistory) && strokeHistory.length > 0) {
          strokeHistory.forEach(stroke => {
            drawStroke(stroke);
          });
        } else {
          console.log("No strokes found or strokeHistory is not an array.");
        }
      }).catch(error => {
        console.error("Failed to fetch strokes:", error);
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

  const calculateBoundingBox = (points) => {
    if (!points || points.length === 0) {
      console.error('No points provided to calculate bounding box');
      return null;
    }

    if (points.length === 1) {
      // If there's only one point, create a small bounding box around it
      const point = points[0];
      const padding = 1; // 1 pixel padding
      return {
        minX: point.x - padding,
        maxX: point.x + padding,
        minY: point.y - padding,
        maxY: point.y + padding
      };
    }

    let minX = points[0].x, maxX = points[0].x;
    let minY = points[0].y, maxY = points[0].y;
  
    points.forEach(point => {
      if (point && typeof point.x === 'number' && typeof point.y === 'number') {
        minX = Math.min(minX, point.x);
        maxX = Math.max(maxX, point.x);
        minY = Math.min(minY, point.y);
        maxY = Math.max(maxY, point.y);
      }
    });
  
    return { minX, maxX, minY, maxY };
  };

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.strokeStyle = strokeStyle;
    ctx.lineWidth = lineWidth;

    canvas.style.cursor = currentTool === "eraser" 
      ? `url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="${lineWidth * 2}" height="${lineWidth * 2}" viewBox="0 0 ${lineWidth * 2} ${lineWidth * 2}"><circle cx="${lineWidth}" cy="${lineWidth}" r="${lineWidth}" fill="rgba(128,128,128,0.5)"/></svg>') ${lineWidth} ${lineWidth}, auto`
      : "crosshair";


    let currentPath = [];

    const startDrawing = (event) => {
      if (!whiteboard || !whiteboard.id) {
        console.error("Whiteboard is not loaded properly");
        return;
      }
      const { offsetX, offsetY } = event;
      currentPath = [{ x: offsetX, y: offsetY }];
      ctx.beginPath();
      ctx.moveTo(offsetX, offsetY);
      setIsDrawing(true);
    };

    const draw = (event) => {
      if (!isDrawing) return;
      const { offsetX, offsetY } = event;
      currentPath.push({ x: offsetX, y: offsetY });

      if (currentTool === "pen") {
        ctx.lineTo(offsetX, offsetY);
        ctx.stroke();
      } else if (currentTool === "eraser") {
        ctx.save();
        ctx.globalCompositeOperation = 'destination-out';
        ctx.arc(offsetX, offsetY, lineWidth * 2, 0, Math.PI * 2, false);
        ctx.fill();
        ctx.restore();
      }
    };

    const stopDrawing = async () => {
      if (!isDrawing) return;
      setIsDrawing(false);
      ctx.closePath();

      if (!whiteboard || !whiteboard.id) {
        console.error("Whiteboard is not set or does not have a valid ID");
        return;
      }

      if (currentTool === "pen") {
        const strokeData = {
          whiteboardID: whiteboard.id,
          ownerID: 1,
          path: currentPath,
          color: strokeStyle,
          width: lineWidth,
        };
        try {
          await addStroke(strokeData);
          console.log("Stroke added:", strokeData);
        } catch (error) {
          console.error("Failed to add stroke:", error);
        }
      } else if (currentTool === "eraser") {
        const boundingBox = calculateBoundingBox(currentPath);
        try {
          await deleteStrokesByBoundingBox(whiteboard.id, boundingBox);
          console.log("Strokes deleted in bounding box:", boundingBox);
          // Redraw the canvas after erasing
          redrawCanvas();
        } catch (error) {
          console.error("Failed to delete strokes:", error);
        }
      } else {
        console.error("Failed to calculate bounding box for eraser");
      }
    };

    canvas.addEventListener('mousedown', startDrawing);
    canvas.addEventListener('mousemove', draw);
    canvas.addEventListener('mouseup', stopDrawing);
    canvas.addEventListener('mouseleave', stopDrawing);

    return () => {
      canvas.removeEventListener('mousedown', startDrawing);
      canvas.removeEventListener('mousemove', draw);
      canvas.removeEventListener('mouseup', stopDrawing);
      canvas.removeEventListener('mouseleave', stopDrawing);
    };
  }, [isDrawing, whiteboard, currentTool, strokeStyle, lineWidth]);

  const redrawCanvas = async () => {
    if (!whiteboard || !whiteboard.id) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    try {
      const strokeHistory = await getStrokesHistoryByWhiteboard(whiteboard.id);
      if (strokeHistory && Array.isArray(strokeHistory) && strokeHistory.length > 0) {
        strokeHistory.forEach(stroke => {
          drawStroke(stroke);
        });
      }
    } catch (error) {
      console.error("Failed to redraw canvas:", error);
    }
  };

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
      <div className="tool-selection">
        <button 
          className={`tool-button ${currentTool === "pen" ? "active" : ""}`}
          onClick={() => setCurrentTool("pen")}
        >
          Pen
        </button>
        <button 
          className={`tool-button ${currentTool === "eraser" ? "active" : ""}`}
          onClick={() => setCurrentTool("eraser")}
        >
          Eraser
        </button>
        <input
          type="color"
          value={strokeStyle}
          onChange={(e) => setStrokeStyle(e.target.value)}
          disabled={currentTool === "eraser"}
          className="color-picker"
        />
        <input
          type="range"
          min="1"
          max="20"
          value={lineWidth}
          onChange={(e) => setLineWidth(parseInt(e.target.value))}
          className="line-width-slider"
        />
      </div>
      <canvas ref={canvasRef} className="whiteboard-canvas"></canvas>
      <button onClick={handleClearBoard} className="clear-board-button">Clear Board</button>
    </div>
  );
}

export default Whiteboard;