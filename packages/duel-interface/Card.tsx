import { useEffect, useRef, useState } from "react";

export interface CardProps {
  virtualId?: string;
  name?: string;
  imageId?: string;
  rotated?: boolean;
  flipped?: boolean;
  /** Small index drawn in the card's top right corner, used to number shields. */
  number?: number;
  interactable?: boolean;
  canAddToBattlezone?: boolean;
  canAddToManazone?: boolean;
  hasTapAbility?: boolean;
  selected?: boolean;
  onAddToBattlezone?: (virtualId: string) => void;
  onAddToManazone?: (virtualId: string) => void;
  onTapAbility?: (virtualId: string) => void;
  isDragging?: boolean;
  draggable?: boolean;
  onDragStart?: (e: React.PointerEvent) => void;
  onRightClick?: () => void;
}

export function Card({
  virtualId,
  name,
  imageId,
  rotated = false,
  flipped = false,
  number,
  interactable = false,
  selected = false,
  canAddToBattlezone = true,
  canAddToManazone = true,
  hasTapAbility = false,
  onAddToBattlezone,
  onAddToManazone,
  onTapAbility,
  isDragging = false,
  draggable = false,
  onDragStart,
  onRightClick,
}: CardProps) {
  const imgRef = useRef<HTMLImageElement>(null);
  const [horizontalMargin, setHorizontalMargin] = useState("2rem");

  useEffect(() => {
    const updateMargin = () => {
      if (imgRef.current && rotated) {
        const height = imgRef.current.offsetHeight;
        // Calculate margin proportional to card height
        // Use smaller margins for smaller cards, larger for bigger cards
        const margin = Math.max(12, Math.min(32, height * 0.15));
        setHorizontalMargin(`${margin}px`);
      }
    };

    updateMargin();
    window.addEventListener("resize", updateMargin);
    return () => window.removeEventListener("resize", updateMargin);
  }, [rotated]);

  // Pointer events give one stream for mouse, touch and pen. Binding mouse and
  // touch separately made a tap fire both paths, because the browser follows a
  // touch with compatibility mouse events, and the second path immediately
  // undid the selection the first one made.
  const handlePointerDown = (e: React.PointerEvent) => {
    if (!draggable || !onDragStart) {
      return;
    }

    e.preventDefault();
    onDragStart(e);
  };

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    if (onRightClick) {
      onRightClick();
    }
  };

  return (
    <>
      <div className="group relative flex-shrink-0 h-full">
        {/* Card image */}
        <img
          ref={imgRef}
          src={`https://scans.shobu.io/${imageId || "backside"}.jpg`}
          alt={name || "Backside card"}
          draggable={false}
          className={`block h-full w-auto flex-shrink-0 ${
            interactable && !isDragging ? "cursor-grab" : ""
          } ${
            rotated && flipped
              ? "-rotate-90"
              : rotated
                ? "rotate-90"
                : flipped
                  ? "rotate-180"
                  : ""
          } ${isDragging ? "opacity-0" : ""} ${
            selected ? "ring-1 ring-blue-100" : ""
          }`}
          onPointerDown={handlePointerDown}
          onContextMenu={handleContextMenu}
          style={{
            // Draggable cards keep the gesture entirely, which is what makes
            // dragging them reliable. Cards that cannot be dragged have nothing
            // to protect, so they hand the gesture back and their row becomes
            // scrollable by touching the cards themselves rather than only the
            // gaps between them.
            touchAction: draggable ? "none" : "auto",
            borderRadius: "5%",
            marginLeft: rotated ? horizontalMargin : undefined,
            marginRight: rotated ? horizontalMargin : undefined,
          }}
        />

        {/* Kept upright and in the same screen corner even when the card is
            flipped, so a shield's number stays readable from either side. */}
        {number !== undefined && (
          <span className="pointer-events-none absolute top-[2%] right-[10%] text-[clamp(0.55rem,1.1vh,0.8rem)] leading-none text-white drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]">
            {number}
          </span>
        )}
      </div>
    </>
  );
}
