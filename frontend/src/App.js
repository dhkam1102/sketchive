import logo from './logo.svg';
import './App.css';
import React from 'react';
import Whiteboard from './Whiteboard';

function App() {
  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Real-Time Whiteboard</h1>
      </header>
      <main className="whiteboard">
        {/* Whiteboard component will go here */}
        <Whiteboard />
      </main>
      <footer className="app-footer">
        {/* Footer or other information */}
      </footer>
    </div>
  );
}


export default App;
