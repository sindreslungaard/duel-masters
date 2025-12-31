import { useState } from "react";
import "./App.css";

function App() {
  const [loading, setLoading] = useState(false);
  const [duel, setDuel] = useState<any>(null);
  const [duelToken, setDuelToken] = useState<string | null>(null);
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
      setDuel(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <p className="bg-red-500">hello</p>
      <button
        onClick={createMatch}
        disabled={loading}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {loading ? "Creating..." : "Create Match"}
      </button>

      {error && <p className="text-red-600 mt-2">Error: {error}</p>}
      {duel && (
        <pre className="mt-2 p-2 bg-gray-100 rounded">
          {JSON.stringify(duel, null, 2)}
        </pre>
      )}
    </>
  );
}

export default App;
