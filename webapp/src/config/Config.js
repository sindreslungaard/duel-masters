export default {
  /**
   * The URL where the REST API is served.
   */
  API_ENDPOINT: `${window.location.origin}/api`,
  /**
   * The URL for WebSockets.
   */
  SOCKET_ENDPOINT: `wss://${window.location.host}`,
  /**
   * The time duration in which multiple messages from one user
   * will count as one message in seconds.
   */
  CHAT_MESSAGE_APPEND_TIMEOUT: 15,
  /**
   * The prefix to be added in front of every localStorage item.
   */
  STORAGE_PREFIX: "shobu",
  /**
   * The key for saving the name of the last created duel.
   */
  STORAGE_KEY_DUEL_NAME: "duel_name",
};
