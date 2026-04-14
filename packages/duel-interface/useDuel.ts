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
  const [reconnecting, setReconnecting] = useState(false);
  const reconnectTimeoutRef = useRef<number | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const isUnmountingRef = useRef(false);

  useEffect(() => {
    isUnmountingRef.current = false;

    const connect = () => {
      // Clear any existing reconnect timeout
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }

      const wsUrl = `${hostUrl}/ws/${duelId}?duelToken=${duelToken}`;
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log("WebSocket connected");
        setConnected(true);
        setError(null);
        setReconnecting(false);
        reconnectAttemptsRef.current = 0; // Reset reconnect attempts on successful connection
      };

      ws.onerror = (event) => {
        console.error("WebSocket error:", event);
        setError("WebSocket connection error");
      };

      ws.onclose = () => {
        console.log("WebSocket closed");
        setConnected(false);
        setReconnecting(true);

        // Only attempt to reconnect if we're not unmounting
        if (!isUnmountingRef.current) {
          const maxReconnectAttempts = 20;
          const baseDelay = 1000; // Start with 1 second
          const maxDelay = 5000; // Max 5 seconds

          if (reconnectAttemptsRef.current < maxReconnectAttempts) {
            // First attempt is immediate, then exponential backoff with cap
            const delay =
              reconnectAttemptsRef.current === 0
                ? 0
                : Math.min(
                    baseDelay * Math.pow(2, reconnectAttemptsRef.current - 1),
                    maxDelay
                  );

            console.log(
              `Reconnecting in ${delay}ms (attempt ${
                reconnectAttemptsRef.current + 1
              }/${maxReconnectAttempts})`
            );

            reconnectAttemptsRef.current++;
            reconnectTimeoutRef.current = setTimeout(connect, delay);
          } else {
            setReconnecting(false);
            setError("Connection lost. Max reconnection attempts reached.");
            alert(
              "Connection lost. Max reconnection attempts reached. Please refresh the page or go back to the lobby."
            );
          }
        }
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
    };

    // Start the initial connection
    connect();

    // Cleanup on unmount
    return () => {
      isUnmountingRef.current = true;

      // Clear any pending reconnect timeout
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }

      // Close the WebSocket if it's open
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
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
    reconnecting,
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
