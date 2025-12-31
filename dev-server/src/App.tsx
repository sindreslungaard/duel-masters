import { useState } from "react";
import "./App.css";
import { SignJWT } from "jose";
import { Duel } from "../../duel-interface";

const DUEL_TOKEN_SECRET = new TextEncoder().encode("duel-secret");

function App() {
  const [loading, setLoading] = useState(false);
  const [duel, setDuel] = useState<any>(null);
  const [hostDuelToken, setHostDuelToken] = useState("");
  const [guestDuelToken, setGuestDuelToken] = useState("");
  const [error, setError] = useState<string | null>(null);

  const createMatch = async () => {
    setLoading(true);
    setError(null);

    const payload = {
      hostId: "1",
      hostDeck: "",
      guestId: "2",
      guestDeck: "",
      name: "Test Match",
      visibility: "public",
      format: "regular",
    };

    try {
      const res = await fetch("/api/match", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }

      const data = await res.json();

      const hostToken = await new SignJWT({ id: "1", username: "Player1" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);
      const guestToken = await new SignJWT({ id: "2", username: "Player2" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);

      setHostDuelToken(hostToken);
      setGuestDuelToken(guestToken);
      setDuel(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button
        onClick={createMatch}
        disabled={loading}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {loading ? "Creating..." : "Create Match"}
      </button>

      {error && <p className="text-red-600 mt-2">Error: {error}</p>}
      {duel && (
        <div className="w-full h-screen flex gap-2">
          <Duel
            hostUrl="ws://localhost:3000"
            duelId={duel.id}
            duelToken={hostDuelToken}
          ></Duel>

          <Duel
            hostUrl="ws://localhost:3000"
            duelId={duel.id}
            duelToken={guestDuelToken}
          ></Duel>
        </div>
      )}
    </>
  );
}

export default App;
