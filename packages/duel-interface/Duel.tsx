import { useEffect, useMemo, useRef, useState } from "react";
import { useDuel } from "./useDuel";
import {
  ActionMessage,
  ActionWarningMessage,
  cardHasFlag,
  CardState,
  DuelChatUserResolver,
  DuelChatUserTriggerRenderer,
  DuelFinishedMessage,
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
import { Chat, type ReceivedChatMessage } from "./Chat";
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

export interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
  playmat?: string;
  /**
   * Whether the opponent's cards are flipped from their natural orientation
   * towards the player looking at the screen. Defaults to false, which orients
   * them towards the opponent the way they would face across a physical table.
   */
  flipOpponentCards?: boolean;
  resolveChatUser?: DuelChatUserResolver;
  renderChatUserTrigger?: DuelChatUserTriggerRenderer;
  /**
   * Current usernames whose user-authored chat messages should be hidden.
   * Existing history is retained and re-filtered when this list changes.
   */
  blockedChatUsers?: readonly string[];
  devTools?: {
    cards: { uid: string; name: string }[];
    activePlayer: "host" | "guest" | "spectator";
    onPlayerSwitch: (player: "host" | "guest" | "spectator") => void;
  };
  onNewTurn?: (myTurn: boolean) => void;
  onLeaveDuel?: () => void;
  onDuelFinished?: (message: DuelFinishedMessage) => void;
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

/** How far a pointer may travel before a press counts as a drag rather than a
 * tap. A finger is much less precise than a mouse and always wobbles a little,
 * so touch needs more slack or every tap would be read as a drag. */
const MOUSE_DRAG_THRESHOLD = 5;
const TOUCH_DRAG_THRESHOLD = 12;

/** Below this the fixed 300px side column leaves too little room for the board,
 * so the chat moves into a drawer and the controls into a bar along the bottom.
 * Covers large phones in landscape and tablets in portrait, while laptops from
 * 1280 up keep the side column. Kept in step with the min-[1200px] classes that
 * size the floating overlays. */
const COMPACT_VIEWPORT_QUERY = "(max-width: 1199px)";

function useCompactViewport() {
  const [isCompact, setIsCompact] = useState(() =>
    typeof window === "undefined"
      ? false
      : window.matchMedia(COMPACT_VIEWPORT_QUERY).matches,
  );

  useEffect(() => {
    const query = window.matchMedia(COMPACT_VIEWPORT_QUERY);
    const update = () => setIsCompact(query.matches);

    update();
    query.addEventListener("change", update);

    return () => query.removeEventListener("change", update);
  }, []);

  return isCompact;
}

/** Leaving the duel, shown the same way on every screen size: a door with an
 * arrow heading out, in red because it forfeits the match. */
function ForfeitButton({ onClick }: { onClick: () => void }) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-label="Forfeit"
      title="Forfeit"
      className="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-md bg-gray-800 text-red-500 transition-colors hover:bg-gray-700 hover:text-red-400"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className="h-5 w-5"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fillRule="evenodd"
          d="M3 3a1 1 0 00-1 1v12a1 1 0 102 0V4a1 1 0 00-1-1zm10.293 9.293a1 1 0 001.414 1.414l3-3a1 1 0 000-1.414l-3-3a1 1 0 10-1.414 1.414L14.586 9H7a1 1 0 100 2h7.586l-1.293 1.293z"
          clipRule="evenodd"
        />
      </svg>
    </button>
  );
}

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

