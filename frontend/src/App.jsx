import { useState, useEffect } from "react"
import JoinScreen from "./components/JoinScreen"
import GameScreen from "./components/GameScreen"
import { useWebSocket } from "./hooks/useWebSocket"

export default function App() {
  const [gameState, setGameState] = useState(null)

  const send = useWebSocket((data) => {
    setGameState(data)

    localStorage.setItem("gameID", data.gameID)
    localStorage.setItem("playerID", data.playerID)

    if (data.message) {
      setTimeout(() => {
        setGameState(prev => ({ ...prev, message: null }))
      }, 3000)
    }
  })

  //  Auto reconnect
  useEffect(() => {
    const gameID = localStorage.getItem("gameID")
    const playerID = localStorage.getItem("playerID")

    if (gameID && playerID) {
      send({ type: "reconnect", gameID, playerID })
    }
  }, [])

  //  JOIN FLOW
  if (!gameState) {
    return (
      <JoinScreen
        onJoin={({ username, gameId, mode }) =>
          send({
            type: "join",
            username,
            gameID: gameId || null,
            mode,
          })
        }
      />
    )
  }

  return <GameScreen state={gameState} send={send} />
}
