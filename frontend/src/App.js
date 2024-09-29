import React from 'react';
import './App.css';
import Whiteboard from './Whiteboard';
import sketchiveLogo from './assets/logo.png';

function App() {
  return (
    <div className="app-container">
      {/* <header className="app-header">
        <img src={sketchiveLogo} alt="Sketchive Logo" className="sketchive-logo" />
      </header> */}
      <main className="app-main">
        <Whiteboard />
      </main>
    </div>
  );
}

export default App;