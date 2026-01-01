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
      <style>{`
        .group:hover .card-glow-animate {
          box-shadow: 0 -20px 50px -5px rgba(255, 255, 255, 0.21),
                      0 -12px 30px -3px rgba(255, 255, 255, 0.15),
                      0 -5px 15px rgba(255, 255, 255, 0.12),
                      0 0 10px rgba(255, 255, 255, 0.045);
        }
      `}</style>
      <div className="group relative pt-10 -mt-10 flex-shrink-0">
        {/* Icons container */}
        {interactable && (
          <div className="absolute top-2 left-1/2 -translate-x-1/2 flex gap-2 z-10 opacity-0 group-hover:opacity-100 transition-opacity">
            {canAddToBattlezone && (
              <div className="group/icon relative">
                <div
                  onClick={() => onAddToBattlezone?.(virtualId!)}
                  className="w-6 h-6 bg-orange-700 rounded-full flex items-center justify-center cursor-pointer hover:bg-orange-800 transition-colors"
                >
                  <span className="text-white text-xs font-bold">
                    <svg
                      className="w-full p-[4px]"
                      viewBox="0 0 16 16"
                      fill="#ffffff"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M1 7.50002C1 3.35788 4.35786 1.52588e-05 8.5 1.52588e-05V5L9.04293 5.54293L11.293 3.29289L12.7072 4.70711L10.4571 6.95714L11 7.50002H16C16 11.6422 12.6421 15 8.5 15V10L7.95714 9.45714L1.45718 15.9571L0.0429688 14.5429L6.54292 8.04292L6.00002 7.50002H1Z"
                        fill="#ffffff"
                      />
                    </svg>
                  </span>
                </div>
                <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-1 bg-gray-900 text-white text-xs rounded whitespace-nowrap opacity-0 group-hover/icon:opacity-100 transition-opacity pointer-events-none">
                  Add to battlezone
                </div>
              </div>
            )}
            {canAddToManazone && (
              <div className="group/icon relative">
                <div
                  className="w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center cursor-pointer hover:bg-blue-600 transition-colors"
                  onClick={() => onAddToManazone?.(virtualId!)}
                >
                  <span className="text-white text-xs font-bold">
                    <svg
                      fill="#ffffff"
                      className="w-full p-[4px]"
                      viewBox="0 0 8 8"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M3 0l-.34.34c-.11.11-2.66 2.69-2.66 4.88 0 1.65 1.35 3 3 3s3-1.35 3-3c0-2.18-2.55-4.77-2.66-4.88l-.34-.34zm-1.5 4.72c.28 0 .5.22.5.5 0 .55.45 1 1 1 .28 0 .5.22.5.5s-.22.5-.5.5c-1.1 0-2-.9-2-2 0-.28.22-.5.5-.5z"
                        transform="translate(1)"
                      />
                    </svg>
                  </span>
                </div>
                <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-1 bg-gray-900 text-white text-xs rounded whitespace-nowrap opacity-0 group-hover/icon:opacity-100 transition-opacity pointer-events-none">
                  Add to manazone
                </div>
              </div>
            )}
            {hasTapAbility && (
              <div className="group/icon relative">
                <div
                  onClick={() => onTapAbility?.(virtualId!)}
                  className="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center cursor-pointer hover:bg-green-600 transition-colors"
                >
                  <span className="text-white text-xs font-bold">↻</span>
                </div>
                <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-1 bg-gray-900 text-white text-xs rounded whitespace-nowrap opacity-0 group-hover/icon:opacity-100 transition-opacity pointer-events-none">
                  Tap ability
                </div>
              </div>
            )}
          </div>
        )}

        {/* Card image */}
        <img
          src={`https://scans.shobu.io/${imageId || "backside"}.jpg`}
          alt={name || "Backside card"}
          className={`h-full flex-shrink-0 rounded-md transition-all duration-300 ${
            interactable ? "card-glow-animate cursor-grab" : ""
          } ${rotated ? "rotate-90 mx-8" : ""} ${flipped ? "rotate-180" : ""}`}
        />
      </div>
    </>
  );
}
