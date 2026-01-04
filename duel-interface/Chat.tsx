import { useState, useRef, useEffect } from "react";
import { ChatMessage } from "./types";

interface ChatProps {
  messages: ChatMessage[];
  onSendMessage: (message: string) => void;
}

export function Chat({ messages, onSendMessage }: ChatProps) {
  const [inputValue, setInputValue] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const isFromServer = (data: ChatMessage) => {
    let sender = data.sender.toLowerCase();
    if (sender === "server" || sender === "server_1" || sender === "server_2")
      return true;

    return false;
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = () => {
    if (inputValue.trim()) {
      onSendMessage(inputValue);
      setInputValue("");
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleSend();
    }
  };

  const getMessageStyle = (message: ChatMessage) => {
    const sender = message.sender.toLowerCase();

    if (sender === "server") {
      return {
        container: "flex justify-center",
        bubble: "text-gray-400 text-center px-3 py-2",
      };
    }

    if (sender === "server_1") {
      return {
        container: "flex justify-start",
        bubble: "bg-orange-200/30 max-w-[80%] rounded-sm px-2 py-1",
      };
    }

    if (sender === "server_2") {
      return {
        container: "flex justify-end",
        bubble: "bg-gray-400/30 max-w-[80%] rounded-sm px-2 py-1",
      };
    }

    // Player message
    return {
      container: "flex justify-start",
      bubble: "bg-black/30 w-full rounded-md px-2 py-1",
    };
  };

  return (
    <div className="flex flex-col h-full">
      {/* Messages area */}
      <div className="flex-1 overflow-y-auto p-2 space-y-2 min-h-0 custom-scrollbar">
        {messages.map((message, index) => {
          const style = getMessageStyle(message);
          const fromServer = isFromServer(message);
          return (
            <div key={index} className={style.container}>
              <div className={style.bubble}>
                <div className="text-sm">
                  {!fromServer && (
                    <>
                      <span className="text-orange-500">{message.sender}:</span>{" "}
                    </>
                  )}
                  <span className="text-white">{message.message}</span>
                </div>
              </div>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      {/* Input area */}
      <div className="p-2">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type to chat"
          className="text-sm w-full bg-black/30 text-white placeholder-gray-400 rounded-md px-3 py-2 outline-none focus:ring-1 focus:ring-white/30"
        />
      </div>
    </div>
  );
}
