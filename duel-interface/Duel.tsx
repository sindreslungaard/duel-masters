import { useEffect, useMemo, useState } from "react";
import { useDuel } from "./useDuel";
import {
  ActionMessage,
  ActionWarningMessage,
  cardHasFlag,
  CardState,
  ChatMessage,
  MatchState,
  PLAYABLE_FLAG,
  ShieldState,
  TAP_ABILITY_FLAG,
  TAPPED_FLAG,
} from "./types";
import { Card } from "./Card";
import { Button } from "./Button";
import { Popup } from "./Popup";
import { Action } from "./Action";
import { Chat } from "./Chat";
import { CardPreview } from "./CardPreview";
import { MultiCardPreview } from "./MultiCardPreview";

const scrollbarStyles = `
  .custom-scrollbar::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }
  
  .custom-scrollbar::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.2);
    border-radius: 4px;
  }
  
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
    border-radius: 4px;
  }
  
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.5);
  }
  
  .custom-scrollbar {
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.3) rgba(0, 0, 0, 0.2);
  }
`;

interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
  devTools?: {
    cards: { uid: string; name: string }[];
    activePlayer: "host" | "guest";
    onPlayerSwitch: (player: "host" | "guest") => void;
  };
}

type DragZone =
  | "hand"
  | "myPlayzone"
  | "opponentPlayzone"
  | "myManazone"
  | "opponentManazone"
  | "myShieldzone"
  | "opponentShieldzone";

interface DragState {
  virtualId: string;
  imageId: string;
  name?: string;
  sourceZone: DragZone;
  mouseX: number;
  mouseY: number;
  rotated?: boolean;
}

interface SelectedCard {
  virtualId: string;
  name: string;
  canPlay: boolean;
  hasTapAbility: boolean;
  zone: "hand" | "battlezone";
}

interface PreviewCard {
  name: string;
  imageId: string;
}

interface PreviewCards {
  text?: string;
  cards: PreviewCard[];
}

interface Action {}

function DevToolSection({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <fieldset className="border border-gray-700 rounded px-3 pb-3 pt-2">
      <legend className="px-2 text-xs text-gray-500">{title}</legend>
      {children}
    </fieldset>
  );
}

