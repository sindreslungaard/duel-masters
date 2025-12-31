import { useEffect, useRef, useState } from "react";
import { MatchState, MatchStateMessage } from "./types";

interface UseDuelOptions {
  duelId: string;
  duelToken: string;
  hostUrl: string;
}

export function useDuel({ duelId, duelToken, hostUrl }: UseDuelOptions) {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [state, setState] = useState<MatchState | null>(null);

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
      console.log("WebSocket message:", event.data);
      try {
        const data = JSON.parse(event.data);

        switch (data.header) {
          case "state_update":
            setState(data.state);
            break;
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

  const sendAddToBattlezone = (virtualId: string) => {
    send({ header: "add_to_playzone", virtualId });
  };

  const sendAddToManazone = (virtualId: string) => {
    send({ header: "add_to_manazone", virtualId });
  };

  const sendTapAbility = (virtualId: string) => {
    send({ header: "tap_ability", virtualId });
  };

  return {
    connected,
    error,
    send,
    sendJoinMatch,
    sendAddToBattlezone,
    sendAddToManazone,
    sendTapAbility,
    state,
  };
}
