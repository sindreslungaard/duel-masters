import { useEffect, useState } from "react";
import { useDuel } from "./useDuel";
import {
  cardHasFlag,
  CardState,
  MatchState,
  PLAYABLE_FLAG,
  ShieldState,
  TAPPED_FLAG,
} from "./types";
import { Card } from "./Card";
import { Button } from "./Button";

interface DuelProps {
  hostUrl: string;
  duelId: string;
  duelToken: string;
}

export function Duel({ duelId, duelToken, hostUrl }: DuelProps) {
  const {
    connected,
    error,
    send,
    sendJoinMatch,
    sendEndTurn,
    sendAddToBattlezone,
    sendAddToManazone,
    sendTapAbility,
    state,
  } = useDuel({
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
    <div className="w-full h-screen text-white flex bg-[url('https://i.imgur.com/mWy5Cnl.gif')] bg-cover bg-center gap-2 p-2">
      <div className="w-[300px] flex flex-col gap-2">
        <div className="flex-1 bg-black/50 rounded-md"></div>
        <Button
          onClick={sendEndTurn}
          disabled={!state.myTurn}
          disabledTooltip="It's not your turn"
        >
          End turn
        </Button>
      </div>
      <div className="flex flex-col h-full w-full">
        <div className="h-[10%] flex gap-5 pb-1">
          {state.opponent.manazone.map(CreateCard({ flipped: true }))}
        </div>
        <div className="h-[10%] flex gap-5 p-1 w-full">
          {state.opponent.shieldzone.map(CreateCard())}
        </div>
        <div className="flex h-[20%] gap-5 p-1 w-full">
          {state.opponent.playzone.map(CreateCard({ flipped: true }))}
        </div>
        <div className="flex h-[20%] gap-5 p-1 w-full">
          {state.me.playzone.map(CreateCard())}
        </div>
        <div className="flex h-[10%] gap-5 p-1 w-full">
          {state.me.shieldzone.map(CreateCard())}
        </div>
        <div className="flex h-[10%] gap-5 p-1 w-full">
          {state.me.manazone.map(CreateCard({ flipped: true }))}
        </div>
        <div className="flex h-[20%] gap-5 pt-1 w-full">
          {state.me.hand.map(
            CreateCard({
              interactable: true,
              canAddToManazone: !state.hasAddedManaThisRound,
              onAddToBattlezone: (virtualId) => {
                sendAddToBattlezone(virtualId);
              },
              onAddToManazone: (virtualId) => {
                sendAddToManazone(virtualId);
              },
              onTapAbility: (virtualId) => {
                sendTapAbility(virtualId);
              },
            })
          )}
        </div>
      </div>
    </div>
  );
}

function CreateCard(
  options: {
    interactable?: boolean;
    canAddToManazone?: boolean;
    flipped?: boolean;
    onAddToBattlezone?: (virtualId: string) => void;
    onAddToManazone?: (virtualId: string) => void;
    onTapAbility?: (virtualId: string) => void;
  } = {}
) {
  return (card: CardState | ShieldState, index: number) => {
    const name = "name" in card && card.name ? card.name : undefined;

    return (
      <Card
        virtualId={card.virtualId}
        name={name}
        imageId={card.uid}
        key={index}
        rotated={cardHasFlag(card.flags, TAPPED_FLAG)}
        interactable={options.interactable}
        canAddToBattlezone={cardHasFlag(card.flags, PLAYABLE_FLAG)}
        canAddToManazone={options.canAddToManazone}
        onAddToBattlezone={options.onAddToBattlezone}
        onAddToManazone={options.onAddToManazone}
        onTapAbility={options.onTapAbility}
        flipped={options.flipped}
      ></Card>
    );
  };
}
