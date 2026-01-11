import { useEffect, useRef, useState } from "react";
import {
  ActionMessage,
  ActionType,
  ActionWarningMessage,
  ChatMessage,
  MatchState,
  MatchStateMessage,
  ShowCardsMessage,
  WaitMessage,
  WarningMessage,
} from "./types";

interface UseDuelOptions {
  duelId: string;
  duelToken: string;
  hostUrl: string;
  onActionMessage?: (message: ActionMessage) => void;
  onActionError?: (message: ActionWarningMessage) => void;
  onActionClose?: () => void;
  onChat?: (message: ChatMessage) => void;
  onWarning?: (message: WarningMessage) => void;
  onWait?: (message: WaitMessage) => void;
  onEndWait?: () => void;
  onShowCards?: (message: ShowCardsMessage) => void;
  onShowCardsNonDismissable?: (message: ShowCardsMessage) => void;
}

export function useDuel({
  duelId,
  duelToken,
  hostUrl,
  onActionMessage,
  onActionError,
  onActionClose,
  onChat,
  onWarning,
  onWait,
  onEndWait,
}: UseDuelOptions) {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [state, setState] = useState<MatchState | null>(null);
  const [opponentDisconnected, setOpponentDisconnected] = useState(false);

  useEffect(() => {
    const wsUrl = `${hostUrl}/ws/${duelId}?duelToken=${duelToken}`;
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
      setConnected(true);
      setError(null);
    };

    ws.onerror = (event) => {
      console.error("WebSocket error:", event);
      setError("WebSocket connection error");
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
      setConnected(false);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        switch (data.header) {
          case "state_update":
            setState(data.state);
            break;
          case "action":
            onActionMessage?.(data);
            break;
          case "action_error":
            onActionError?.(data);
            break;
          case "close_action":
            onActionClose?.();
            break;
          case "chat":
            onChat?.(data);
            break;
          case "warn":
            onWarning?.(data);
            break;
          case "wait":
            onWait?.(data);
            break;
          case "end_wait":
            onEndWait?.();
            break;
          case "show_cards":
            onActionMessage?.({
              header: data.header,
              actionType: ActionType.ShowCards,
              text: data.message,
              showCards: {
                cards: data.cards,
                dismissable: true,
              },
            });
            break;
          case "show_cards_non_dismissible":
            onActionMessage?.({
              header: data.header,
              actionType: ActionType.ShowCards,
              text: data.message,
              showCards: {
                cards: data.cards,
                dismissable: false,
              },
            });
            break;
          case "opponent_disconnected":
            setOpponentDisconnected(true);
            break;

          case "opponent_reconnected":
            setOpponentDisconnected(false);
            break;

          default:
            console.log("Unknown message type:", data.header);
        }
      } catch (err) {
        console.error("Error parsing message:", err);
      }
    };

    // Cleanup on unmount
    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, [duelId, duelToken, hostUrl]);

  const send = (data: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data));
    } else {
      console.error("WebSocket is not connected");
    }
  };

  const sendJoinMatch = () => {
    send({ header: "join_match" });
  };

  const sendEndTurn = () => {
    send({ header: "end_turn" });
  };

  const sendAddToBattlezone = (virtualId: string) => {
    send({ header: "add_to_playzone", virtualId });
  };

  const sendAddToManazone = (virtualId: string) => {
    send({ header: "add_to_manazone", virtualId });
  };

  const sendAttackPlayer = (virtualId: string) => {
    send({ header: "attack_player", virtualId });
  };

  const sendAttackCreature = (virtualId: string) => {
    send({ header: "attack_creature", virtualId });
  };

  const sendTapAbility = (virtualId: string) => {
    send({ header: "tap_ability", virtualId });
  };

  const sendAction = (data: {
    cards: string[];
    cancel: boolean;
    count?: number;
  }) => {
    send({ header: "action", ...data });
  };

  const sendChat = (message: string) => {
    send({ header: "chat", message });
  };

  return {
    connected,
    error,
    state,
    opponentDisconnected,
    send,
    sendJoinMatch,
    sendEndTurn,
    sendAddToBattlezone,
    sendAddToManazone,
    sendAttackPlayer,
    sendAttackCreature,
    sendTapAbility,
    sendAction,
    sendChat,
  };
}
