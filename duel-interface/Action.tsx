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
}: ActionProps) {
  const [selectedCardIds, setSelectedCardIds] = useState(new Set<string>());

  const selectCard = (cardId: string) => {
    if (selectedCardIds.has(cardId)) {
      setSelectedCardIds((prev) => {
        const next = new Set(prev);

        next.delete(cardId);
        return next;
      });
    } else {
      setSelectedCardIds((prev) => new Set(prev).add(cardId));
    }
  };

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
        <div className="px-6 py-6 pt-4">
          <div className="text-sm text-gray-100">{text}</div>
          <div className="flex gap-2 p-2 mt-4 bg-black/30 rounded-md">
            {cards?.map((card, index) => (
              <div key={index} className="w-30">
                <img
                  onClick={() => selectCard(card.virtualId)}
                  onContextMenu={(e) => {
                    e.preventDefault();
                    if (onCardRightClick && card.uid) {
                      onCardRightClick(card.uid, card.name);
                    }
                  }}
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
