import ErrorDialog from "../components/dialogs/ErrorDialog";
import WarningDialog from "../components/dialogs/WarningDialog";

/**
 * This BaseMixin contains methods that are used in multiple components
 * throughout the project.
 */
export default {
  data() {
    return {
      isLoading: false,
    };
  },
  methods: {
    /**
     * Shows an error dialog with an message and a reload
     * button.
     *
     * @param {String} message
     */
    showError(message) {
      console.log(message);
      this.$vfm.show({
        component: ErrorDialog,
        bind: {
          message,
        },
        on: {
          "closed": () => window.location.reload(),
        },
      });
    },
    /**
     * Shows a warning dialog with an message and a close
     * button.
     *
     * @param {String} message
     */
    showWarning(message) {
      this.$vfm.show({
        component: WarningDialog,
        bind: {
          message,
        },
      });
    },
    /**
     * Connects to a WebSocket with the given endpoint and setups error handlers.
     *
     * @param {String} endpoint
     */
    connectToSocket(endpoint) {
      this.isLoading = true;
      const socket = new WebSocket(this.$config.SOCKET_ENDPOINT + endpoint);

      socket.onopen = () => {
        socket.send(localStorage.getItem("token"));
        this.isLoading = false;
      };

      socket.onerror = () => {
        this.showError("Lost connection to the server");
      };

      return socket;
    },
    /**
     * Sends an object to the given WebSocket.
     *
     * @param {WebSocket} socket
     * @param {Object} data
     */
    sendMessage(socket, data) {
      socket.send(JSON.stringify(data));
    },
  },
};
