import { useEffect, useRef, useState } from "react";

interface PopupProps {
  visible: boolean;
  onClose?: () => void;
  showCloseButton?: boolean;
  closeOnOutsideClick?: boolean;
  maxWidth?: string;
  maxHeight?: string;
  contentClassName?: string;
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
  maxHeight = "85vh",
  contentClassName = "overflow-y-auto flex-1",
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

  /** Keeps the popup reachable: it may hang off an edge, but never so far that
   * there is nothing left to grab or read. Without this a popup dragged off
   * screen, or left offset when the window shrank, could not be recovered and
   * the prompt behind it could not be answered. */
  const clampToViewport = (next: { x: number; y: number }) => {
    const popup = popupRef.current;
    if (!popup) return next;

    // offsetLeft/offsetTop ignore the transform, so they describe where the
    // centred popup sits before any drag offset. Measuring from the transformed
    // rect instead would need the current offset, which is not reliably in scope
    // while a drag is in flight.
    const restLeft = popup.offsetLeft;
    const restTop = popup.offsetTop;
    const width = popup.offsetWidth;

    // Always leave this much of the popup on screen in each direction.
    const margin = 64;

    return {
      x: Math.min(
        Math.max(next.x, margin - restLeft - width),
        window.innerWidth - margin - restLeft,
      ),
      // The header must stay reachable, so the popup may never go above the top
      // edge, and never so far down that the header scrolls out of sight.
      y: Math.min(
        Math.max(next.y, -restTop),
        window.innerHeight - margin - restTop,
      ),
    };
  };

  const handlePointerDown = (e: React.PointerEvent) => {
    // Don't start dragging if clicking on the close button
    if ((e.target as HTMLElement).closest("button")) {
      return;
    }

    setIsDragging(true);
    setDragStart({
      x: e.clientX - position.x,
      y: e.clientY - position.y,
    });
  };

  useEffect(() => {
    const handlePointerMove = (e: PointerEvent) => {
      if (!isDragging) return;

      setPosition(
        clampToViewport({
          x: e.clientX - dragStart.x,
          y: e.clientY - dragStart.y,
        }),
      );
    };

    const handleEnd = () => {
      setIsDragging(false);
    };

    if (isDragging) {
      document.addEventListener("pointermove", handlePointerMove);
      document.addEventListener("pointerup", handleEnd);
      document.addEventListener("pointercancel", handleEnd);
    }

    return () => {
      document.removeEventListener("pointermove", handlePointerMove);
      document.removeEventListener("pointerup", handleEnd);
      document.removeEventListener("pointercancel", handleEnd);
    };
  }, [isDragging, dragStart]);

  // A resize or an orientation change can leave an offset popup off screen.
  useEffect(() => {
    if (!visible) return;

    const handleResize = () => setPosition((current) => clampToViewport(current));

    window.addEventListener("resize", handleResize);
    window.addEventListener("orientationchange", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("orientationchange", handleResize);
    };
  }, [visible]);

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
      className="fixed inset-0 bg-black/60 flex items-center justify-center p-2 sm:p-4"
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
            className="flex-shrink-0 px-4 py-3 sm:px-6 sm:py-4 border-b border-gray-800 cursor-move select-none"
            style={{ touchAction: "none" }}
            onPointerDown={handlePointerDown}
          >
            <h2 className="pr-10 text-base font-semibold text-white sm:text-xl">
              {title}
            </h2>
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
        <div className={contentClassName}>{children}</div>
      </div>
    </div>
  );
}
