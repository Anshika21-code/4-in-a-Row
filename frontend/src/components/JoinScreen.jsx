import { useState } from "react"

export default function JoinScreen({ onJoin }) {
  const [username, setUsername] = useState("")
  const [gameId, setGameId] = useState("")
  const [mode, setMode] = useState("human") // human | bot

  return (
    <div className="center">
      <h2>4 in a Row</h2>

      <input
        placeholder="Username"
        value={username}
        onChange={e => setUsername(e.target.value)}
      />

      <input
        placeholder="Game ID (optional)"
        value={gameId}
        onChange={e => setGameId(e.target.value)}
      />

      {/*  MODE SELECT */}
      <div>
        <label>
          <input
            type="radio"
            value="human"
            checked={mode === "human"}
            onChange={() => setMode("human")}
          />
          Play vs Human
        </label>
        <br />
        <label>
          <input
            type="radio"
            value="bot"
            checked={mode === "bot"}
            onChange={() => setMode("bot")}
          />
          Play vs Bot 
        </label>
      </div>

      <button
        onClick={() => onJoin({ username, gameId, mode })}
        disabled={!username}
      >
        Create / Join
      </button>
    </div>
  )
}
