class LobbyMessageHandler {
  constructor(socket, component) {
    this.socket = socket;
    this.component = component;

    this.socket.onmessage = (event) => {
      this.handleMessage(event);
    };
  }

  handleMessage(event) {
    const data = JSON.parse(event.data);
    const header = data.header;

    const handlerMap = {
      mping: () => this.handlePing(),
      hello: () => this.handleHello(),
      chat: () => this.handleChat(data),
      users: () => this.handleUsers(data),
      matches: () => this.handleMatches(data)
    };

    if (!handlerMap[header]) {
      console.error(`Received unknown header ${header}`);
    }

    handlerMap[header]();
  }

  handlePing() {
    this.component.sendMessage(this.socket, {
      header: "mpong"
    });
  }

  handleHello() {
    this.component.sendMessage(this.socket, {
      header: "subscribe"
    });
  }

  handleChat(data) {
    for (const message of data.messages) {
      this.component.addChatMessage(message);
    }
  }

  handleUsers(data) {
    const users = {
      admin: [],
      player: []
    };

    data.users.forEach(user => {
      let category = "player";

      const chatRoles = user.permissions.filter(x =>
        x.includes("chat.role.")
      );

      if (chatRoles.length > 0) {
        category = chatRoles[0].split("chat.role.")[1];
      }

      users[category].push(user);
    });

    this.component.users = users;
  }

  handleMatches(data) {
    this.component.matches = data.matches;
  }
};

export default LobbyMessageHandler;
