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
  onCardRightClick?: (imageId: string, name?: string) => void;

  // Action details
  actionType: ActionType;
  cards?: CardState[];
  text: string;
  minSelections: number;
  maxSelections: number;
  cancellable: boolean;
  unselectableCards: CardState[];
  choices: string[] | null;
}

export function Action({
  title,
  text,
  error,
  cards,
  cancellable,
  visible,
  onChoose,
  onClose,
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

  const [selectedCardIds, setSelectedCardIds] = useState(new Set<string>());
  const [count, setCount] = useState(minSelections);
  const [isBrushing, setIsBrushing] = useState(false);
  const [brushedCards, setBrushedCards] = useState(new Set<string>());
  const [selectedSearchValue, setSelectedSearchValue] = useState("-1");
  const touchInProgressRef = useRef(false);

  const handleBrushEnd = () => {
    setIsBrushing(false);
    setBrushedCards(new Set());
    // Reset touch flag after a delay to ensure mouse events are blocked
    setTimeout(() => {
      touchInProgressRef.current = false;
    }, 300);
  };

  const toggleCard = (cardId: string) => {
    // Only toggle each card once per brush session
    if (brushedCards.has(cardId)) return;

    setBrushedCards((prev) => new Set(prev).add(cardId));

    // Toggle the card
    if (selectedCardIds.has(cardId)) {
      // Always allow deselection
      setSelectedCardIds((prev) => {
        const next = new Set(prev);
        next.delete(cardId);
        return next;
      });
    } else {
      // Only allow selection if under max
      if (selectedCardIds.size < maxSelections) {
        setSelectedCardIds((prev) => new Set(prev).add(cardId));
      }
    }
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

  const cardCount = cards?.length || 0;
  const gridCols = Math.max(3, Math.min(cardCount, 6));

  return (
    <Popup
      title={title}
      visible={visible}
      showCloseButton={cancellable}
      zIndex={1000}
      closeOnOutsideClick={false}
      onClose={onClose}
    >
      <div className="px-6 py-6 pt-4 select-none" onTouchMove={handleTouchMove}>
        {/* Normal card selection */}
        {actionType === ActionType.None && (
          <>
            <div className="text-sm text-gray-100">{text}</div>
            <div
              className="grid gap-2 p-2 mt-4 bg-black/30 rounded-md w-fit"
              style={{
                gridTemplateColumns: `repeat(${gridCols}, minmax(0, 1fr))`,
              }}
            >
              {cards?.map((card, index) => (
                <div
                  key={index}
                  className="w-full"
                  data-card-id={card.virtualId}
                  onMouseEnter={() => handleCardHover(card.virtualId)}
                  onMouseDown={(e) => handleCardMouseDown(card.virtualId, e)}
                  onTouchStart={(e) => handleCardMouseDown(card.virtualId, e)}
                >
                  <img
                    onContextMenu={(e) => {
                      e.preventDefault();
                      if (onCardRightClick && card.uid) {
                        onCardRightClick(card.uid, card.name);
                      }
                    }}
                    onDragStart={(e) => e.preventDefault()}
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
            </div>
            <div className="flex items-center gap-4 mt-4">
              <Button
                onClick={() =>
                  onChoose({
                    cards: (cards || [])
                      .map((card) => card.virtualId)
                      .filter((id) => selectedCardIds.has(id)),
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
              <div className="flex-1 text-right text-xs text-gray-300 italic">
                Click and drag to (de)select faster
              </div>
            </div>
          </>
        )}

        {/* Question action */}
        {actionType === ActionType.Question && (
          <>
            <div className="text-sm text-gray-100">{text}</div>
            {choices && choices.length > 0 ? (
              <div className="gap-4 mt-6 flex flex-wrap">
                {choices.map((choice, i) => (
                  <div className="flex-1">
                    <Button
                      key={i}
                      onClick={() =>
                        onChoose({ cards: [], cancel: false, count: i })
                      }
                    >
                      {choice}
                    </Button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="gap-4 mt-6 flex flex-wrap">
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
            )}
          </>
        )}

        {/* Count action */}
        {actionType === ActionType.Count && (
          <>
            <div className="text-sm text-gray-100">{text}</div>
            <div className="mt-6 flex gap-4">
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
          </>
        )}

        {/* Searchable action */}
        {actionType === ActionType.Searchable &&
          choices &&
          choices.length > 0 && (
            <>
              <div className="text-sm text-gray-100">{text}</div>
              <div className="mt-6 flex gap-4">
                <select
                  className="bg-gray-800 text-white px-2 py-[0.4rem] rounded border border-gray-700 focus:outline-none focus:border-blue-500 text-xs"
                  id="action-searchable-selector"
                  value={selectedSearchValue}
                  onChange={(e) => setSelectedSearchValue(e.target.value)}
                >
                  <option value="-1" disabled>
                    Search and select
                  </option>
                  {choices?.map((choice, i) => (
                    <option key={i} value={`${i}`}>
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
            </>
          )}

        {/* Order card selection */}
        {actionType === ActionType.Order && (
          <>
            <div className="text-sm text-gray-100">{text}</div>
            <div
              className="grid gap-2 p-2 mt-4 bg-black/30 rounded-md w-fit"
              style={{
                gridTemplateColumns: `repeat(${gridCols}, minmax(0, 1fr))`,
              }}
            >
              {cards?.map((card, index) => {
                const orderNumber =
                  [...selectedCardIds].indexOf(card.virtualId) + 1;
                const isSelected = selectedCardIds.has(card.virtualId);

                return (
                  <div
                    key={index}
                    className="w-full relative"
                    data-card-id={card.virtualId}
                    onMouseEnter={() => handleCardHover(card.virtualId)}
                    onMouseDown={(e) => handleCardMouseDown(card.virtualId, e)}
                    onTouchStart={(e) => handleCardMouseDown(card.virtualId, e)}
                  >
                    <img
                      onContextMenu={(e) => {
                        e.preventDefault();
                        if (onCardRightClick && card.uid) {
                          onCardRightClick(card.uid, card.name);
                        }
                      }}
                      onDragStart={(e) => e.preventDefault()}
                      draggable={false}
                      className={`rounded-md ${
                        isSelected ? "ring-1 ring-blue-100" : ""
                      }`}
                      src={`https://scans.shobu.io/${card.uid}.jpg`}
                      alt={card.name}
                      style={{ borderRadius: "5%" }}
                    />
                    {isSelected && (
                      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                        <div className="text-white text-6xl font-bold drop-shadow-[0_0_8px_rgba(0,0,0,0.9)]">
                          {orderNumber}
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
            <div className="flex items-center gap-4 mt-4">
              <Button
                onClick={() =>
                  onChoose({
                    cards: (cards || [])
                      .map((card) => card.virtualId)
                      .filter((id) => selectedCardIds.has(id)),
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
              <div className="flex-1 text-right text-xs text-gray-300 italic">
                Click and drag to (de)select faster
              </div>
            </div>
          </>
        )}

        {error && (
          <div className="mt-4 flex items-center gap-2 text-sm text-red-500">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              className="w-5 h-5 flex-shrink-0"
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
    </Popup>
  );
}
