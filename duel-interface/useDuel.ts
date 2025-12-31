import { useEffect, useRef, useState } from "react";

interface UseDuelOptions {
  duelId: string;
  duelToken: string;
  hostUrl: string;
}

export function useDuel({ duelId, duelToken, hostUrl }: UseDuelOptions) {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [gameState, setGameState] = useState<any>(null);
  const wsRef = useRef<WebSocket | null>(null);

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
        setMessages((prev) => [...prev, data]);

        // Update game state if the message contains it
        if (data.gameState) {
          setGameState(data.gameState);
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

  return {
    connected,
    error,
    messages,
    gameState,
    send,
    sendJoinMatch,
  };
}
