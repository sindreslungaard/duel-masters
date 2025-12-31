interface DuelProps {
  duelId: string;
  duelToken: string;
}

export function Duel({ duelId, duelToken }: DuelProps) {
  return (
    <div>
      Duel Interface for duel {duelId} with token {duelToken}
    </div>
  );
}
