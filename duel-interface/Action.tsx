import { useEffect, useRef } from "react";
import { Popup } from "./Popup";
import { ActionType, CardState } from "./types";
import { Button } from "./Button";

export interface SelectableCard extends CardState {
  selected?: boolean;
}

export interface ActionProps {
  title: string;
  visible: boolean;
  error?: string;
  onChoose: (data: { cards: string[]; cancel: false; count?: number }) => void;
  onClose: () => void;

  // Action details
  actionType: ActionType;
  cards: SelectableCard[];
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
}: ActionProps) {
  return (
    <>
      <Popup
        title={title}
        visible={visible}
        showCloseButton={cancellable}
        zIndex={1000}
        closeOnOutsideClick={false}
      >
        <div className="px-6 py-6 pt-4">
          <div className="text-sm text-gray-100">{text}</div>
          <div className="flex gap-2 p-2 mt-4 bg-black/50 rounded-md">
            {cards.map((card, index) => (
              <div key={index} className="w-30">
                <img
                  className="rounded-md"
                  src={`https://scans.shobu.io/${card.uid}.jpg`}
                  alt={card.name}
                />
              </div>
            ))}
          </div>
          <div className="flex gap-4 mt-4">
            <Button
              onClick={() =>
                onChoose({
                  cards: cards
                    .filter((card) => card.selected)
                    .map((card) => card.virtualId),
                  cancel: false,
                  count: 0,
                })
              }
            >
              Choose
            </Button>
            <Button variant="gray" onClick={onClose}>
              Close
            </Button>
          </div>
          {error && <div className="mt-4 text-sm text-red-500">{error}</div>}
        </div>
      </Popup>
    </>
  );
}
