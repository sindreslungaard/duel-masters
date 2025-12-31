import { useEffect, useState } from "react";
import { useDuel } from "./useDuel";
import {
  cardHasFlag,
  CardState,
  MatchState,
  ShieldState,
  TAPPED_FLAG,
} from "./types";
import { Card } from "./Card";

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

  useEffect(() => {
    if (connected) {
      sendJoinMatch();
    }
  }, [connected]);

  if (!state) {
    return <div>Waiting for both players to join...</div>;
  }

  return (
    <div className="w-full h-screen bg-gray-900 text-white p-4">
      <div className="flex flex-col h-full">
        <div className="h-[10%] flex gap-5">
          {state.opponent.manazone.map(CreateCard)}
        </div>
        <div className="h-[10%] flex gap-5">
          {state.opponent.shieldzone.map(CreateCard)}
        </div>
        <div className="flex h-[20%] gap-5">
          {state.opponent.playzone.map(CreateCard)}
        </div>
        <div className="flex h-[20%] gap-5">
          {state.me.playzone.map(CreateCard)}
        </div>
        <div className="flex h-[10%] gap-5">
          {state.me.shieldzone.map(CreateCard)}
        </div>
        <div className="flex h-[10%] gap-5">
          {state.me.manazone.map(CreateCard)}
        </div>
        <div className="flex h-[20%] gap-5">
          {state.me.hand.map(CreateCard)}
        </div>
      </div>
    </div>
  );
}

function CreateCard(card: CardState | ShieldState, index: number) {
  const name = "name" in card && card.name ? card.name : undefined;

  return (
    <Card
      virtualId={card.virtualId}
      name={name}
      imageId={card.uid}
      key={index}
      rotated={cardHasFlag(card.flags, TAPPED_FLAG)}
    ></Card>
  );
}
