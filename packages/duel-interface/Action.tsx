import { useEffect, useRef, useState } from "react";
import { Popup } from "./Popup";
import { ActionType, CardState } from "./types";
import { Button } from "./Button";
import { CountInput } from "./CountInput";

export interface ActionProps {
  title: string;
  visible: boolean;
  error?: string;
  onChoose: (data: {
    cards: string[];
    cancel: boolean;
    count?: number;
  }) => void;
  onClose: () => void;
  onDismiss?: () => void;
  onCardRightClick?: (imageId: string, name?: string) => void;

  // Action details
  actionType?: ActionType;
  cards?: CardState[];
  cardsObject?: Record<string, CardState[]>;
  /** virtualId -> shield number, as shown in the shield zone. */
  shieldMap?: Record<string, number>;
  showCards?: {
    cards: string[];
    dismissable: boolean;
  };
  text: string;
  minSelections: number;
  maxSelections: number;
  cancellable: boolean;
  unselectableCards?: CardState[];
  choices: string[] | null;
}

export function Action({
  title,
  text,
  error,
  cards,
  unselectableCards,
  cardsObject,
  shieldMap,
  showCards,
  cancellable,
  visible,
  onChoose,
  onClose,
  onDismiss,
  onCardRightClick,
  minSelections,
  maxSelections,
  actionType,
  choices,
}: ActionProps) {
  if (actionType === ActionType.Order) {
    minSelections = cards ? cards.length : 0;
    maxSelections = cards ? cards.length : 0;
  }

  const [selectedCardsObjectKey, setSelectedCardsObjectKey] = useState<
    string | null
  >(cardsObject ? Object.keys(cardsObject)[0] : null);
  const [selectedCardIds, setSelectedCardIds] = useState(new Set<string>());
  const [count, setCount] = useState(minSelections);
  const [isBrushing, setIsBrushing] = useState(false);
  const [selectedSearchValue, setSelectedSearchValue] = useState("-1");
  const brushedCardIdsRef = useRef(new Set<string>());
  // A finger that lands on a card might be starting a brush or might be about
  // to scroll the list, so the first card waits here until the gesture says
  // which. Releasing without moving is a tap and selects it; reaching a second
  // card is a brush; the browser taking the gesture for a scroll drops it.
  const pendingCardIdRef = useRef<string | null>(null);

  const flushPendingCard = () => {
    const pendingCardId = pendingCardIdRef.current;

    if (pendingCardId === null) {
      return;
    }

    pendingCardIdRef.current = null;
    toggleCard(pendingCardId);
  };

  const handleBrushEnd = () => {
    flushPendingCard();
    setIsBrushing(false);
    brushedCardIdsRef.current.clear();
  };

  // The browser claims the gesture once it starts scrolling. Whatever the
  // finger was resting on was never a selection, so it is dropped rather than
  // flushed, and a scroll that began on a card leaves the selection alone.
  const handleBrushCancel = () => {
    pendingCardIdRef.current = null;
    setIsBrushing(false);
    brushedCardIdsRef.current.clear();
  };

  const toggleCard = (cardId: string) => {
    // Only toggle each card once per brush session
    if (brushedCardIdsRef.current.has(cardId)) return;

    brushedCardIdsRef.current.add(cardId);

    setSelectedCardIds((previous) => {
      const next = new Set(previous);

      // Always allow deselection
      if (next.has(cardId)) {
        next.delete(cardId);
      } else if (next.size < maxSelections) {
        next.add(cardId);
      }

      return next;
    });
  };

  // One pointer stream for mouse, touch and pen. Binding mouse and touch
  // separately meant a tap ran both paths, because a touch is followed by
  // compatibility mouse events, and the second toggle undid the first.
  const handleCardPointerDown = (cardId: string, e: React.PointerEvent) => {
    setIsBrushing(true);

    // A mouse press is unambiguous: there is no scrolling to compete with it.
    if (e.pointerType === "mouse") {
      toggleCard(cardId);
      return;
    }

    pendingCardIdRef.current = cardId;
  };

  const handleCardHover = (cardId: string) => {
    if (!isBrushing) return;
    toggleCard(cardId);
  };

  useEffect(() => {
    if (isBrushing) {
      window.addEventListener("pointerup", handleBrushEnd);
      window.addEventListener("pointercancel", handleBrushCancel);

      return () => {
        window.removeEventListener("pointerup", handleBrushEnd);
        window.removeEventListener("pointercancel", handleBrushCancel);
      };
    }
  }, [isBrushing]);

  // A touch pointer is implicitly captured by the element it started on, so
  // pointerenter never fires for the cards it is dragged across. Hit testing the
  // coordinates is what makes brushing work with a finger.
  const handlePointerMove = (e: React.PointerEvent) => {
    if (!isBrushing || e.pointerType === "mouse") return;

    const element = document.elementFromPoint(e.clientX, e.clientY);
    const cardElement = element?.closest("[data-card-id]") as HTMLElement;

    if (cardElement) {
      const cardId = cardElement.getAttribute("data-card-id");
      if (cardId) {
        // Still on the card the finger landed on. The gesture has not proven
        // itself a brush yet, so that card stays pending.
        if (pendingCardIdRef.current === cardId) {
          return;
        }

        flushPendingCard();
        handleCardHover(cardId);
      }
    }
  };

  const selectedCardsObjectCards =
    cardsObject?.[selectedCardsObjectKey || ""] || [];
  const cardCount = cardsObject
    ? selectedCardsObjectCards.length
    : showCards
      ? showCards.cards.length
      : (cards?.length || 0) + (unselectableCards?.length || 0);
  // The grid used a fixed column count at natural card width, which overflowed
  // the popup horizontally on a phone with no way to scroll to the cards that
  // fell off the edge. Auto-fitting to a minimum card width instead means the
  // column count follows the space actually available.
  //
  // The cap below exists only to stop a handful of cards from being stretched
  // across the whole popup, so it scales with the number of cards. The popup
  // itself is then sized to the same figure, so the grid always reaches both
  // edges of the content area instead of leaving the cap's slack as dead space.
  const gridMaxCardsAcross = Math.max(3, cardCount);
  const gridMaxWidthRem = gridMaxCardsAcross * 8;
  // The content area's own horizontal padding, which the popup has to add back
  // on top of the grid for the two to line up.
  const contentPaddingRem = 3;

  const currentSelectableCardIds = new Set(
    (cardsObject ? Object.values(cardsObject).flat() : cards || []).map(
      (card) => card.virtualId,
    ),
  );
  const selectedCurrentCardIds = [...selectedCardIds].filter((cardId) =>
    currentSelectableCardIds.has(cardId),
  );

  const cardGridStyle = {
    // auto-fit keeps cards from ever being pushed outside the popup: columns
    // shrink to the 5rem floor and then wrap onto the next row. The floor is
    // capped at 100% so that a popup narrower than one card still cannot be
    // overflowed by the track itself.
    gridTemplateColumns: `repeat(auto-fit, minmax(min(5rem, 100%), 1fr))`,
    maxWidth: `${gridMaxWidthRem}rem`,
  };
  const isCardSelection = actionType === ActionType.None || !actionType;

  // A non-dismissable preview (show_cards_non_dismissible) blocks the match
  // event loop until the player answers it, so acknowledging one has to reach
  // the server. A dismissable preview (show_cards) is fire and forget, and
  // answering it would leave a stray action for the next prompt to consume.
  const acknowledgeShowCards = () => {
    if (showCards && !showCards.dismissable) {
      onClose();
    }

    onDismiss?.();
  };

  return (
    <Popup
      title={title}
      visible={visible}
      showCloseButton={showCards ? true : cancellable}
      zIndex={1000}
      closeOnOutsideClick={false}
      // Prompts that show cards earn the extra room on a large monitor, where
      // the old fixed 600px left the cards needlessly small. A prompt with few
      // enough cards to hit the grid's width cap shrinks to that instead, so
      // the cards reach both edges rather than trailing off into empty popup.
      maxWidth={
        cardCount > 0
          ? `min(96vw, ${gridMaxWidthRem + contentPaddingRem}rem, 900px)`
          : "min(96vw, 600px)"
      }
      contentClassName="flex min-h-0 flex-1 overflow-hidden"
      onClose={showCards ? acknowledgeShowCards : onClose}
    >
      {/* min-w-0 is what keeps the cards inside the popup. A flex item defaults
          to min-width:auto, so this column refused to be narrower than the card
          grid's own minimum width, grew past the popup, and everything beyond
          the edge was clipped by the popup's overflow-hidden. */}
      <div className="flex min-h-0 min-w-0 flex-1 flex-col select-none">
        <div
          className="custom-scrollbar min-h-0 flex-1 overflow-y-auto px-4 py-4 sm:px-6"
          onPointerMove={handlePointerMove}
        >
          <div className="text-sm text-gray-100">{text}</div>

          {actionType === ActionType.ShowCards && (
            <div
              className="mt-4 grid w-full gap-2 rounded-md bg-black/30 p-2"
              style={cardGridStyle}
            >
              {showCards?.cards.map((imageId, index) => (
                <div key={index} className="w-full">
                  <img
                    onContextMenu={(event) => {
                      event.preventDefault();
                      onCardRightClick?.(imageId, "");
                    }}
                    draggable={false}
                    className="w-full rounded-md"
                    src={`https://scans.shobu.io/${imageId}.jpg`}
                    alt="Card preview"
                    style={{ borderRadius: "5%" }}
                  />
                </div>
              ))}
            </div>
          )}

          {isCardSelection && (
            <>
              {cardsObject && (
                <div className="mt-4 flex gap-4">
                  <select
                    className="bg-gray-800 text-white pl-2 pr-8 py-[0.4rem] rounded border border-gray-700 focus:outline-none focus:border-blue-500 text-xs appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2212%22%20height%3D%2212%22%20viewBox%3D%220%200%2012%2012%22%3E%3Cpath%20fill%3D%22%23ffffff%22%20d%3D%22M6%209L1%204h10z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:12px_12px] bg-[position:right_0.5rem_center] bg-no-repeat"
                    id="action-card-group-selector"
                    value={selectedCardsObjectKey || ""}
                    onChange={(event) =>
                      setSelectedCardsObjectKey(event.target.value)
                    }
                  >
                    {Object.keys(cardsObject).map((key) => (
                      <option key={key} value={key}>
                        {key}
                      </option>
                    ))}
                  </select>
                </div>
              )}

              <div
                className="mt-4 grid w-full gap-2 rounded-md bg-black/30 p-2"
                style={cardGridStyle}
              >
                {cardsObject &&
                  selectedCardsObjectCards.map((card) => (
                    <div
                      key={card.virtualId}
                      className="relative w-full"
                      data-card-id={card.virtualId}
                      // Brushing runs sideways along a row while the list
                      // scrolls up and down, so the browser keeps the vertical
                      // pan and leaves horizontal movement to the brush.
                      style={{ touchAction: "pan-y" }}
                      onPointerEnter={() => handleCardHover(card.virtualId)}
                      onPointerDown={(event) =>
                        handleCardPointerDown(card.virtualId, event)
                      }
                    >
                      <img
                        onContextMenu={(event) => {
                          event.preventDefault();
                          if (onCardRightClick && card.uid) {
                            onCardRightClick(card.uid, card.name);
                          }
                        }}
                        onDragStart={(event) => event.preventDefault()}
                        draggable={false}
                        className={`w-full rounded-md ${
                          selectedCardIds.has(card.virtualId)
                            ? "ring-1 ring-blue-100"
                            : ""
                        }`}
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                      <ShieldNumber shieldMap={shieldMap} card={card} />
                    </div>
                  ))}

                {cardsObject && !selectedCardsObjectCards.length && (
                  <div className="text-sm text-gray-400">No cards to show</div>
                )}

                {!cardsObject &&
                  cards?.map((card) => (
                    <div
                      key={card.virtualId}
                      className="relative w-full"
                      data-card-id={card.virtualId}
                      style={{ touchAction: "pan-y" }}
                      onPointerEnter={() => handleCardHover(card.virtualId)}
                      onPointerDown={(event) =>
                        handleCardPointerDown(card.virtualId, event)
                      }
                    >
                      <img
                        onContextMenu={(event) => {
                          event.preventDefault();
                          if (onCardRightClick && card.uid) {
                            onCardRightClick(card.uid, card.name);
                          }
                        }}
                        onDragStart={(event) => event.preventDefault()}
                        draggable={false}
                        className={`w-full rounded-md ${
                          selectedCardIds.has(card.virtualId)
                            ? "ring-1 ring-blue-100"
                            : ""
                        }`}
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                      <ShieldNumber shieldMap={shieldMap} card={card} />
                    </div>
                  ))}

                {!cardsObject &&
                  unselectableCards?.map((card) => (
                    <div key={card.virtualId} className="relative w-full">
                      <img
                        onContextMenu={(event) => {
                          event.preventDefault();
                          if (onCardRightClick && card.uid) {
                            onCardRightClick(card.uid, card.name);
                          }
                        }}
                        draggable={false}
                        className="w-full cursor-not-allowed rounded-md opacity-50 grayscale"
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                      <ShieldNumber shieldMap={shieldMap} card={card} />
                    </div>
                  ))}
              </div>
            </>
          )}

          {actionType === ActionType.Order && (
            <div
              className="mt-4 grid w-full gap-2 rounded-md bg-black/30 p-2"
              style={cardGridStyle}
            >
              {cards?.map((card) => {
                const orderNumber =
                  [...selectedCardIds].indexOf(card.virtualId) + 1;
                const isSelected = selectedCardIds.has(card.virtualId);

                return (
                  <div
                    key={card.virtualId}
                    className="relative w-full"
                    data-card-id={card.virtualId}
                    style={{ touchAction: "pan-y" }}
                    onPointerEnter={() => handleCardHover(card.virtualId)}
                    onPointerDown={(event) =>
                      handleCardPointerDown(card.virtualId, event)
                    }
                  >
                    <img
                      onContextMenu={(event) => {
                        event.preventDefault();
                        if (onCardRightClick && card.uid) {
                          onCardRightClick(card.uid, card.name);
                        }
                      }}
                      onDragStart={(event) => event.preventDefault()}
                      draggable={false}
                      className={`w-full rounded-md ${
                        isSelected ? "ring-1 ring-blue-100" : ""
                      }`}
                      src={`https://scans.shobu.io/${card.uid}.jpg`}
                      alt={card.name}
                      style={{ borderRadius: "5%" }}
                    />
                    {isSelected && (
                      <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
                        <div className="text-6xl font-bold text-white drop-shadow-[0_0_8px_rgba(0,0,0,0.9)]">
                          {orderNumber}
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>

        <div className="shrink-0 border-t border-gray-800 bg-gray-900 px-4 py-3 sm:px-6 sm:py-4">
          {actionType === ActionType.ShowCards && (
            <div className="flex">
              <Button onClick={acknowledgeShowCards}>
                Acknowledge and Close
              </Button>
            </div>
          )}

          {(isCardSelection || actionType === ActionType.Order) && (
            <div className="flex items-center gap-4">
              <Button
                onClick={() =>
                  onChoose({
                    cards: selectedCurrentCardIds,
                    cancel: false,
                    count: 0,
                  })
                }
              >
                Choose
              </Button>
              {cancellable && (
                <Button variant="gray" onClick={onClose}>
                  Close
                </Button>
              )}
              <div className="hidden flex-1 text-right text-xs italic text-gray-300 sm:block">
                Click and drag to (de)select faster
              </div>
            </div>
          )}

          {actionType === ActionType.Question &&
            (choices && choices.length > 0 ? (
              <div className="flex flex-wrap gap-4">
                {choices.map((choice, index) => (
                  <div key={index} className="flex-1">
                    <Button
                      onClick={() =>
                        onChoose({ cards: [], cancel: false, count: index })
                      }
                    >
                      {choice}
                    </Button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-wrap gap-4">
                <div className="flex-1">
                  <Button
                    onClick={() => onChoose({ cards: [], cancel: false })}
                  >
                    Yes
                  </Button>
                </div>
                <div className="flex-1">
                  <Button
                    variant="gray"
                    onClick={() => onChoose({ cards: [], cancel: true })}
                  >
                    No
                  </Button>
                </div>
              </div>
            ))}

          {actionType === ActionType.Count && (
            <div className="flex gap-4">
              <CountInput
                value={count}
                onChange={setCount}
                min={minSelections}
                max={maxSelections}
              />
              <Button
                onClick={() => onChoose({ cards: [], cancel: false, count })}
              >
                Choose
              </Button>
            </div>
          )}

          {actionType === ActionType.Searchable &&
            choices &&
            choices.length > 0 && (
              <div className="flex gap-4">
                <select
                  className="bg-gray-800 text-white pl-2 pr-8 py-[0.4rem] rounded border border-gray-700 focus:outline-none focus:border-blue-500 text-xs appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2212%22%20height%3D%2212%22%20viewBox%3D%220%200%2012%2012%22%3E%3Cpath%20fill%3D%22%23ffffff%22%20d%3D%22M6%209L1%204h10z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:12px_12px] bg-[position:right_0.5rem_center] bg-no-repeat"
                  id="action-searchable-selector"
                  value={selectedSearchValue}
                  onChange={(event) =>
                    setSelectedSearchValue(event.target.value)
                  }
                >
                  <option value="-1" disabled>
                    Search and select
                  </option>
                  {choices.map((choice, index) => (
                    <option key={index} value={`${index}`}>
                      {choice}
                    </option>
                  ))}
                </select>
                <Button
                  disabled={selectedSearchValue === "-1"}
                  onClick={() => {
                    if (selectedSearchValue !== "-1") {
                      onChoose({
                        cards: [],
                        count: parseInt(selectedSearchValue),
                        cancel: false,
                      });
                    }
                  }}
                >
                  Choose
                </Button>
              </div>
            )}

          {error && (
            <div className="mt-3 flex items-center gap-2 text-sm text-red-500">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                className="h-5 w-5 shrink-0"
              >
                <path
                  fillRule="evenodd"
                  d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-5a.75.75 0 01.75.75v4.5a.75.75 0 01-1.5 0v-4.5A.75.75 0 0110 5zm0 10a1 1 0 100-2 1 1 0 000 2z"
                  clipRule="evenodd"
                />
              </svg>
              {error}
            </div>
          )}
        </div>
      </div>
    </Popup>
  );
}

// ShieldNumber mirrors the badge the shield zone draws in Card.tsx, so a
// shield offered as a selection candidate (for example when choosing which
// one to break) keeps the same index the player already knows it by.
function ShieldNumber({
  shieldMap,
  card,
}: {
  shieldMap?: Record<string, number>;
  card: CardState;
}) {
  const number = shieldMap?.[card.virtualId];

  if (number === undefined) {
    return null;
  }

  return (
    <span className="pointer-events-none absolute top-[2%] right-[10%] text-[clamp(0.55rem,1.1vh,0.8rem)] leading-none text-white drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]">
      {number}
    </span>
  );
}
