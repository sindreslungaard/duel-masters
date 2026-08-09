import {
  memo,
  startTransition,
  type ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";
import type {
  ChatMessage,
  DuelChatUser,
  DuelChatUserResolver,
  DuelChatUserTriggerRenderer,
} from "./types";

export type ReceivedChatMessage = ChatMessage & {
  receivedAt: number;
};

interface ChatProps {
  messages: ReceivedChatMessage[];
  onSendMessage: (message: string) => void;
  resolveUser?: DuelChatUserResolver;
  renderUserTrigger?: DuelChatUserTriggerRenderer;
}

const MESSAGE_GROUP_WINDOW_MS = 5 * 60 * 1000;

function normalizeUsername(username: string): string {
  return username.trim().toLowerCase();
}

function isFromServer(data: Pick<ChatMessage, "sender">): boolean {
  const sender = normalizeUsername(data.sender);

  return sender === "server" || sender === "server_1" || sender === "server_2";
}

function getInitials(username: string): string {
  const parts = username.trim().split(" ").filter(Boolean);

  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();

  return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase();
}

function formatRelativeTime(timestamp: number): string {
  const diffMs = Math.max(0, Date.now() - timestamp);
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return "just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return new Date(timestamp).toLocaleDateString();
}

function ChatAvatar({ user }: { user: DuelChatUser }) {
  const avatarUrl = user.avatarUrl ?? null;
  const [imageState, setImageState] = useState<{
    url: string | null;
    loaded: boolean;
  }>({ url: null, loaded: false });
  const imageLoaded = imageState.url === avatarUrl && imageState.loaded;

  return (
    <span className="relative flex w-10 h-10 flex-shrink-0 overflow-hidden rounded-full bg-muted">
      <span className="flex w-full h-full items-center justify-center text-xs">
        {getInitials(user.username)}
      </span>
      {avatarUrl && (
        <img
          src={avatarUrl}
          alt={user.username}
          className={`absolute inset-0 aspect-square w-full h-full ${
            imageLoaded ? "opacity-100" : "opacity-0"
          }`}
          onLoad={() => setImageState({ url: avatarUrl, loaded: true })}
          onError={() => setImageState({ url: avatarUrl, loaded: false })}
        />
      )}
    </span>
  );
}

function ChatComponent({
  messages,
  onSendMessage,
  resolveUser,
  renderUserTrigger,
}: ChatProps) {
  const [inputValue, setInputValue] = useState("");
  const [cachedUsers, setCachedUsers] = useState(
    () => new Map<string, DuelChatUser>(),
  );
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const mountedRef = useRef(false);
  const requestedUsernamesRef = useRef(new Set<string>());

  useEffect(() => {
    mountedRef.current = true;

    return () => {
      mountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    if (!resolveUser) return;

    const latestSenderByUsername = new Map<string, string>();

    messages.forEach((message) => {
      if (isFromServer(message)) return;

      latestSenderByUsername.set(
        normalizeUsername(message.sender),
        message.sender,
      );
    });

    for (const [username, sender] of latestSenderByUsername) {
      if (
        cachedUsers.has(username) ||
        requestedUsernamesRef.current.has(username)
      ) {
        continue;
      }

      requestedUsernamesRef.current.add(username);

      // Fire and forget after rendering the initials fallback. Profile work
      // must never be awaited by WebSocket handling or gameplay actions.
      void Promise.resolve()
        .then(() => resolveUser(sender))
        .then((resolvedUser) => {
          if (!resolvedUser || !mountedRef.current) return;

          startTransition(() => {
            setCachedUsers((current) => {
              const existing = current.get(username);

              // Another lookup may have populated this user while awaiting.
              if (existing) return current;

              const next = new Map(current);
              next.set(username, resolvedUser);
              return next;
            });
          });
        })
        .catch((error) => {
          console.error(
            `Failed to resolve chat user ${sender}:`,
            error,
          );
        });
    }
  }, [cachedUsers, messages, resolveUser]);

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

    return null;
  };

  const renderTrigger = (
    user: DuelChatUser,
    className: string,
    children: ReactNode,
  ) => {
    if (renderUserTrigger) {
      return renderUserTrigger({ user, className, children });
    }

    return <span className={className}>{children}</span>;
  };

  return (
    <div className="flex flex-col h-full">
      {/* Messages area */}
      <div className="flex-1 overflow-y-auto p-4 space-y-1 min-h-0 custom-scrollbar">
        {messages.map((message, index) => {
          const fromServer = isFromServer(message);

          if (!fromServer) {
            const previousMessage = index > 0 ? messages[index - 1] : null;
            const shouldGroup =
              previousMessage !== null &&
              !isFromServer(previousMessage) &&
              previousMessage.sender === message.sender &&
              message.receivedAt - previousMessage.receivedAt <
                MESSAGE_GROUP_WINDOW_MS;
            const user = cachedUsers.get(normalizeUsername(message.sender)) ?? {
              username: message.sender,
            };

            if (shouldGroup) {
              return (
                <div key={index} className="flex gap-3 -mx-2 px-2">
                  <div className="w-10 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm break-words leading-relaxed">
                      {message.message}
                    </p>
                  </div>
                </div>
              );
            }

            return (
              <div
                key={index}
                className="mt-4 flex items-start gap-3 -mx-2 px-2 first:mt-0"
              >
                {renderTrigger(
                  user,
                  "shrink-0 rounded-full",
                  <ChatAvatar user={user} />,
                )}
                <div className="flex-1 min-w-0">
                  <div className="flex items-baseline gap-2">
                    {renderTrigger(
                      user,
                      "text-sm font-semibold hover:underline",
                      user.username,
                    )}
                    <span className="text-[10px] text-muted-foreground">
                      {formatRelativeTime(message.receivedAt)}
                    </span>
                  </div>
                  <p className="text-sm break-words mt-0.5">
                    {message.message}
                  </p>
                </div>
              </div>
            );
          }

          const style = getMessageStyle(message);

          if (!style) return null;

          return (
            <div key={index} className={style.container}>
              <div className={style.bubble}>
                <div className="text-sm">
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

export const Chat = memo(ChatComponent);
