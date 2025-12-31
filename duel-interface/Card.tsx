export interface CardProps {
  virtualId?: string;
  name?: string;
  imageId?: string;
  rotated?: boolean;
}

export function Card({ virtualId, name, imageId, rotated = false }: CardProps) {
  return (
    <div>
      <img
        src={`https://scans.shobu.io/${imageId || "backside"}.jpg`}
        alt={name}
        className={`h-full rounded-xl ${rotated ? "rotate-90 mx-8" : ""}`}
      />
    </div>
  );
}
