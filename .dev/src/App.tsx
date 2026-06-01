import { useState, useEffect } from "react";
import "./App.css";
import { SignJWT } from "jose";
import { Duel } from "../../packages/duel-interface";

const DUEL_TOKEN_SECRET = new TextEncoder().encode("duel-secret");
const DECK_SIZE = 40;
const UNIQUE_CARDS_PER_DECK = 10;

type CardSummary = {
  uid: string;
  name: string;
};

const buildDeck = (availableCards: CardSummary[], offset = 0): string[] => {
  if (availableCards.length === 0) {
    return [];
  }

  const uniqueCards = Array.from(
    { length: Math.min(UNIQUE_CARDS_PER_DECK, availableCards.length) },
    (_, index) => availableCards[(offset + index) % availableCards.length].uid,
  );

  const deck: string[] = [];
  while (deck.length < DECK_SIZE) {
    for (const uid of uniqueCards) {
      deck.push(uid);
      if (deck.length === DECK_SIZE) {
        break;
      }
    }
  }

  return deck;
};

function App() {
  const [loading, setLoading] = useState(false);
  const [duel, setDuel] = useState<any>(null);
  const [hostDuelToken, setHostDuelToken] = useState("");
  const [guestDuelToken, setGuestDuelToken] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [activePlayer, setActivePlayer] = useState<
    "host" | "guest" | "spectator"
  >("host");
  const [spectatorDuelToken, setSpectatorDuelToken] = useState("");
  const [cards, setCards] = useState<CardSummary[]>([]);

  const createMatch = async () => {
    setLoading(true);
    setError(null);

    try {
      const cardsRes = await fetch("/api/cards");

      if (!cardsRes.ok) {
        throw new Error(`HTTP error! status: ${cardsRes.status}`);
      }

      const availableCards = (await cardsRes.json()) as CardSummary[];
      setCards(availableCards);

      const payload = {
        hostId: "1",
        hostDeck: buildDeck(availableCards, 0),
        guestId: "2",
        guestDeck: buildDeck(availableCards, UNIQUE_CARDS_PER_DECK),
        name: "Test Match",
        visibility: "public",
        format: "regular",
      };

      const matchRes = await fetch("/api/match", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!matchRes.ok) {
        throw new Error(`HTTP error! status: ${matchRes.status}`);
      }

      const match = await matchRes.json();

      const hostToken = await new SignJWT({ id: "1", username: "Player1" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);
      const guestToken = await new SignJWT({ id: "2", username: "Player2" })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);
      const spectatorToken = await new SignJWT({
        id: "3",
        username: "Spectator",
      })
        .setProtectedHeader({ alg: "HS256" })
        .sign(DUEL_TOKEN_SECRET);

      setHostDuelToken(hostToken);
      setGuestDuelToken(guestToken);
      setSpectatorDuelToken(spectatorToken);
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
      } else if (event.key === "3") {
        setActivePlayer("spectator");
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
              hostUrl="ws://localhost:3001"
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
              hostUrl="ws://localhost:3001"
              duelId={duel.id}
              duelToken={guestDuelToken}
              devTools={{
                cards,
                activePlayer,
                onPlayerSwitch: setActivePlayer,
              }}
            />
          </div>

          <div className={activePlayer === "spectator" ? "block" : "hidden"}>
            <Duel
              hostUrl="ws://localhost:3001"
              duelId={duel.id}
              duelToken={spectatorDuelToken}
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
