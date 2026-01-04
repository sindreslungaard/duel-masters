import { useState, useEffect } from "react";
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
  const [activePlayer, setActivePlayer] = useState<"host" | "guest">("host");
  const [cards, setCards] = useState<{ uid: string; name: string }[]>([]);

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
      const [matchRes, cardsRes] = await Promise.all([
        fetch("/api/match", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        }),
        fetch("/api/cards"),
      ]);

      if (!matchRes.ok) {
        throw new Error(`HTTP error! status: ${matchRes.status}`);
      }

      if (!cardsRes.ok) {
        throw new Error(`HTTP error! status: ${cardsRes.status}`);
      }

      const match = await matchRes.json();
      const cards = await cardsRes.json();
      setCards(cards);

      const hostToken = await new SignJWT({ id: "1", username: "Player1" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);
      const guestToken = await new SignJWT({ id: "2", username: "Player2" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);

      setHostDuelToken(hostToken);
      setGuestDuelToken(guestToken);
      setDuel(match);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    createMatch();
  }, []);

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      if (event.key === "1") {
        setActivePlayer("host");
      } else if (event.key === "2") {
        setActivePlayer("guest");
      }
    };

    window.addEventListener("keydown", handleKeyPress);
    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, []);

  return (
    <>
      {loading && (
        <p className="text-white mt-20 text-center">Loading match...</p>
      )}
      {error && (
        <p className="text-red-600 mt-20 text-center">Error: {error}</p>
      )}

      {duel && (
        <div className="w-full h-screen">
          <div className={activePlayer === "host" ? "block" : "hidden"}>
            <Duel
              hostUrl="ws://localhost:3000"
              duelId={duel.id}
              duelToken={hostDuelToken}
              devTools={{
                cards,
                activePlayer,
                onPlayerSwitch: setActivePlayer,
              }}
            />
          </div>

          <div className={activePlayer === "guest" ? "block" : "hidden"}>
            <Duel
              hostUrl="ws://localhost:3000"
              duelId={duel.id}
              duelToken={guestDuelToken}
              devTools={{
                cards,
                activePlayer,
                onPlayerSwitch: setActivePlayer,
              }}
            />
          </div>
        </div>
      )}
    </>
  );
}

export default App;
