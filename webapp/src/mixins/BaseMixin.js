import ErrorDialog from "../components/dialogs/ErrorDialog";
import WarningDialog from "../components/dialogs/WarningDialog";

/**
 * This BaseMixin contains methods that are used in multiple components
 * throughout the project.
 */
export default {
  data() {
    return {
      isLoading: true,
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
      this.$modal.show(ErrorDialog, { message }, {}, {
        "closed": () => window.location.reload(),
      });
    },
    /**
     * Shows a warning dialog with an message and a close
     * button.
     *
     * @param {String} message
     */
    showWarning(message) {
      this.$modal.show(WarningDialog, { message });
    },
    /**
     * Connects to a WebSocket with the given endpoint and setups error handlers.
     *
     * @param {String} endpoint
     */
    connectToSocket(endpoint) {
      const socket = new WebSocket(this.$config.SOCKET_ENDPOINT + endpoint);

      socket.onopen = () => {
        socket.send(localStorage.getItem("token"));
        this.isLoading = false;
      };

      socket.onclose = () => {
        this.showError("Lost connection to the server");
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
