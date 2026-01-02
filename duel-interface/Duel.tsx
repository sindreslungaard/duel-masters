import { useEffect, useMemo, useState } from "react";
import { useDuel } from "./useDuel";
import {
  cardHasFlag,
  CardState,
  MatchState,
  PLAYABLE_FLAG,
  ShieldState,
  TAP_ABILITY_FLAG,
  TAPPED_FLAG,
} from "./types";
import { Card } from "./Card";
import { Button } from "./Button";

interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
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

export function Duel({ duelId, duelToken, hostUrl }: DuelProps) {
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
    state,
  } = useDuel({
    hostUrl,
    duelId,
    duelToken,
  });

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
  }, [selectedCardId, state]);

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
    if (dragStartPosition && !dragState) {
      setSelectedCardId(dragStartPosition.virtualId);
      console.log("selected card", dragStartPosition.virtualId);
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
      // Handle attack on playzone
      console.log("Attack playzone", dragState.virtualId);
    } else if (
      dragState.sourceZone === "myPlayzone" &&
      dropZone === "opponentShieldzone"
    ) {
      // Handle attack on shields
      console.log("Attack shields", dragState.virtualId);
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
      const zones: DragZone[] = ["myPlayzone"];
      if (!state?.hasAddedManaThisRound) {
        zones.push("myManazone");
      }
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

  return (
    <div
      className="w-full h-screen text-white flex bg-[url('https://i.imgur.com/mWy5Cnl.gif')] bg-cover bg-center gap-2 p-2"
      style={dragState ? { cursor: "grabbing" } : {}}
    >
      <div className="w-[300px] flex flex-col gap-2">
        <div className="flex-1 bg-black/50 rounded-md"></div>

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
                      onClick={() => sendAddToManazone(selectedCard.virtualId)}
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
                      disabled={!selectedCard.canPlay}
                    >
                      Attack player
                    </Button>
                  </div>
                  <div className="flex-1 min-w-0">
                    <Button
                      onClick={() => sendAttackCreature(selectedCard.virtualId)}
                      disabled={!selectedCard.canPlay}
                    >
                      Attack creature
                    </Button>
                  </div>
                  <div className="flex-1 min-w-0">
                    <Button
                      onClick={() => sendTapAbility(selectedCard.virtualId)}
                      disabled={!selectedCard.canPlay}
                    >
                      Tap ability
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        <div className="bg-black/50 p-2 rounded-md">
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
          className="h-[10%] flex gap-5 pb-1 relative"
          data-dropzone="opponentManazone"
        >
          <div
            className="absolute inset-0 z-0"
            data-dropzone="opponentManazone"
          />
          {isValidDropZone("opponentManazone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.opponent.manazone.map(
              CreateCard({ flipped: true, dragState, zone: "opponentManazone" })
            )}
          </div>
        </div>
        <div
          className="h-[10%] flex gap-5 p-1 w-full relative"
          data-dropzone="opponentShieldzone"
        >
          <div
            className="absolute inset-0 z-0"
            data-dropzone="opponentShieldzone"
          />
          {isValidDropZone("opponentShieldzone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.opponent.shieldzone.map(
              CreateCard({ dragState, zone: "opponentShieldzone" })
            )}
          </div>
        </div>
        <div
          className="flex h-[20%] gap-5 p-1 w-full relative"
          data-dropzone="opponentPlayzone"
        >
          <div
            className="absolute inset-0 z-0"
            data-dropzone="opponentPlayzone"
          />
          {isValidDropZone("opponentPlayzone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.opponent.playzone.map(
              CreateCard({ flipped: true, dragState, zone: "opponentPlayzone" })
            )}
          </div>
        </div>
        <div
          className="flex h-[20%] gap-5 p-1 w-full relative"
          data-dropzone="myPlayzone"
        >
          <div className="absolute inset-0 z-0" data-dropzone="myPlayzone" />
          {isValidDropZone("myPlayzone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.me.playzone.map(
              CreateCard({
                dragState,
                zone: "myPlayzone",
                draggable: true,
                onDragStart: handleCardDragStart,
              })
            )}
          </div>
        </div>
        <div
          className="flex h-[10%] gap-5 p-1 w-full relative"
          data-dropzone="myShieldzone"
        >
          <div className="absolute inset-0 z-0" data-dropzone="myShieldzone" />
          {isValidDropZone("myShieldzone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.me.shieldzone.map(
              CreateCard({ dragState, zone: "myShieldzone" })
            )}
          </div>
        </div>
        <div
          className="flex h-[10%] gap-5 p-1 w-full relative"
          data-dropzone="myManazone"
        >
          <div className="absolute inset-0 z-0" data-dropzone="myManazone" />
          {isValidDropZone("myManazone") && (
            <div className="absolute inset-0 bg-green-500/30 border-2 border-green-500 rounded-md pointer-events-none z-10" />
          )}
          <div className="relative z-10 flex gap-5 w-full">
            {state.me.manazone.map(
              CreateCard({ flipped: true, dragState, zone: "myManazone" })
            )}
          </div>
        </div>
        <div
          className="flex h-[20%] gap-5 pt-1 w-full relative"
          data-dropzone="hand"
        >
          <div className="absolute inset-0 z-0" data-dropzone="hand" />
          <div className="relative z-10 flex gap-5 w-full">
            {state.me.hand.map(
              CreateCard({
                interactable: true,
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
                draggable: true,
                onDragStart: handleCardDragStart,
              })
            )}
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
            className={`h-[150px] rounded-md opacity-90 ${
              dragState.rotated ? "rotate-90" : ""
            }`}
            style={{ pointerEvents: "none" }}
          />
        </div>
      )}
    </div>
  );
}

function CreateCard(
  options: {
    interactable?: boolean;
    canAddToManazone?: boolean;
    flipped?: boolean;
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
      ></Card>
    );
  };
}
