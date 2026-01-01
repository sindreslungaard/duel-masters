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
  onAddToBattlezone?: (virtualId: string) => void;
  onAddToManazone?: (virtualId: string) => void;
  onTapAbility?: (virtualId: string) => void;
}

export function Card({
  virtualId,
  name,
  imageId,
  rotated = false,
  flipped = false,
  interactable = false,
  canAddToBattlezone = true,
  canAddToManazone = true,
  hasTapAbility = false,
  onAddToBattlezone,
  onAddToManazone,
  onTapAbility,
}: CardProps) {
  return (
    <>
      <div className="group relative pt-10 -mt-10 flex-shrink-0">
        {/* Card image */}
        <img
          src={`https://scans.shobu.io/${imageId || "backside"}.jpg`}
          alt={name || "Backside card"}
          className={`h-full flex-shrink-0 rounded-md transition-all duration-300 ${
            interactable ? "cursor-grab" : ""
          } ${rotated ? "rotate-90 mx-8" : ""} ${flipped ? "rotate-180" : ""}`}
        />
      </div>
    </>
  );
}
