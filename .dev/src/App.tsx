import { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import "./App.css";
import { SignJWT } from "jose";
import {
  Duel,
  type DuelChatUser,
  type DuelChatUserResolver,
  type DuelChatUserTriggerProps,
} from "../../packages/duel-interface";

const DUEL_TOKEN_SECRET = new TextEncoder().encode("duel-secret");
const DECK_SIZE = 40;
const UNIQUE_CARDS_PER_DECK = 10;

type CardSummary = {
  uid: string;
  name: string;
};

type DevDuel = {
  id: string;
};

type DevChatProfile = DuelChatUser & {
  joined: string;
  gamesPlayed: number;
  badges: { name: string; color: string }[];
};

function createAvatarDataUrl(
  initials: string,
  startColor: string,
  endColor: string,
): string {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 96 96">
      <defs>
        <linearGradient id="background" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stop-color="${startColor}" />
          <stop offset="1" stop-color="${endColor}" />
        </linearGradient>
      </defs>
      <rect width="96" height="96" rx="48" fill="url(#background)" />
      <circle cx="72" cy="23" r="20" fill="white" fill-opacity="0.12" />
      <text x="48" y="58" text-anchor="middle" fill="white"
        font-family="Arial, sans-serif" font-size="32" font-weight="700">
        ${initials}
      </text>
    </svg>
  `;

  return `data:image/svg+xml,${encodeURIComponent(svg)}`;
}

const DEV_CHAT_PROFILES: readonly DevChatProfile[] = [
  {
    username: "Player1",
    avatarUrl: createAvatarDataUrl("P1", "#f97316", "#dc2626"),
    joined: "Jan 2024",
    gamesPlayed: 128,
    badges: [{ name: "Fire Duelist", color: "#fb923c" }],
  },
  /*   {
    username: "Player2",
    avatarUrl: createAvatarDataUrl("P2", "#3b82f6", "#7c3aed"),
    joined: "Mar 2024",
    gamesPlayed: 96,
    badges: [{ name: "Water Adept", color: "#60a5fa" }],
  }, */
  {
    username: "Spectator",
    avatarUrl: createAvatarDataUrl("SP", "#10b981", "#0f766e"),
    joined: "Jun 2025",
    gamesPlayed: 42,
    badges: [{ name: "Community Member", color: "#34d399" }],
  },
];

const resolveDevChatUser: DuelChatUserResolver = async (username) => {
  await new Promise((resolve) => window.setTimeout(resolve, 300));

  return (
    DEV_CHAT_PROFILES.find(
      (profile) => profile.username.toLowerCase() === username.toLowerCase(),
    ) ?? null
  );
};

function DevProfilePreviewTrigger({
  user,
  children,
  className,
}: DuelChatUserTriggerProps) {
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const triggerRef = useRef<HTMLButtonElement>(null);
  const previewRef = useRef<HTMLDivElement>(null);
  const profile =
    DEV_CHAT_PROFILES.find(
      (candidate) =>
        candidate.username.toLowerCase() === user.username.toLowerCase(),
    ) ?? null;

  useEffect(() => {
    if (!open) return;

    const closeOnOutsideClick = (event: PointerEvent) => {
      const target = event.target as Node;

      if (
        !triggerRef.current?.contains(target) &&
        !previewRef.current?.contains(target)
      ) {
        setOpen(false);
      }
    };
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };
    const closeOnViewportChange = () => setOpen(false);

    document.addEventListener("pointerdown", closeOnOutsideClick);
    document.addEventListener("keydown", closeOnEscape);
    window.addEventListener("resize", closeOnViewportChange);
    window.addEventListener("scroll", closeOnViewportChange, true);

    return () => {
      document.removeEventListener("pointerdown", closeOnOutsideClick);
      document.removeEventListener("keydown", closeOnEscape);
      window.removeEventListener("resize", closeOnViewportChange);
      window.removeEventListener("scroll", closeOnViewportChange, true);
    };
  }, [open]);

  if (!profile || user.username.startsWith("Guest ")) {
    return <span className={className}>{children}</span>;
  }

  const togglePreview = () => {
    if (!open && triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      const previewWidth = 288;
      const previewHeight = 210;
      const left = Math.max(
        8,
        Math.min(rect.left, window.innerWidth - previewWidth - 8),
      );
      const top =
        rect.bottom + previewHeight + 8 <= window.innerHeight
          ? rect.bottom + 8
          : Math.max(8, rect.top - previewHeight - 8);

      setPosition({ top, left });
    }

    setOpen((current) => !current);
  };

  return (
    <>
      <button
        ref={triggerRef}
        type="button"
        aria-expanded={open}
        aria-haspopup="dialog"
        className={`block border-0 bg-transparent p-0 text-left leading-none text-inherit cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 ${className}`}
        onClick={togglePreview}
      >
        {children}
      </button>
      {open &&
        createPortal(
          <div
            ref={previewRef}
            role="dialog"
            aria-label={`${user.username} profile preview`}
            className="fixed z-[100] w-72 space-y-4 rounded-xl border border-white/15 bg-zinc-900 p-4 text-white shadow-2xl"
            style={position}
          >
            <div className="flex items-center gap-3">
              {profile?.avatarUrl ? (
                <img
                  src={profile.avatarUrl}
                  alt=""
                  className="size-12 rounded-full"
                />
              ) : (
                <div className="flex size-12 items-center justify-center rounded-full bg-white/10 font-semibold">
                  {user.username.slice(0, 2).toUpperCase()}
                </div>
              )}
              <div className="min-w-0">
                <div className="truncate font-semibold">{user.username}</div>
                <div className="text-xs text-zinc-400">
                  Joined {profile?.joined ?? "locally"}
                </div>
              </div>
            </div>
            <div className="rounded-lg bg-white/5 p-2 text-center">
              <div className="font-medium">{profile?.gamesPlayed ?? 0}</div>
              <div className="text-xs text-zinc-400">Games</div>
            </div>
            {profile && profile.badges.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {profile.badges.map((badge) => (
                  <span
                    key={badge.name}
                    className="rounded-full border px-2.5 py-1 text-xs font-medium"
                    style={{
                      borderColor: `${badge.color}55`,
                      backgroundColor: `${badge.color}14`,
                      color: badge.color,
                    }}
                  >
                    {badge.name}
                  </span>
                ))}
              </div>
            )}
            <div className="text-xs text-zinc-500">
              Local profile preview for duel-interface development.
            </div>
          </div>,
          document.body,
        )}
    </>
  );
}

const renderDevProfilePreview = (props: DuelChatUserTriggerProps) => (
  <DevProfilePreviewTrigger {...props} />
);

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
  const [duel, setDuel] = useState<DevDuel | null>(null);
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

      const match = (await matchRes.json()) as DevDuel;

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
              resolveChatUser={resolveDevChatUser}
              renderChatUserTrigger={renderDevProfilePreview}
              devTools={{
                cards,
                activePlayer,
                onPlayerSwitch: setActivePlayer,
              }}
              onDuelFinished={(message) => {
                alert(
                  `Duel finished! Winner: ${message.winner?.username || "Unknown"}, Turns: ${message.turns}, Duration: ${message.durationSeconds} seconds`,
                );
              }}
            />
          </div>

          <div className={activePlayer === "guest" ? "block" : "hidden"}>
            <Duel
              hostUrl="ws://localhost:3001"
              duelId={duel.id}
              duelToken={guestDuelToken}
              resolveChatUser={resolveDevChatUser}
              renderChatUserTrigger={renderDevProfilePreview}
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
              resolveChatUser={resolveDevChatUser}
              renderChatUserTrigger={renderDevProfilePreview}
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
