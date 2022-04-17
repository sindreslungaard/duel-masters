const config = JSON.parse(localStorage.getItem("config") ?? "null") ?? {
  host: window.location.host,
  ws_protocol: location.protocol == "https:" ? "wss://" : "ws://",
  api: "/api",
};

export default config;