export function Duel({ duelId, duelToken, hostUrl, devTools }: DuelProps) {
  const [action, setAction] = useState<ActionMessage | null>(null);
  const [actionError, setActionError] = useState<ActionWarningMessage | null>(
    null
  );
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [wait, setWait] = useState("");
  const [warningMessage, setWarningMessage] = useState("");
  const [dots, setDots] = useState<"." | ".." | "...">(".");

  const {
    connected,
    error,
    send,
    sendJoinMatch,
    sendEndTurn,
    sendAddToBattlezone,
    sendAddToManazone,
    sendAttackPlayer,
    sendAttackCreature,
    sendTapAbility,
    sendAction,
    sendChat,
    state,
  } = useDuel({
    hostUrl,
    duelId,
    duelToken,
    onActionMessage: setAction,
    onActionError: setActionError,
    onActionClose: () => {
      setAction(null);
      setActionError(null);
    },
    onChat: (data) => {
      setChatMessages((prev) => [...prev, data]);
    },
    onWarning: (data) => {
      setWarningMessage(data.message);
    },
    onWait: (data) => {
      setWait(data.message);
    },
    onEndWait: () => {
      setWait("");
    },
  });

  // Popup modal related refs
  const [opponentDisconnected, setOpponentDisconnected] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const [reconnecting, setReconnecting] = useState(false);

  const [previewCard, setPreviewCard] = useState<PreviewCard | null>(null);
  const [previewCards, setPreviewCards] = useState<PreviewCards | null>(null);
  const [multiCardView, setMultiCardView] = useState<{
    cards: { imageId: string; name: string }[];
    title: string;
  } | null>(null);

  const [showPopup1, setShowPopup1] = useState(true);

  const [selectedCardId, setSelectedCardId] = useState<string | null>(null);
  const [selectedCard, setSelectedCard] = useState<SelectedCard | null>(null);

  useEffect(() => {
    if (selectedCardId) {
      let zone: "hand" | "battlezone" = "hand";
      let card = state?.me.hand.find((c) => c.virtualId === selectedCardId);

      if (!card) {
        card = state?.me.playzone.find((c) => c.virtualId === selectedCardId);
        zone = "battlezone";
      }

      if (!card) {
        setSelectedCard(null);
        return;
      }

      const canPlay = cardHasFlag(card.flags, PLAYABLE_FLAG);
      const hasTapAbility = cardHasFlag(card.flags, TAP_ABILITY_FLAG);

      setSelectedCard({
        virtualId: card.virtualId,
        name: card.name || "",
        canPlay,
        hasTapAbility,
        zone,
      });
    } else {
      setSelectedCard(null);
    }
  }, [selectedCardId]);

  useEffect(() => {
    setSelectedCardId(null);
    setSelectedCard(null);
  }, [state]);

  useEffect(() => {
    const interval = setInterval(() => {
      setDots((prev) => {
        if (prev === ".") return "..";
        if (prev === "..") return "...";
        return ".";
      });
    }, 500);

    return () => clearInterval(interval);
  }, []);

  const [dragState, setDragState] = useState<DragState | null>(null);
  const [dropZone, setDropZone] = useState<DragZone | null>(null);
  const [dragStartPosition, setDragStartPosition] = useState<{
    x: number;
    y: number;
    virtualId: string;
    imageId: string;
    name?: string;
    sourceZone: DragZone;
    rotated: boolean;
  } | null>(null);

  const DRAG_THRESHOLD = 5; // pixels

  const handleCardDragStart = (
    virtualId: string,
    imageId: string,
    name: string | undefined,
    sourceZone: DragZone,
    rotated: boolean,
    e: React.MouseEvent | React.TouchEvent
  ) => {
    // Ignore right clicks
    if ("button" in e && e.button === 2) {
      return;
    }

    const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
    const clientY = "touches" in e ? e.touches[0].clientY : e.clientY;

    // Record start position but don't start dragging yet
    setDragStartPosition({
      x: clientX,
      y: clientY,
      virtualId,
      imageId,
      name,
      sourceZone,
      rotated,
    });
  };

  const handleMouseMove = (e: MouseEvent | TouchEvent) => {
    const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
    const clientY = "touches" in e ? e.touches[0].clientY : e.clientY;

    // Check if we should start dragging based on threshold
    if (dragStartPosition && !dragState) {
      const deltaX = Math.abs(clientX - dragStartPosition.x);
      const deltaY = Math.abs(clientY - dragStartPosition.y);

      if (deltaX > DRAG_THRESHOLD || deltaY > DRAG_THRESHOLD) {
        // Start dragging
        setDragState({
          virtualId: dragStartPosition.virtualId,
          imageId: dragStartPosition.imageId,
          name: dragStartPosition.name,
          sourceZone: dragStartPosition.sourceZone,
          mouseX: clientX,
          mouseY: clientY,
          rotated: dragStartPosition.rotated,
        });
        // Select the card being dragged
        setSelectedCardId(dragStartPosition.virtualId);
        setDragStartPosition(null);
      }
      return;
    }

    if (!dragState) return;

    setDragState({
      ...dragState,
      mouseX: clientX,
      mouseY: clientY,
    });

    // Check all drop zones and find which one contains the cursor
    let foundZone: DragZone | null = null;
    const dropZones = document.querySelectorAll("[data-dropzone]");
    dropZones.forEach((zone) => {
      const rect = zone.getBoundingClientRect();
      if (
        clientX >= rect.left &&
        clientX <= rect.right &&
        clientY >= rect.top &&
        clientY <= rect.bottom
      ) {
        foundZone = zone.getAttribute("data-dropzone") as DragZone;
      }
    });

    setDropZone(foundZone);
  };

  const handleMouseUp = () => {
    // If we have a dragStartPosition but no dragState, it's a click
    if (dragStartPosition && !dragState && state?.myTurn) {
      // Toggle selection: unselect if already selected, select if not
      setSelectedCardId((prev) =>
        prev === dragStartPosition.virtualId
          ? null
          : dragStartPosition.virtualId
      );
      setDragStartPosition(null);
      return;
    }

    if (!dragState || !dropZone) {
      setDragState(null);
      setDropZone(null);
      setDragStartPosition(null);
      return;
    }

    // Handle drop actions based on source and target zones
    if (dragState.sourceZone === "hand" && dropZone === "myManazone") {
      sendAddToManazone(dragState.virtualId);
    } else if (dragState.sourceZone === "hand" && dropZone === "myPlayzone") {
      sendAddToBattlezone(dragState.virtualId);
    } else if (
      dragState.sourceZone === "myPlayzone" &&
      dropZone === "opponentPlayzone"
    ) {
      sendAttackCreature(dragState.virtualId);
    } else if (
      dragState.sourceZone === "myPlayzone" &&
      dropZone === "opponentShieldzone"
    ) {
      sendAttackPlayer(dragState.virtualId);
    }

    setDragState(null);
    setDropZone(null);
    setDragStartPosition(null);
  };

  useEffect(() => {
    if (dragState || dragStartPosition) {
      window.addEventListener("mousemove", handleMouseMove);
      window.addEventListener("touchmove", handleMouseMove);
      window.addEventListener("mouseup", handleMouseUp);
      window.addEventListener("touchend", handleMouseUp);

      return () => {
        window.removeEventListener("mousemove", handleMouseMove);
        window.removeEventListener("touchmove", handleMouseMove);
        window.removeEventListener("mouseup", handleMouseUp);
        window.removeEventListener("touchend", handleMouseUp);
      };
    }
  }, [dragState, dropZone, dragStartPosition]);

  useEffect(() => {
    if (connected) {
      sendJoinMatch();
    }
  }, [connected]);

  if (!state) {
    return <div>Waiting for both players to join...</div>;
  }

  const getValidDropZones = (sourceZone: DragZone): DragZone[] => {
    if (sourceZone === "hand") {
      const zones: DragZone[] = ["myPlayzone", "myManazone"];
      return zones;
    }
    if (sourceZone === "myPlayzone") {
      return ["opponentPlayzone", "opponentShieldzone"];
    }
    return [];
  };

  const isValidDropZone = (zone: DragZone): boolean => {
    if (!dragState) return false;
    return getValidDropZones(dragState.sourceZone).includes(zone);
  };

  const getDropZoneColor = (zone: DragZone): "green" | "red" | null => {
    if (!dragState || !isValidDropZone(zone)) return null;

    // Check if the action is actually allowed
    if (dragState.sourceZone === "hand") {
      // Find the card being dragged
      const draggedCard = state?.me.hand.find(
        (c) => c.virtualId === dragState.virtualId
      );

      if (zone === "myPlayzone") {
        // Check if card can be played
        return draggedCard && cardHasFlag(draggedCard.flags, PLAYABLE_FLAG)
          ? "green"
          : "red";
      }

      if (zone === "myManazone") {
        // Check if mana can be added
        return state?.hasAddedManaThisRound ? "red" : "green";
      }
    }

    // Check if attacking opponent's battlezone when it's empty
    if (dragState.sourceZone === "myPlayzone" && zone === "opponentPlayzone") {
      // If opponent has no creatures, highlight red
      return state?.opponent.playzone.length === 0 ? "red" : "green";
    }

    // For other zones (attacking shields), always green if valid
    return "green";
  };

  return (
    <>
      <style>{scrollbarStyles}</style>
      <div
        className="w-full h-screen text-white flex bg-[linear-gradient(45deg,rgb(29,33,42),rgb(20,16,21))] bg-cover bg-no-repeat gap-2 p-1 custom-scrollbar"
        style={dragState ? { cursor: "grabbing" } : {}}
      >
        {/* <Popup
        title="Your Title Here"
        visible={showPopup1}
        onClose={() => setShowPopup1(false)}
        maxWidth="500px"
        zIndex={1000}
      >
        <div className="p-6">Your content here</div>
      </Popup> */}

        <div className="w-[300px] flex flex-col gap-2">
          {/* Devtools */}
          {devTools && (
            <div className="bg-black/30 rounded-md overflow-hidden p-3 text-sm">
              <p className="mb-3 font-semibold">Development Tools</p>

              <DevToolSection title="Player Switch">
                <div className="flex gap-2">
                  <div className="flex-1">
                    <Button
                      variant={
                        devTools.activePlayer === "host" ? "default" : "gray"
                      }
                      onClick={() => devTools.onPlayerSwitch("host")}
                    >
                      Player 1
                    </Button>
                  </div>
                  <div className="flex-1">
                    <Button
                      variant={
                        devTools.activePlayer === "guest" ? "default" : "gray"
                      }
                      onClick={() => devTools.onPlayerSwitch("guest")}
                    >
                      Player 2
                    </Button>
                  </div>
                </div>
              </DevToolSection>

              <div className="mt-3">
                <DevToolSection title="Initialize">
                  <Button variant="gray" onClick={() => sendChat("/init all")}>
                    Initialize zones with 1 of each race
                  </Button>
                </DevToolSection>
              </div>

              <div className="mt-3">
                <DevToolSection title="Add Cards">
                  <div className="flex gap-2">
                    <div className="flex-1">
                      <select
                        className="w-full bg-gray-800 text-white px-2 py-[0.4rem] rounded border border-gray-700 focus:outline-none focus:border-blue-500 text-xs"
                        id="card-selector"
                        defaultValue=""
                      >
                        <option value="" disabled>
                          Select a card...
                        </option>
                        {devTools?.cards.map((card) => (
                          <option key={card.uid} value={card.uid}>
                            {card.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="">
                      <Button
                        variant="gray"
                        onClick={() => {
                          const select = document.getElementById(
                            "card-selector"
                          ) as HTMLSelectElement;
                          if (select.value) {
                            sendChat(`/add ${select.value}`);
                          }
                        }}
                      >
                        Add
                      </Button>
                    </div>
                  </div>
                </DevToolSection>
              </div>
            </div>
          )}

          {/* Chat */}
          <div className="flex-1 bg-black/30 rounded-md overflow-hidden">
            <Chat messages={chatMessages} onSendMessage={sendChat} />
          </div>

          {/* Actions */}
          <div className="bg-black/50 p-2 rounded-md h-[72px] text-gray-400">
            {selectedCard && state.myTurn && (
              <div className="flex flex-col gap-2">
                <div className="flex-1 text-xs whitespace-nowrap overflow-hidden text-ellipsis">
                  {selectedCard.name}
                </div>
                {selectedCard.zone === "hand" && (
                  <div className="flex gap-2">
                    {/* Hand zone */}
                    <div className="flex-1 min-w-0">
                      <Button
                        onClick={() =>
                          sendAddToBattlezone(selectedCard.virtualId)
                        }
                        disabled={!selectedCard.canPlay}
                      >
                        Summon
                      </Button>
                    </div>
                    <div className="flex-1 min-w-0">
                      <Button
                        onClick={() =>
                          sendAddToManazone(selectedCard.virtualId)
                        }
                        disabled={state.hasAddedManaThisRound}
                      >
                        Add to manazone
                      </Button>
                    </div>
                  </div>
                )}

                {selectedCard.zone === "battlezone" && (
                  <div className="flex gap-2">
                    <div className="flex-1 min-w-0">
                      <Button
                        onClick={() => sendAttackPlayer(selectedCard.virtualId)}
                      >
                        Attack Player
                      </Button>
                    </div>
                    <div className="flex-1 min-w-0">
                      <Button
                        onClick={() =>
                          sendAttackCreature(selectedCard.virtualId)
                        }
                      >
                        Attack Creature
                      </Button>
                    </div>
                    {selectedCard.hasTapAbility && (
                      <div className="flex-1 min-w-0">
                        <Button
                          onClick={() => sendTapAbility(selectedCard.virtualId)}
                        >
                          Tap Ability
                        </Button>
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}
          </div>

          {/* End turn / forfeit */}
          <div className="bg-black/30 p-2 rounded-md">
            <Button
              onClick={sendEndTurn}
              disabled={!state.myTurn}
              disabledTooltip="It's not your turn"
            >
              End turn
            </Button>
          </div>
        </div>
        <div className="flex flex-1 flex-col h-full w-full">
          <div
            className="h-[10%] flex gap-5 pb-1 relative overflow-x-auto"
            data-dropzone="opponentManazone"
          >
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentManazone"
            />
            {getDropZoneColor("opponentManazone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("opponentManazone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full">
              {state.opponent.manazone.map(
                CreateCard({
                  flipped: true,
                  dragState,
                  zone: "opponentManazone",
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="h-[10%] flex gap-5 p-1 w-full relative overflow-x-auto"
            data-dropzone="opponentShieldzone"
          >
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentShieldzone"
            />
            {getDropZoneColor("opponentShieldzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("opponentShieldzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full">
              {state.opponent.shieldzone.map(
                CreateCard({
                  dragState,
                  zone: "opponentShieldzone",
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="flex h-[20%] gap-5 p-1 w-full relative overflow-x-auto"
            data-dropzone="opponentPlayzone"
          >
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentPlayzone"
            />
            {getDropZoneColor("opponentPlayzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("opponentPlayzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full">
              {state.opponent.playzone.map(
                CreateCard({
                  flipped: true,
                  dragState,
                  zone: "opponentPlayzone",
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="flex h-[20%] gap-5 p-1 w-full relative overflow-x-auto"
            data-dropzone="myPlayzone"
          >
            <div className="absolute inset-0 z-0" data-dropzone="myPlayzone" />
            {getDropZoneColor("myPlayzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("myPlayzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full">
              {state.me.playzone.map(
                CreateCard({
                  selected: (id: string) => id === selectedCardId,
                  interactable: state?.myTurn,
                  dragState,
                  zone: "myPlayzone",
                  draggable: state.myTurn,
                  onDragStart: handleCardDragStart,
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="flex h-[10%] gap-5 p-1 w-full relative overflow-x-auto"
            data-dropzone="myShieldzone"
          >
            <div
              className="absolute inset-0 z-0 "
              data-dropzone="myShieldzone"
            />
            {getDropZoneColor("myShieldzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("myShieldzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full">
              {state.me.shieldzone.map(
                CreateCard({
                  dragState,
                  zone: "myShieldzone",
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="flex h-[10%] gap-5 p-1 w-full relative overflow-x-auto"
            data-dropzone="myManazone"
          >
            <div className="absolute inset-0 z-0" data-dropzone="myManazone" />
            {getDropZoneColor("myManazone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-10" />
            )}
            {getDropZoneColor("myManazone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-10" />
            )}
            <div className="relative z-10 flex gap-5 w-full ">
              {state.me.manazone.map(
                CreateCard({
                  flipped: true,
                  dragState,
                  zone: "myManazone",
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
          <div
            className="flex h-[20%] gap-5 pt-1 w-full relative overflow-x-auto"
            data-dropzone="hand"
          >
            <div className="absolute inset-0 z-0" data-dropzone="hand" />
            <div className="relative z-10 flex gap-5 w-full p-px">
              {state.me.hand.map(
                CreateCard({
                  selected: (id: string) => id === selectedCardId,
                  interactable: state?.myTurn,
                  canAddToManazone: !state.hasAddedManaThisRound,
                  onAddToBattlezone: (virtualId) => {
                    sendAddToBattlezone(virtualId);
                  },
                  onAddToManazone: (virtualId) => {
                    sendAddToManazone(virtualId);
                  },
                  onTapAbility: (virtualId) => {
                    sendTapAbility(virtualId);
                  },
                  dragState,
                  zone: "hand",
                  draggable: state.myTurn,
                  onDragStart: handleCardDragStart,
                  onRightClick: (imageId, name) =>
                    setPreviewCard({ imageId, name: name || "" }),
                })
              )}
            </div>
          </div>
        </div>

        {/* Player Info Panel - Right Side */}
        <div className="fixed right-[0.5vw] top-1/2 -translate-y-2/3 w-[12vw] min-w-[100px] max-w-[160px] flex flex-col gap-[0.5vh] md:gap-[5vh] z-20">
          {/* Opponent Section */}
          <div className="rounded-lg flex flex-col gap-[1vh]">
            {/* Opponent Hand Count */}
            <div className="text-center">
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh]">
                Hand [{state.opponent.handCount}]
              </p>
            </div>

            {/* Opponent Deck */}
            <div>
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh] text-center">
                Deck [{state.opponent.deck}]
              </p>
              <div className="relative h-[12vh] min-h-[60px] max-h-[110px] flex items-center justify-center">
                <img
                  src="https://scans.shobu.io/backside.jpg"
                  alt="Deck back"
                  className="h-full"
                  style={{ borderRadius: "5%" }}
                />
              </div>
            </div>

            {/* Opponent Graveyard */}
            <div>
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh] text-center">
                Graveyard [{state.opponent.graveyard.length}]
              </p>
              <div className="relative h-[12vh] min-h-[60px] max-h-[110px] flex items-center justify-center">
                {state.opponent.graveyard.length > 0 ? (
                  <img
                    src={`https://scans.shobu.io/${
                      state.opponent.graveyard[
                        state.opponent.graveyard.length - 1
                      ].uid
                    }.jpg`}
                    alt="Top graveyard card"
                    className="h-full cursor-pointer hover:scale-105 transition-transform"
                    style={{ borderRadius: "5%" }}
                    onClick={() => {
                      setMultiCardView({
                        cards: state.opponent.graveyard.map((card) => ({
                          imageId: card.uid,
                          name: card.name,
                        })),
                        title: "Opponent's Graveyard",
                      });
                    }}
                    onContextMenu={(e) => {
                      e.preventDefault();
                      setMultiCardView({
                        cards: state.opponent.graveyard.map((card) => ({
                          imageId: card.uid,
                          name: card.name,
                        })),
                        title: "Opponent's Graveyard",
                      });
                    }}
                  />
                ) : (
                  <img
                    src="https://scans.shobu.io/backside.jpg"
                    alt="Empty graveyard"
                    className="h-full opacity-30"
                    style={{ borderRadius: "5%" }}
                  />
                )}
              </div>
            </div>
          </div>

          {/* Player Section */}
          <div className="rounded-lg flex flex-col gap-[1vh]">
            {/* Player Graveyard */}
            <div>
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh] text-center">
                Graveyard [{state.me.graveyard.length}]
              </p>
              <div className="relative h-[12vh] min-h-[60px] max-h-[110px] flex items-center justify-center">
                {state.me.graveyard.length > 0 ? (
                  <img
                    src={`https://scans.shobu.io/${
                      state.me.graveyard[state.me.graveyard.length - 1].uid
                    }.jpg`}
                    alt="Top graveyard card"
                    className="h-full cursor-pointer hover:scale-105 transition-transform"
                    style={{ borderRadius: "5%" }}
                    onClick={() => {
                      setMultiCardView({
                        cards: state.me.graveyard.map((card) => ({
                          imageId: card.uid,
                          name: card.name,
                        })),
                        title: "My Graveyard",
                      });
                    }}
                    onContextMenu={(e) => {
                      e.preventDefault();
                      setMultiCardView({
                        cards: state.me.graveyard.map((card) => ({
                          imageId: card.uid,
                          name: card.name,
                        })),
                        title: "My Graveyard",
                      });
                    }}
                  />
                ) : (
                  <img
                    src="https://scans.shobu.io/backside.jpg"
                    alt="Empty graveyard"
                    className="h-full opacity-30"
                    style={{ borderRadius: "5%" }}
                  />
                )}
              </div>
            </div>

            {/* Player Deck */}
            <div>
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh] text-center">
                Deck [{state.me.deck}]
              </p>
              <div className="relative h-[12vh] min-h-[60px] max-h-[110px] flex items-center justify-center">
                <img
                  src="https://scans.shobu.io/backside.jpg"
                  alt="Deck back"
                  className="h-full"
                  style={{ borderRadius: "5%" }}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Floating card that follows cursor */}
        {dragState && (
          <div
            className="fixed pointer-events-none z-50"
            style={{
              left: dragState.mouseX,
              top: dragState.mouseY,
              transform: "translate(-50%, -50%)",
              pointerEvents: "none",
            }}
          >
            <img
              src={`https://scans.shobu.io/${dragState.imageId}.jpg`}
              alt={dragState.name || "Card"}
              className={`h-[150px] opacity-90 ${
                dragState.rotated ? "rotate-90" : ""
              }`}
              style={{ pointerEvents: "none", borderRadius: "5%" }}
            />
          </div>
        )}
      </div>

      <CardPreview
        visible={!!previewCard}
        imageId={previewCard?.imageId || null}
        name={previewCard?.name || null}
        onClose={() => setPreviewCard(null)}
      />

      <MultiCardPreview
        visible={!!multiCardView}
        cards={multiCardView?.cards || []}
        title={multiCardView?.title || ""}
        onClose={() => setMultiCardView(null)}
        onCardClick={(imageId, name) => {
          setPreviewCard({ imageId, name });
          setMultiCardView(null);
        }}
      />

      <Popup
        visible={!!warningMessage.length}
        onClose={() => setWarningMessage("")}
        title="Warning"
        maxWidth="500px"
        closeOnOutsideClick={false}
      >
        <div className="p-6 text-white ">{warningMessage}</div>
      </Popup>

      <Popup
        visible={!!wait.length}
        onClose={() => setWait("")}
        title="Wait"
        maxWidth="500px"
        closeOnOutsideClick={false}
        showCloseButton={false}
      >
        <div className="p-6 text-white ">
          {wait}
          {dots}
        </div>
      </Popup>

      {action && (
        <Action
          title="Action Required"
          visible={true}
          error={actionError ? actionError.message : undefined}
          actionType={action.actionType}
          cards={action.cards}
          text={action.text}
          minSelections={action.minSelections}
          maxSelections={action.maxSelections}
          cancellable={action.cancellable}
          unselectableCards={action.unselectableCards}
          choices={action.choices}
          onChoose={sendAction}
          onClose={() => sendAction({ cards: [], cancel: true })}
          onCardRightClick={(imageId, name) =>
            setPreviewCard({ imageId, name: name || "" })
          }
        ></Action>
      )}
    </>
  );
}

function CreateCard(
  options: {
    interactable?: boolean;
    canAddToManazone?: boolean;
    flipped?: boolean;
    selected?: (virtualId: string) => boolean;
    onAddToBattlezone?: (virtualId: string) => void;
    onAddToManazone?: (virtualId: string) => void;
    onTapAbility?: (virtualId: string) => void;
    dragState?: DragState | null;
    zone?: DragZone;
    draggable?: boolean;
    onDragStart?: (
      virtualId: string,
      imageId: string,
      name: string | undefined,
      sourceZone: DragZone,
      rotated: boolean,
      e: React.MouseEvent | React.TouchEvent
    ) => void;
    onRightClick?: (imageId: string, name?: string) => void;
  } = {}
) {
  return (card: CardState | ShieldState, index: number) => {
    const name = "name" in card && card.name ? card.name : undefined;
    const rotated = cardHasFlag(card.flags, TAPPED_FLAG);
    const isDragging = options.dragState?.virtualId === card.virtualId;

    return (
      <Card
        virtualId={card.virtualId}
        name={name}
        imageId={card.uid}
        key={index}
        rotated={rotated}
        selected={options.selected ? options.selected(card.virtualId) : false}
        interactable={options.interactable}
        canAddToBattlezone={cardHasFlag(card.flags, PLAYABLE_FLAG)}
        canAddToManazone={options.canAddToManazone}
        onAddToBattlezone={options.onAddToBattlezone}
        onAddToManazone={options.onAddToManazone}
        onTapAbility={options.onTapAbility}
        flipped={options.flipped}
        isDragging={isDragging}
        draggable={options.draggable}
        onDragStart={(e) => {
          if (options.onDragStart && options.zone && card.uid) {
            options.onDragStart(
              card.virtualId,
              card.uid,
              name,
              options.zone,
              rotated,
              e
            );
          }
        }}
        onRightClick={() => {
          if (options.onRightClick && card.uid) {
            options.onRightClick(card.uid, name);
          }
        }}
      ></Card>
    );
  };
}
