import { useEffect, useRef, useState } from "react";
import { Popup } from "./Popup";
import { ActionType, CardState } from "./types";
import { Button } from "./Button";

export interface ActionProps {
  title: string;
  visible: boolean;
  error?: string;
  onChoose: (data: { cards: string[]; cancel: false; count?: number }) => void;
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
  choices: string[];
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
  maxSelections,
}: ActionProps) {
  const [selectedCardIds, setSelectedCardIds] = useState(new Set<string>());
  const [isBrushing, setIsBrushing] = useState(false);
  const [brushedCards, setBrushedCards] = useState(new Set<string>());
  const mouseDownHandledRef = useRef(false);

  const selectCard = (cardId: string) => {
    if (selectedCardIds.has(cardId)) {
      setSelectedCardIds((prev) => {
        const next = new Set(prev);

        next.delete(cardId);
        return next;
      });
    } else {
      // Don't allow selecting more than maxSelections
      if (selectedCardIds.size >= maxSelections) {
        return;
      }
      setSelectedCardIds((prev) => new Set(prev).add(cardId));
    }
  };

  const handleBrushStart = () => {
    setIsBrushing(true);
    setBrushedCards(new Set());
  };

  const handleBrushEnd = () => {
    setIsBrushing(false);
    setBrushedCards(new Set());
  };

  const handleCardHover = (cardId: string) => {
    if (!isBrushing) return;
    
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

  const handleCardMouseDown = (cardId: string) => {
    // Start brushing and immediately handle this card
    mouseDownHandledRef.current = true;
    setIsBrushing(true);
    setBrushedCards(new Set([cardId]));
    
    // Toggle the card
    if (selectedCardIds.has(cardId)) {
      setSelectedCardIds((prev) => {
        const next = new Set(prev);
        next.delete(cardId);
        return next;
      });
    } else {
      if (selectedCardIds.size < maxSelections) {
        setSelectedCardIds((prev) => new Set(prev).add(cardId));
      }
    }
  };

  const handleCardClick = (cardId: string, e: React.MouseEvent) => {
    // Prevent double-toggle if mousedown already handled it
    if (mouseDownHandledRef.current) {
      mouseDownHandledRef.current = false;
      e.preventDefault();
      return;
    }
    selectCard(cardId);
  };

  useEffect(() => {
    if (isBrushing) {
      const handleMouseUp = () => handleBrushEnd();
      window.addEventListener('mouseup', handleMouseUp);
      return () => window.removeEventListener('mouseup', handleMouseUp);
    }
  }, [isBrushing]);

  const cardCount = cards?.length || 0;
  const gridCols = Math.max(3, Math.min(cardCount, 6));

  return (
    <>
      <Popup
        title={title}
        visible={visible}
        showCloseButton={cancellable}
        zIndex={1000}
        closeOnOutsideClick={false}
        onClose={onClose}
      >
        <div 
          className="px-6 py-6 pt-4 select-none" 
          onMouseDown={handleBrushStart}
        >
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
                onMouseEnter={() => handleCardHover(card.virtualId)}
                onMouseDown={() => handleCardMouseDown(card.virtualId)}
              >
                <img
                  onClick={(e) => handleCardClick(card.virtualId, e)}
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
          <div className="flex gap-4 mt-4">
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
          </div>
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
    </>
  );
}
