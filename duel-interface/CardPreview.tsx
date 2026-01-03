import { useEffect } from "react";

interface CardPreviewProps {
  visible: boolean;
  imageUrl: string;
  onClose: () => void;
}

export function CardPreview({ visible, imageUrl, onClose }: CardPreviewProps) {
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

  return (
    <div
      className="fixed inset-0 bg-black/60 flex items-center justify-center p-4"
      style={{ zIndex: 1000 }}
      onClick={handleOverlayClick}
    >
      <div className="flex flex-col items-center gap-4">
        <img
          src={imageUrl}
          alt="Card preview"
          className="rounded-lg shadow-2xl"
          style={{
            maxHeight: "80vh",
            maxWidth: "90vw",
            height: "auto",
            width: "auto",
          }}
        />
        <button
          onClick={onClose}
          className="px-6 py-2 bg-gray-800 hover:bg-gray-700 text-white rounded-md transition-colors"
        >
          Close
        </button>
      </div>
    </div>
  );
}
