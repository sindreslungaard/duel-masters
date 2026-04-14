import { useEffect, useRef, useState } from "react";

export interface CardProps {
  virtualId?: string;
  name?: string;
  imageId?: string;
  rotated?: boolean;
  flipped?: boolean;
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
  onDragStart?: (e: React.MouseEvent | React.TouchEvent) => void;
  onRightClick?: () => void;
}

export function Card({
  virtualId,
  name,
  imageId,
  rotated = false,
  flipped = false,
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

  const handleMouseDown = (e: React.MouseEvent | React.TouchEvent) => {
    if (draggable && onDragStart) {
      e.preventDefault();
      // Prevent context menu on touch devices
      if ("touches" in e) {
        e.stopPropagation();
      }
      onDragStart(e);
    }
  };

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    if (onRightClick) {
      onRightClick();
    }
  };

  return (
    <>
      <div className="group relative pt-10 -mt-10 flex-shrink-0">
        {/* Card image */}
        <img
          ref={imgRef}
          src={`https://scans.shobu.io/${imageId || "backside"}.jpg`}
          alt={name || "Backside card"}
          draggable={false}
          className={`h-full flex-shrink-0 ${
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
          onMouseDown={handleMouseDown}
          onTouchStart={handleMouseDown}
          onContextMenu={handleContextMenu}
          style={{
            touchAction: "none",
            borderRadius: "5%",
            marginLeft: rotated ? horizontalMargin : undefined,
            marginRight: rotated ? horizontalMargin : undefined,
          }}
        />
      </div>
    </>
  );
}
