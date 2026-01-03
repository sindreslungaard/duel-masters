import { useEffect } from "react";
import { Button } from "./Button";

interface CardPreviewProps {
  visible: boolean;
  imageId: string | null;
  name: string | null;
  onClose: () => void;
}

export function CardPreview({
  visible,
  imageId,
  name,
  onClose,
}: CardPreviewProps) {
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

  const imageUrl = imageId ? `https://scans.shobu.io/${imageId}.jpg` : "";

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
          className="shadow-2xl"
          style={{
            maxHeight: "80vh",
            maxWidth: "90vw",
            height: "auto",
            width: "auto",
            borderRadius: "5%",
          }}
        />
        <Button variant="gray" onClick={onClose}>
          Close Preview
        </Button>
      </div>
    </div>
  );
}
