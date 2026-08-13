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
  const touchInProgressRef = useRef(false);

  const handleBrushEnd = () => {
    setIsBrushing(false);
    brushedCardIdsRef.current.clear();
    // Reset touch flag after a delay to ensure mouse events are blocked
    window.setTimeout(() => {
      touchInProgressRef.current = false;
    }, 300);
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

  const handleCardMouseDown = (
    cardId: string,
    e?: React.MouseEvent | React.TouchEvent
  ) => {
    // Ignore mouse events if a touch is in progress
    if (e && !("touches" in e) && touchInProgressRef.current) {
      return;
    }

    if (e && "touches" in e) {
      touchInProgressRef.current = true;
    }

    setIsBrushing(true);
    toggleCard(cardId);
  };

  const handleCardHover = (cardId: string) => {
    if (!isBrushing) return;
    toggleCard(cardId);
  };

  useEffect(() => {
    if (isBrushing) {
      const handleMouseUp = () => handleBrushEnd();
      const handleTouchEnd = () => handleBrushEnd();

      window.addEventListener("mouseup", handleMouseUp);
      window.addEventListener("touchend", handleTouchEnd);

      return () => {
        window.removeEventListener("mouseup", handleMouseUp);
        window.removeEventListener("touchend", handleTouchEnd);
      };
    }
  }, [isBrushing]);

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isBrushing) return;

    const touch = e.touches[0];
    const element = document.elementFromPoint(touch.clientX, touch.clientY);
    const cardElement = element?.closest("[data-card-id]") as HTMLElement;

    if (cardElement) {
      const cardId = cardElement.getAttribute("data-card-id");
      if (cardId) {
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
  const gridCols = Math.max(3, Math.min(cardCount, 6));

  const currentSelectableCardIds = new Set(
    (cardsObject ? Object.values(cardsObject).flat() : cards || []).map(
      (card) => card.virtualId,
    ),
  );
  const selectedCurrentCardIds = [...selectedCardIds].filter((cardId) =>
    currentSelectableCardIds.has(cardId),
  );

  const cardGridStyle = {
    gridTemplateColumns: `repeat(${gridCols}, minmax(0, 1fr))`,
  };
  const isCardSelection = actionType === ActionType.None || !actionType;

  return (
    <Popup
      title={title}
      visible={visible}
      showCloseButton={showCards ? true : cancellable}
      zIndex={1000}
      closeOnOutsideClick={false}
      contentClassName="flex min-h-0 flex-1 overflow-hidden"
      onClose={showCards ? onDismiss : onClose}
    >
      <div className="flex min-h-0 flex-1 flex-col select-none">
        <div
          className="custom-scrollbar min-h-0 flex-1 overflow-y-auto px-6 py-4"
          onTouchMove={handleTouchMove}
        >
          <div className="text-sm text-gray-100">{text}</div>

          {actionType === ActionType.ShowCards && (
            <div
              className="mt-4 grid w-fit gap-2 rounded-md bg-black/30 p-2"
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
                    className="rounded-md"
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
                className="mt-4 grid w-fit gap-2 rounded-md bg-black/30 p-2"
                style={cardGridStyle}
              >
                {cardsObject &&
                  selectedCardsObjectCards.map((card) => (
                    <div
                      key={card.virtualId}
                      className="w-full"
                      data-card-id={card.virtualId}
                      onMouseEnter={() => handleCardHover(card.virtualId)}
                      onMouseDown={(event) =>
                        handleCardMouseDown(card.virtualId, event)
                      }
                      onTouchStart={(event) =>
                        handleCardMouseDown(card.virtualId, event)
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
                        className={`rounded-md ${
                          selectedCardIds.has(card.virtualId)
                            ? "ring-1 ring-blue-100"
                            : ""
                        }`}
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                    </div>
                  ))}

                {cardsObject && !selectedCardsObjectCards.length && (
                  <div className="text-sm text-gray-400">No cards to show</div>
                )}

                {!cardsObject &&
                  cards?.map((card) => (
                    <div
                      key={card.virtualId}
                      className="w-full"
                      data-card-id={card.virtualId}
                      onMouseEnter={() => handleCardHover(card.virtualId)}
                      onMouseDown={(event) =>
                        handleCardMouseDown(card.virtualId, event)
                      }
                      onTouchStart={(event) =>
                        handleCardMouseDown(card.virtualId, event)
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
                        className={`rounded-md ${
                          selectedCardIds.has(card.virtualId)
                            ? "ring-1 ring-blue-100"
                            : ""
                        }`}
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                    </div>
                  ))}

                {!cardsObject &&
                  unselectableCards?.map((card) => (
                    <div key={card.virtualId} className="w-full">
                      <img
                        onContextMenu={(event) => {
                          event.preventDefault();
                          if (onCardRightClick && card.uid) {
                            onCardRightClick(card.uid, card.name);
                          }
                        }}
                        draggable={false}
                        className="cursor-not-allowed rounded-md opacity-50 grayscale"
                        src={`https://scans.shobu.io/${card.uid}.jpg`}
                        alt={card.name}
                        style={{ borderRadius: "5%" }}
                      />
                    </div>
                  ))}
              </div>
            </>
          )}

          {actionType === ActionType.Order && (
            <div
              className="mt-4 grid w-fit gap-2 rounded-md bg-black/30 p-2"
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
                    onMouseEnter={() => handleCardHover(card.virtualId)}
                    onMouseDown={(event) =>
                      handleCardMouseDown(card.virtualId, event)
                    }
                    onTouchStart={(event) =>
                      handleCardMouseDown(card.virtualId, event)
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
                      className={`rounded-md ${
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

        <div className="shrink-0 border-t border-gray-800 bg-gray-900 px-6 py-4">
          {actionType === ActionType.ShowCards && (
            <div className="flex">
              <Button onClick={() => onDismiss?.()}>
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
              <div className="flex-1 text-right text-xs italic text-gray-300">
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
