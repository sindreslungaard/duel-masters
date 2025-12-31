import { useEffect, useState } from "react";
import { useDuel } from "./useDuel";
import { MatchState } from "./types";

interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
}

export function Duel({ duelId, duelToken, hostUrl }: DuelProps) {
  const { connected, error, send, sendJoinMatch, state } = useDuel({
    hostUrl,
    duelId,
    duelToken,
  });

  return (
    <div className="w-full h-screen bg-gray-900 text-white p-4">
      <div className="flex flex-col">
        <div className="flex-1">1</div>
        <div className="flex-1">1</div>
        <div className="flex-1">1</div>
        <div className="flex-1">1</div>
        <div className="flex-1">1</div>
      </div>
    </div>
  );
}
