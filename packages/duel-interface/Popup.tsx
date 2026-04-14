import { useEffect, useRef, useState } from "react";

interface PopupProps {
  visible: boolean;
  onClose?: () => void;
  showCloseButton?: boolean;
  closeOnOutsideClick?: boolean;
  maxWidth?: string;
  maxHeight?: string;
  zIndex?: number;
  title?: string;
  children: React.ReactNode;
}

export function Popup({
  visible,
  onClose,
  showCloseButton = true,
  closeOnOutsideClick = false,
  maxWidth = "600px",
  maxHeight = "80vh",
  zIndex = 1000,
  title,
  children,
}: PopupProps) {
  const popupRef = useRef<HTMLDivElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  // Reset position when popup becomes visible
  useEffect(() => {
    if (visible) {
      setPosition({ x: 0, y: 0 });
    }
  }, [visible]);

  const handlePointerDown = (e: React.MouseEvent | React.TouchEvent) => {
    // Don't start dragging if clicking on the close button
    if ((e.target as HTMLElement).closest("button")) {
      return;
    }

    const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
    const clientY = "touches" in e ? e.touches[0].clientY : e.clientY;

    setIsDragging(true);
    setDragStart({
      x: clientX - position.x,
      y: clientY - position.y,
    });
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDragging) return;

      setPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      });
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (!isDragging) return;

      setPosition({
        x: e.touches[0].clientX - dragStart.x,
        y: e.touches[0].clientY - dragStart.y,
      });
    };

    const handleEnd = () => {
      setIsDragging(false);
    };

    if (isDragging) {
      document.addEventListener("mousemove", handleMouseMove);
      document.addEventListener("mouseup", handleEnd);
      document.addEventListener("touchmove", handleTouchMove);
      document.addEventListener("touchend", handleEnd);
    }

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleEnd);
      document.removeEventListener("touchmove", handleTouchMove);
      document.removeEventListener("touchend", handleEnd);
    };
  }, [isDragging, dragStart]);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && visible && showCloseButton && onClose) {
        onClose();
      }
    };

    if (visible) {
      document.addEventListener("keydown", handleEscape);
      // Prevent body scroll when popup is open
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "";
    };
  }, [visible, onClose]);

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (
      closeOnOutsideClick &&
      onClose &&
      popupRef.current &&
      !popupRef.current.contains(e.target as Node)
    ) {
      onClose();
    }
  };

  if (!visible) return null;

  return (
    <div
      className="fixed inset-0 bg-black/60 flex items-center justify-center p-4"
      style={{ zIndex }}
      onClick={handleOverlayClick}
    >
      <div
        ref={popupRef}
        className="bg-gray-900 rounded-lg shadow-2xl relative w-full overflow-hidden flex flex-col"
        style={{
          maxWidth,
          maxHeight,
          transform: `translate(${position.x}px, ${position.y}px)`,
          transition: isDragging ? "none" : "transform 0.1s ease-out",
        }}
      >
        {/* Header */}
        {title && (
          <div
            className="flex-shrink-0 px-6 py-4 border-b border-gray-800 cursor-move select-none"
            onMouseDown={handlePointerDown}
            onTouchStart={handlePointerDown}
          >
            <h2 className="text-xl font-semibold text-white pr-8">{title}</h2>
          </div>
        )}

        {/* Close button */}
        {showCloseButton && onClose && (
          <button
            onClick={onClose}
            className="cursor-pointer absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white z-10"
            aria-label="Close"
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
        )}

        {/* Content */}
        <div className="overflow-y-auto flex-1">{children}</div>
      </div>
    </div>
  );
}