export function Duel({
  duelId,
  duelToken,
  hostUrl,
  playmat,
  flipOpponentCards = false,
  resolveChatUser,
  renderChatUserTrigger,
  blockedChatUsers,
  devTools,
  onLeaveDuel,
  onDuelFinished,
  onNewTurn,
}: DuelProps) {
  const [action, setAction] = useState<ActionMessage | null>(null);
  const [actionRevision, setActionRevision] = useState(0);
  const [actionError, setActionError] = useState<ActionWarningMessage | null>(
    null,
  );
  const [chatMessages, setChatMessages] = useState<ReceivedChatMessage[]>([]);
  const [wait, setWait] = useState("");
  const [warningMessage, setWarningMessage] = useState("");
  const [dots, setDots] = useState<"." | ".." | "...">(".");
  const [duelFinishedCountdown, setDuelFinishedCountdown] = useState<
    number | null
  >(null);
  const [duelFinishedRedirectFailed, setDuelFinishedRedirectFailed] =
    useState(false);
  const [confirmForfeit, setConfirmForfeit] = useState(false);

  const {
    connected,
    error,
    duelFinished,
    send,
    sendJoinMatch,
    sendEndTurn,
    sendAddToBattlezone,
    sendAddToManazone,
    sendAttackPlayer,
    sendAttackCreature,
    sendTapAbility,
    sendResign,
    sendAction,
    sendChat,
    state,
    opponentDisconnected,
    reconnecting,
  } = useDuel({
    hostUrl,
    duelId,
    duelToken,
    onActionMessage: (nextAction) => {
      setAction(nextAction);
      setActionError(null);
      setActionRevision((revision) => revision + 1);
    },
    onActionError: setActionError,
    onActionClose: () => {
      setAction(null);
      setActionError(null);
    },
    onChat: (data) => {
      setChatMessages((prev) => [...prev, { ...data, receivedAt: Date.now() }]);
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

  // Whose turn it is only reaches the client as part of the match state, so a
  // new turn is the turn changing from the one previously seen rather than an
  // event of its own.
  const myTurn = state?.myTurn;
  const previousMyTurnRef = useRef<boolean | null>(null);

  useEffect(() => {
    if (myTurn === undefined) {
      return;
    }

    const previousMyTurn = previousMyTurnRef.current;
    previousMyTurnRef.current = myTurn;

    // The first state of the connection describes the turn already in progress,
    // whether the duel just started or the player reconnected mid-duel.
    if (previousMyTurn === null || previousMyTurn === myTurn) {
      return;
    }

    onNewTurn?.(myTurn);
  }, [myTurn, onNewTurn]);

  useEffect(() => {
    if (!duelFinished) {
      return;
    }

    setWarningMessage("");
    setWait("");
    setAction(null);
    setActionError(null);
    setConfirmForfeit(false);
    setDuelFinishedCountdown(5);
    setDuelFinishedRedirectFailed(false);

    let countdown = 5;
    const interval = window.setInterval(() => {
      countdown -= 1;

      if (countdown <= 0) {
        window.clearInterval(interval);
        setDuelFinishedCountdown(0);
        onDuelFinished?.(duelFinished);

        window.setTimeout(() => {
          setDuelFinishedRedirectFailed(true);
        }, 3000);
        return;
      }

      setDuelFinishedCountdown(countdown);
    }, 1000);

    return () => {
      window.clearInterval(interval);
    };
  }, [duelFinished, onDuelFinished]);

  const duelFinishedWinner =
    duelFinished?.winner?.username ?? duelFinished?.winner?.uid ?? "No winner";

  const duelFinishedStatusText =
    duelFinishedCountdown !== null && duelFinishedCountdown > 0
      ? `Redirecting in ${duelFinishedCountdown}...`
      : "Redirecting...";

  const isCompact = useCompactViewport();
  const [chatOpen, setChatOpen] = useState(false);

  // The drawer only exists on narrow screens; leaving it "open" while resizing
  // to desktop would otherwise strand the state.
  useEffect(() => {
    if (!isCompact) {
      setChatOpen(false);
    }
  }, [isCompact]);

  const [previewCard, setPreviewCard] = useState<PreviewCard | null>(null);
  const [multiCardView, setMultiCardView] = useState<{
    cards: { imageId: string; name: string }[];
    title: string;
  } | null>(null);

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
    threshold: number;
  } | null>(null);

  const handleCardDragStart = (
    virtualId: string,
    imageId: string,
    name: string | undefined,
    sourceZone: DragZone,
    rotated: boolean,
    e: React.PointerEvent,
  ) => {
    // Ignore right clicks
    if (e.button === 2) {
      return;
    }

    // Record start position but don't start dragging yet
    setDragStartPosition({
      x: e.clientX,
      y: e.clientY,
      virtualId,
      imageId,
      name,
      sourceZone,
      rotated,
      // A finger is far less precise than a mouse, so it needs more slack
      // before a tap counts as a drag.
      threshold: e.pointerType === "mouse" ? MOUSE_DRAG_THRESHOLD : TOUCH_DRAG_THRESHOLD,
    });
  };

  const handleMouseMove = (e: PointerEvent) => {
    const clientX = e.clientX;
    const clientY = e.clientY;

    // Check if we should start dragging based on threshold
    if (dragStartPosition && !dragState) {
      const deltaX = Math.abs(clientX - dragStartPosition.x);
      const deltaY = Math.abs(clientY - dragStartPosition.y);

      if (
        deltaX > dragStartPosition.threshold ||
        deltaY > dragStartPosition.threshold
      ) {
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
          : dragStartPosition.virtualId,
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

  const handlePointerCancel = () => {
    setDragState(null);
    setDropZone(null);
    setDragStartPosition(null);
  };

  useEffect(() => {
    if (dragState || dragStartPosition) {
      window.addEventListener("pointermove", handleMouseMove);
      window.addEventListener("pointerup", handleMouseUp);
      // A pointer can be taken away mid gesture, by the system or by a browser
      // scroll gesture winning. Without this the card would stay stuck to the
      // cursor with no way to drop it.
      window.addEventListener("pointercancel", handlePointerCancel);

      return () => {
        window.removeEventListener("pointermove", handleMouseMove);
        window.removeEventListener("pointerup", handleMouseUp);
        window.removeEventListener("pointercancel", handlePointerCancel);
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

  const isSpectating = state.spectator;

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
        (c) => c.virtualId === dragState.virtualId,
      );

      if (zone === "myPlayzone") {
        // Check if card can be played. Attacking or using a tap ability moves
        // the match into the attack step, after which nothing can be summoned.
        return draggedCard &&
          cardHasFlag(draggedCard.flags, PLAYABLE_FLAG) &&
          !state?.hasAttackedThisRound
          ? "green"
          : "red";
      }

      if (zone === "myManazone") {
        // Check if mana can be added
        return !state?.hasAddedManaThisRound &&
          !state?.hasAttackedThisRound &&
          state?.canChargeManaThisRound
          ? "green"
          : "red";
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

  /** The card actions and End turn button. On a narrow screen these move out of
   * the side column into a bar along the bottom, because they are the controls a
   * player needs constantly and the column itself is hidden there. */
  const renderControls = (compact: boolean) => {
    if (isSpectating) {
      return null;
    }

    // Attacking or activating a tap ability moves the match into the attack
    // step, and the server refuses summons, casts and mana charging from there
    // on. Disabling the buttons blocks the click outright, so the reason the
    // server would have replied with is surfaced as a tooltip instead.
    const cantSummon = state.hasAttackedThisRound;
    const summonDisabledTooltip = cantSummon
      ? "You can't summon creatures after attacking or using a tap ability"
      : "You don't have enough mana to play this card";

    const cantChargeMana =
      state.hasAddedManaThisRound ||
      state.hasAttackedThisRound ||
      !state.canChargeManaThisRound;
    const chargeManaDisabledTooltip = state.hasAddedManaThisRound
      ? "You've already added mana to your manazone this turn"
      : state.hasAttackedThisRound
        ? "You can't charge mana after attacking or using a tap ability"
        : "You can't charge mana after playing creatures or spells";

    const cardActions = selectedCard && state.myTurn && (
      <div className={compact ? "flex items-center gap-2" : "flex flex-col gap-2"}>
        <div
          className={`overflow-hidden text-ellipsis whitespace-nowrap text-xs ${
            compact
              ? // The board already highlights the selected card, so on a phone
                // the name yields to the buttons rather than squeezing them.
                "hidden min-w-0 shrink sm:block sm:max-w-[8rem]"
              : "flex-1"
          }`}
        >
          {selectedCard.name}
        </div>
        {selectedCard.zone === "hand" && (
          <div className="flex flex-1 gap-2">
            <div className="flex-1 min-w-0">
              <Button
                onClick={() => sendAddToBattlezone(selectedCard.virtualId)}
                disabled={!selectedCard.canPlay || cantSummon}
                disabledTooltip={summonDisabledTooltip}
              >
                Summon
              </Button>
            </div>
            <div className="flex-1 min-w-0">
              <Button
                onClick={() => sendAddToManazone(selectedCard.virtualId)}
                disabled={cantChargeMana}
                disabledTooltip={chargeManaDisabledTooltip}
              >
                {compact ? "Mana" : "Add to manazone"}
              </Button>
            </div>
          </div>
        )}

        {selectedCard.zone === "battlezone" && (
          <div className="flex flex-1 gap-2">
            <div className="flex-1 min-w-0">
              <Button
                onClick={() => sendAttackPlayer(selectedCard.virtualId)}
              >
                {compact ? "Attack" : "Attack Player"}
              </Button>
            </div>
            <div className="flex-1 min-w-0">
              <Button
                onClick={() => sendAttackCreature(selectedCard.virtualId)}
              >
                {compact ? "Battle" : "Attack Creature"}
              </Button>
            </div>
            {selectedCard.hasTapAbility && (
              <div className="flex-1 min-w-0">
                <Button
                  onClick={() => sendTapAbility(selectedCard.virtualId)}
                >
                  {compact ? "Tap" : "Tap Ability"}
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
    );

    if (compact) {
      return (
        <div className="flex items-center gap-2">
          <div className="min-w-0 flex-1">{cardActions}</div>
          <div className="w-[6.5rem] shrink-0">
            <Button
              onClick={sendEndTurn}
              disabled={!state.myTurn}
              disabledTooltip="It's not your turn"
            >
              End turn
            </Button>
          </div>
        </div>
      );
    }

    return (
      <>
        <div className="bg-black/50 p-2 rounded-md h-[72px] text-gray-400">
          {cardActions}
        </div>
        <div className="bg-black/30 p-2 rounded-md">
          <Button
            onClick={sendEndTurn}
            disabled={!state.myTurn}
            disabledTooltip="It's not your turn"
          >
            End turn
          </Button>
        </div>
      </>
    );
  };

  return (
    <>
      <style>{scrollbarStyles}</style>
      <div
        className="w-full h-screen text-white flex bg-[linear-gradient(45deg,rgb(29,33,42),rgb(20,16,21))] bg-cover bg-no-repeat gap-2 p-1 custom-scrollbar"
        style={{
          height: "100dvh",
          ...(playmat && {
            backgroundImage: `url(${JSON.stringify(playmat)}), linear-gradient(45deg, rgb(29, 33, 42), rgb(20, 16, 21))`,
            backgroundPosition: "center",
          }),
          ...(dragState && { cursor: "grabbing" }),
        }}
      >
        {/* On a narrow screen this column is replaced by a chat drawer and a
            bottom control bar, so the board can use the full width. */}
        <div
          className={
            isCompact
              ? "fixed inset-y-0 left-0 z-40 flex w-[86vw] max-w-[340px] flex-col gap-2 bg-[rgb(20,16,21)] p-2 shadow-2xl transition-transform duration-200" +
                (chatOpen ? " translate-x-0" : " -translate-x-full")
              : "w-[300px] flex flex-col gap-2"
          }
        >
          {isCompact && (
            <div className="flex shrink-0 items-center justify-between">
              <span className="text-sm font-semibold">Chat</span>
              <button
                type="button"
                onClick={() => setChatOpen(false)}
                aria-label="Close chat"
                className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-full bg-gray-800 text-gray-300 transition-colors hover:bg-gray-700"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                    clipRule="evenodd"
                  />
                </svg>
              </button>
            </div>
          )}

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
                      Host
                    </Button>
                  </div>
                  <div className="flex-1">
                    <Button
                      variant={
                        devTools.activePlayer === "guest" ? "default" : "gray"
                      }
                      onClick={() => devTools.onPlayerSwitch("guest")}
                    >
                      Guest
                    </Button>
                  </div>
                  <div className="flex-1">
                    <Button
                      variant={
                        devTools.activePlayer === "spectator"
                          ? "default"
                          : "gray"
                      }
                      onClick={() => devTools.onPlayerSwitch("spectator")}
                    >
                      Spectator
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
                        className="w-full bg-gray-800 text-white pl-2 pr-8 py-[0.4rem] rounded border border-gray-700 focus:outline-none focus:border-blue-500 text-xs appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2212%22%20height%3D%2212%22%20viewBox%3D%220%200%2012%2012%22%3E%3Cpath%20fill%3D%22%23ffffff%22%20d%3D%22M6%209L1%204h10z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:12px_12px] bg-[position:right_0.5rem_center] bg-no-repeat"
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
                            "card-selector",
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
            <Chat
              messages={chatMessages}
              onSendMessage={sendChat}
              resolveUser={resolveChatUser}
              renderUserTrigger={renderChatUserTrigger}
              blockedUsers={blockedChatUsers}
            />
          </div>

          {/* Card actions and End turn, unless a narrow screen has moved
              them to the bottom bar */}
          {!isCompact && renderControls(false)}
        </div>

        {/* Tapping the dimmed board closes the drawer */}
        {isCompact && chatOpen && (
          <div
            className="fixed inset-0 z-30 bg-black/60"
            onClick={() => setChatOpen(false)}
          />
        )}

        <div
          className="flex flex-1 flex-col h-full w-full"
          style={isCompact && !isSpectating ? { paddingBottom: "3.25rem" } : undefined}
        >
          <div className="h-[10%] relative" data-dropzone="opponentManazone">
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentManazone"
            />
            {getDropZoneColor("opponentManazone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("opponentManazone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full pb-1">
                {state.opponent.manazone.map(
                  CreateCard({
                    flipped: !flipOpponentCards,
                    dragState,
                    zone: "opponentManazone",
                    onRightClick: (imageId, name) =>
                      setPreviewCard({ imageId, name: name || "" }),
                  }),
                )}
              </div>
            </div>
          </div>
          <div
            className="h-[10%] w-full relative"
            data-dropzone="opponentShieldzone"
          >
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentShieldzone"
            />
            {getDropZoneColor("opponentShieldzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("opponentShieldzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full p-1">
                {state.opponent.shieldzone.map(
                  CreateCard({
                    flipped: !flipOpponentCards,
                    shieldMap: state.opponent.shieldMap,
                    dragState,
                    zone: "opponentShieldzone",
                    onRightClick: (imageId, name) =>
                      setPreviewCard({ imageId, name: name || "" }),
                  }),
                )}
              </div>
            </div>
          </div>
          <div
            className="h-[20%] w-full relative"
            data-dropzone="opponentPlayzone"
          >
            <div
              className="absolute inset-0 z-0"
              data-dropzone="opponentPlayzone"
            />
            {getDropZoneColor("opponentPlayzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("opponentPlayzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full p-1">
                {state.opponent.playzone.map(
                  CreateCard({
                    flipped: !flipOpponentCards,
                    dragState,
                    zone: "opponentPlayzone",
                    onRightClick: (imageId, name) =>
                      setPreviewCard({ imageId, name: name || "" }),
                  }),
                )}
              </div>
            </div>
          </div>
          <div className="h-[20%] w-full relative" data-dropzone="myPlayzone">
            <div className="absolute inset-0 z-0" data-dropzone="myPlayzone" />
            {getDropZoneColor("myPlayzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("myPlayzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full p-1">
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
                  }),
                )}
              </div>
            </div>
          </div>
          <div className="h-[10%] w-full relative" data-dropzone="myShieldzone">
            <div
              className="absolute inset-0 z-0 "
              data-dropzone="myShieldzone"
            />
            {getDropZoneColor("myShieldzone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("myShieldzone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full p-1">
                {state.me.shieldzone.map(
                  CreateCard({
                    shieldMap: state.me.shieldMap,
                    dragState,
                    zone: "myShieldzone",
                    onRightClick: (imageId, name) =>
                      setPreviewCard({ imageId, name: name || "" }),
                  }),
                )}
              </div>
            </div>
          </div>
          <div className="h-[10%] w-full relative" data-dropzone="myManazone">
            <div className="absolute inset-0 z-0" data-dropzone="myManazone" />
            {getDropZoneColor("myManazone") === "green" && (
              <div className="absolute inset-0 bg-green-500/30 pointer-events-none z-20" />
            )}
            {getDropZoneColor("myManazone") === "red" && (
              <div className="absolute inset-0 bg-red-500/30 pointer-events-none z-20" />
            )}
            <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
              <div className="inline-flex w-max justify-start gap-5 h-full p-1">
                {state.me.manazone.map(
                  CreateCard({
                    flipped: true,
                    dragState,
                    zone: "myManazone",
                    onRightClick: (imageId, name) =>
                      setPreviewCard({ imageId, name: name || "" }),
                  }),
                )}
              </div>
            </div>
          </div>
          <div className="h-[20%] w-full relative" data-dropzone="hand">
            {isSpectating ? (
              <div className="absolute inset-0 bg-black/30 rounded-md flex items-center justify-center">
                <div className="flex flex-col items-center gap-3 text-center">
                  <span className="text-xs font-semibold uppercase tracking-widest text-gray-400">
                    Spectating
                  </span>
                  <div className="flex items-center gap-8">
                    <span className="text-sm font-semibold text-white">
                      {state.me.username}
                    </span>
                    <span className="text-gray-500 text-lg font-bold">vs</span>
                    <span className="text-sm font-semibold text-white">
                      {state.opponent.username}
                    </span>
                  </div>
                </div>
              </div>
            ) : (
              <>
                <div className="absolute inset-0 z-0" data-dropzone="hand" />
                <div className="absolute inset-0 z-10 overflow-x-auto overflow-y-hidden">
                  <div className="inline-flex w-max justify-start gap-5 h-full pt-1 p-px">
                    {state.me.hand.map(
                      CreateCard({
                        selected: (id: string) => id === selectedCardId,
                        interactable: state?.myTurn,
                        canAddToManazone:
                          !state.hasAddedManaThisRound &&
                          !state.hasAttackedThisRound &&
                          state.canChargeManaThisRound,
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
                      }),
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
        </div>

        {/* Bottom control bar, narrow screens only. The board reserves room for
            it so it never covers the hand. */}
        {isCompact && !isSpectating && (
          <div className="fixed inset-x-0 bottom-0 z-30 flex h-[3.25rem] items-center gap-2 border-t border-white/10 bg-[rgb(20,16,21)]/95 px-2">
            <button
              type="button"
              onClick={() => setChatOpen(true)}
              aria-label="Open chat"
              className="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-md bg-gray-800 text-gray-200 transition-colors hover:bg-gray-700"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M18 10c0 3.866-3.582 7-8 7a8.84 8.84 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7z"
                  clipRule="evenodd"
                />
              </svg>
            </button>

            <div className="min-w-0 flex-1">{renderControls(true)}</div>

            {/* Forfeit sits here rather than in the top corner so every control
                a player needs is along one edge, within thumb reach. */}
            <ForfeitButton onClick={() => setConfirmForfeit(true)} />
          </div>
        )}

        {/* Spectators get the chat toggle without the control bar */}
        {isCompact && isSpectating && (
          <button
            type="button"
            onClick={() => setChatOpen(true)}
            aria-label="Open chat"
            className="fixed bottom-2 left-2 z-30 flex h-10 w-10 cursor-pointer items-center justify-center rounded-full bg-gray-800 text-gray-200 shadow-lg transition-colors hover:bg-gray-700"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fillRule="evenodd"
                d="M18 10c0 3.866-3.582 7-8 7a8.84 8.84 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        )}

        {/* Forfeit - Top Right. On a narrow screen it lives in the bottom bar
            instead, beside End turn. */}
        {!isSpectating && !isCompact && (
          <div className="fixed right-[0.5vw] top-[0.5vh] z-30">
            <ForfeitButton onClick={() => setConfirmForfeit(true)} />
          </div>
        )}

        {/* Player Info Panel - Right Side */}
        <div className="fixed right-[0.5vw] top-1/2 -translate-y-2/3 w-[16vw] min-w-[56px] max-w-[80px] min-[1200px]:w-[12vw] min-[1200px]:min-w-[100px] min-[1200px]:max-w-[160px] flex flex-col gap-[0.5vh] md:gap-[5vh] z-20">
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
              <div className="relative h-[12vh] min-h-[40px] max-h-[70px] min-[1200px]:min-h-[60px] min-[1200px]:max-h-[110px] flex items-center justify-center">
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
              <div className="relative h-[12vh] min-h-[40px] max-h-[70px] min-[1200px]:min-h-[60px] min-[1200px]:max-h-[110px] flex items-center justify-center">
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
            {/* Player Hand Count (spectator only) */}
            {isSpectating && (
              <div className="text-center">
                <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh]">
                  Hand [{state.me.handCount}]
                </p>
              </div>
            )}
            {/* Player Graveyard */}
            <div>
              <p className="text-[clamp(0.6rem,1.2vh,0.85rem)] text-white mb-[0.5vh] text-center">
                Graveyard [{state.me.graveyard.length}]
              </p>
              <div className="relative h-[12vh] min-h-[40px] max-h-[70px] min-[1200px]:min-h-[60px] min-[1200px]:max-h-[110px] flex items-center justify-center">
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
              <div className="relative h-[12vh] min-h-[40px] max-h-[70px] min-[1200px]:min-h-[60px] min-[1200px]:max-h-[110px] flex items-center justify-center">
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

      <Popup
        visible={!!duelFinished}
        title="Duel Finished"
        maxWidth="500px"
        closeOnOutsideClick={false}
        showCloseButton={false}
      >
        <div className="p-6 text-white">
          <p className="text-lg font-semibold">
            {duelFinishedWinner} won the duel
          </p>
          <p className="mt-3 text-sm text-gray-400">{duelFinishedStatusText}</p>
          {duelFinishedRedirectFailed && (
            <p className="mt-3 text-sm text-gray-400">
              Failed to redirect, please close the page manually
            </p>
          )}
        </div>
      </Popup>

      <Popup
        visible={confirmForfeit}
        onClose={() => setConfirmForfeit(false)}
        title="Forfeit"
        maxWidth="500px"
        closeOnOutsideClick={true}
      >
        <div className="p-6 text-white">
          <p>
            Are you sure you want to forfeit? Your opponent wins the duel
            immediately.
          </p>
          <div className="flex gap-3 mt-6">
            <Button variant="gray" onClick={() => setConfirmForfeit(false)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                setConfirmForfeit(false);
                sendResign();
              }}
            >
              Forfeit
            </Button>
          </div>
        </div>
      </Popup>

      <Popup
        visible={reconnecting}
        title="Disconnected"
        maxWidth="500px"
        closeOnOutsideClick={false}
        showCloseButton={false}
      >
        <div className="p-6 text-white ">
          <p>Attempting to reconnect{dots}</p>
          <div className="flex mt-6">
            <Button variant="destructive" onClick={onLeaveDuel}>
              Leave Duel
            </Button>
          </div>
        </div>
      </Popup>

      <Popup
        visible={opponentDisconnected}
        title="Opponent Disconnected"
        maxWidth="500px"
        closeOnOutsideClick={false}
        showCloseButton={false}
      >
        <div className="p-6 text-white">
          <p>
            Your opponent disconnected or left the match. Waiting for them to
            reconnect{dots}
          </p>
          <div className="flex mt-6">
            <Button variant="destructive" onClick={onLeaveDuel}>
              Leave Duel
            </Button>
          </div>
        </div>
      </Popup>

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
          key={actionRevision}
          title={action.showCards ? "Card Preview" : "Action Required"}
          visible={true}
          error={actionError ? actionError.message : undefined}
          actionType={action.actionType}
          cards={Array.isArray(action.cards) ? action.cards : undefined}
          cardsObject={
            typeof action.cards === "object" && !Array.isArray(action.cards)
              ? action.cards
              : undefined
          }
          showCards={action.showCards}
          text={action.text}
          minSelections={action.minSelections || 0}
          maxSelections={action.maxSelections || 0}
          cancellable={action.cancellable || false}
          unselectableCards={action.unselectableCards}
          choices={action.choices || null}
          onChoose={sendAction}
          onClose={() => sendAction({ cards: [], cancel: true })}
          onDismiss={() => setAction(null)}
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
    /** virtualId -> shield number, as sent by the server in the match state. */
    shieldMap?: Record<string, number>;
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
      e: React.PointerEvent,
    ) => void;
    onRightClick?: (imageId: string, name?: string) => void;
  } = {},
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
        number={options.shieldMap?.[card.virtualId]}
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
              e,
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
