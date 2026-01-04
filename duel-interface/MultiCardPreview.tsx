import { useEffect } from "react";
import { Button } from "./Button";

interface MultiCardPreviewCard {
  imageId: string;
  name: string;
}

interface MultiCardPreviewProps {
  visible: boolean;
  cards: MultiCardPreviewCard[];
  title: string;
  onClose: () => void;
  onCardClick?: (imageId: string, name: string) => void;
}

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

export function MultiCardPreview({
  visible,
  cards,
  title,
  onClose,
  onCardClick,
}: MultiCardPreviewProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && visible) {
        onClose();
      }
    };

    if (visible) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "";
    };
  }, [visible, onClose]);

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  if (!visible) return null;

  const handleCardClick = (card: MultiCardPreviewCard, e: React.MouseEvent) => {
    e.preventDefault();
    if (onCardClick) {
      onCardClick(card.imageId, card.name);
    }
  };

  return (
    <>
      <style>{scrollbarStyles}</style>
      <div
        className="fixed inset-0 bg-black/60 flex items-center justify-center p-4"
        style={{ zIndex: 1000 }}
        onClick={handleOverlayClick}
      >
        <div className="flex flex-col items-center gap-4 max-w-[90vw] max-h-[90vh]">
          <h2 className="text-white text-2xl font-bold">{title}</h2>
          <div
            className="overflow-y-auto p-4 bg-black/40 rounded-lg custom-scrollbar"
            style={{
              maxHeight: "calc(90vh - 120px)",
            }}
          >
            <div
              className="grid gap-3 justify-items-center"
              style={{
                gridTemplateColumns: `repeat(${Math.min(cards.length, 6)}, 150px)`,
              }}
            >
              {cards.map((card, index) => (
                <div key={index} className="flex flex-col items-center">
                  <img
                    src={`https://scans.shobu.io/${card.imageId}.jpg`}
                    alt={card.name}
                    className="shadow-lg hover:scale-105 transition-transform cursor-pointer"
                    style={{
                      width: "150px",
                      height: "auto",
                      borderRadius: "5%",
                    }}
                    onClick={(e) => handleCardClick(card, e)}
                    onContextMenu={(e) => handleCardClick(card, e)}
                  />
                </div>
              ))}
            </div>
          </div>
          <Button variant="gray" onClick={onClose}>
            Close Preview
          </Button>
        </div>
      </div>
    </>
  );
}
