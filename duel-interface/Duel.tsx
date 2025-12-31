import { useDuel } from "./useDuel";

interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
}

export function Duel({ duelId, duelToken, hostUrl }: DuelProps) {
  const { connected, error, messages, gameState, send, sendJoinMatch } =
    useDuel({
      hostUrl,
      duelId,
      duelToken,
    });

  return (
    <div className="w-full h-full bg-gray-900 text-white p-4">
      <div className="mb-4">
        <h2 className="text-2xl font-bold">Duel Interface</h2>
        <p>Duel ID: {duelId}</p>
        <p>Status: {connected ? "🟢 Connected" : "🔴 Disconnected"}</p>
        {error && <p className="text-red-500">Error: {error}</p>}
      </div>

      <div className="mb-4">
        <button
          onClick={sendJoinMatch}
          disabled={!connected}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-500"
        >
          Join Match
        </button>
      </div>

      <div className="mb-4">
        <h3 className="text-xl font-semibold mb-2">Game State</h3>
        <pre className="bg-gray-800 p-2 rounded overflow-auto">
          {gameState ? JSON.stringify(gameState, null, 2) : "No game state yet"}
        </pre>
      </div>

      <div>
        <h3 className="text-xl font-semibold mb-2">
          Messages ({messages.length})
        </h3>
        <div className="bg-gray-800 p-2 rounded h-64 overflow-auto">
          {messages.map((msg, index) => (
            <div key={index} className="mb-2 border-b border-gray-700 pb-2">
              <pre className="text-sm">{JSON.stringify(msg, null, 2)}</pre>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
