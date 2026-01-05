import { useState } from "react";

export default function JoinScreen({ onJoin }) {
  const [username, setUsername] = useState("");
  const [mode, setMode] = useState("bot");

  return (
    <div className="join-page">
      <div className="join-card">
        <h1 className="join-title"> 4 in a Row</h1>
        <p className="join-subtitle">
          Challenge the bot or play with a friend
        </p>

        <input
          className="join-input"
          placeholder="Enter your username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />

        <div className="mode-group">
          <label className={`mode ${mode === "bot" ? "active" : ""}`}>
            <input
              type="radio"
              checked={mode === "bot"}
              onChange={() => setMode("bot")}
            />
             Play vs Bot
          </label>

          <label className={`mode ${mode === "human" ? "active" : ""}`}>
            <input
              type="radio"
              checked={mode === "human"}
              onChange={() => setMode("human")}
            />
             Play vs Human
          </label>
        </div>

        <button
          className="start-btn"
          disabled={!username}
          onClick={() => onJoin({ username, mode })}
        >
          Start Game
        </button>
      </div>
    </div>
  );
}
