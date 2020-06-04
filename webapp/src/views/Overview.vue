<template>
  <div>

      <div v-show="errorMessage || wizardVisible" class="overlay"></div>

      <div v-show="errorMessage" class="error">
        <p>{{ errorMessage }}</p>
        <div @click="refreshPage()" class="btn">Reconnect</div>
      </div>

      <div v-show="wizardVisible" class="new-duel">
          <div class="wizard">
              <div class="spacer">
                <span class="headline">Create a new duel</span>
                <br><br>
                <form>
                    <input v-model="wizard.name" type="text" placeholder="Name">
                    <br><br>
                    <span class="helper">Visibility</span>
                    <select v-model="wizard.visibility">
                        <option value="public">Show in list of duels</option>
                        <option value="private">Hide from list of duels</option>
                    </select>

                    <span v-if="wizardError" class="errorMsg">{{ wizardError }}</span>

                    <div @click="createDuel()" class="btn">
                        Create
                    </div>
                    <div @click="toggleWizard()" class="btn cancel">
                        Cancel
                    </div>
                </form>
              </div>
          </div>
      </div>

      <main>

				  <Header style="width: 100%"></Header>

					<div class="spaced">
                        

						<div class="categories">
							<h3 class="user-list">Online</h3>
							<h3 class="chat">Chat</h3>
							<h3 class="duels" style="position: relative;">Duels<span @click="toggleWizard()" class="new-duel-btn">New Duel</span></h3>
						</div>


						<!-- Users online -->
						<div class="box user-list">
							<div class="spaced">

                                <div v-if="wsLoading">Loading{{ loadingDots }}</div>

                                <div class="user-name-container" v-for="(category, index) in users" :key="index">

                                    <div class="user-category"><span>{{ category.category }}</span></div>

                                    <Username v-for="(user, index) in category.users" :key="index" :hub="user.hub" :color="user.color">{{ user.username }}</Username>
                                    
                                    <br>

                                </div>

							</div>
						</div>






						<!-- Chat -->
						<div class="box chat">

                            <div v-if="wsLoading" class="spaced" style="position: absolute">Loading{{ loadingDots }}</div>

							<div class="chatbox">

								<div id="messages" class="messages spaced">
									<div class="messages-helper">

                                        <div v-for="(msg, i) in chatMessages" :key="i">
                                            <Username :color="msg.color">{{ msg.username }}</Username>
                                            <div class="user-messages">
                                                <div v-for="(message, j) in msg.messages" :key="j">
                                                    <span>{{ message }}</span>
                                                </div>
                                            </div>
                                        </div>

									</div>
								</div>
								<form @submit.prevent="sendChat(chatMessage)">
									<input type="text" v-model="chatMessage" placeholder="Type to chat">
								</form>  
							</div>

						</div>








						<!-- Duels -->
						<div class="box duels">

                            <div v-if="wsLoading" class="spaced">Loading{{ loadingDots }}</div>

							<table>
                                <tr v-if="!wsLoading && matches.length < 1"><td>No matches to show, click the button above to create one.</td></tr>
								<tr v-for="(match, index) in matches" :key="index">
									<td><Username :color="match.color">{{ match.owner }}</Username></td>
									<td>{{ match.name }}</td>
									<td><div @click="$router.push('/duel/' + match.id)" :class="'btn' + (match.spectate ? '' : ' save')">{{ match.spectate ? "Spectate" : "Join match" }}</div></td>
								</tr>
							</table>

						</div>


					</div>

      </main>   

  </div>
</template>

<script>
import { call } from '../remote'
import Header from '../components/Header.vue'
import Username from '../components/Username.vue'

const send = (client, message) => {
  client.send(JSON.stringify(message))
}

export default {
  name: 'overview',
  components: {
			Header,
			Username
  },
  computed: {
      username: () => localStorage.getItem('username')
  },
  data() {
      return {
          ws: null,
          wizardVisible: false,
          wizardError: "",
          wizard: {
              name: "",
              description: "",
              visibility: "public"
          },
          chatMessage: "",
          chatMessages: [],
          users: [],
          matches: [],
          errorMessage: "",
          wsLoading: true,
          loadingDots: "."
      }
  },
  methods: {
      refreshPage() {
          location.reload()
      },
      toggleWizard() {
          this.wizardError = ""
          this.wizard = {
              name: "",
              description: "",
              visibility: "public"
          }
          this.wizardVisible = !this.wizardVisible
      },
      async createDuel() {

          if(this.wizard.name.length < 5 || this.wizard.length > 30) {
              this.wizardError = "Duel name must be between 5-30 characters"
              return
          }

          try {
            let res = await call({
                path: "/match",
                method: "POST",
                body: this.wizard
            })

            this.$router.push({ path: '/duel/' + res.data.id })

          } catch(e) {
              try {
                  console.log(e)
                  this.wizardError = e.response.data.message
              } catch(err) {
                  console.log(err)
                  this.wizardError = "Unable to communicate with the server. Please try again later."
              }
          }

      },
      sendChat(message) {
        if(!message) {
            return
        }
        this.chatMessage = ""
        this.ws.send(JSON.stringify({ header: "chat", message }))
      },
      chat(data) {
        
        let createNew = true

        if(this.chatMessages.length > 0) {
            let lastMsg = this.chatMessages[this.chatMessages.length - 1]

            if(lastMsg.username == data.username && lastMsg.timestamp > data.timestamp - 15) {
                lastMsg.messages.push(data.message)
                createNew = false
            }
        }

        if(createNew) {
            this.chatMessages.push({
                username: data.username,
                color: data.color,
                timestamp: data.timestamp,
                messages: [data.message]
            })
        }
        
        this.$nextTick(() => {
            let container = document.getElementById('messages')
            container.scrollTop = container.scrollHeight
        })
      }
  },
  created() {

        // Loading dots
        setInterval(() => {
            if(this.loadingDots.length >= 4)
                this.loadingDots = ""
            else this.loadingDots += "."
        }, 500)

        // Connect to the server
        try {
            const ws = new WebSocket("ws://" + window.location.hostname + "/ws/lobby")
            this.ws = ws

            ws.onopen = () => {
                ws.send(localStorage.getItem("token"))
                this.wsLoading = false
            }

            ws.onclose = () => {
                this.errorMessage = "Lost connection to the server"
            }

            ws.onerror = () => {
                this.errorMessage = "Lost connection to the server"
            }

            ws.onmessage = (event) => {

                const data = JSON.parse(event.data)

                switch(data.header) {

                    case "mping": {
                        send(ws, {
                            header: "mpong"
                        })
                        break
                    }

                    case "hello": {
                        send(ws, {
                            header: "subscribe"
                        })
                        break
                    }

                    case "chat": {
                        for(let message of data.messages) {
                            this.chat(message)
                        }
                        break
                    }

                    case "users": {
                        this.users = [{
                            category: "player",
                            users: []
                        }]

                        for(let user of data.users) {
                            
                            let chatroles = user.permissions.filter(x => x.includes("chat.role."))

                            if(chatroles.length > 0) {

                                let role = chatroles[0].split("chat.role.")[1]

                                let category = this.users.find(x => x.category == role)

                                let toPushCategory = false
                                if(!category) {
                                    category = {
                                        category: role,
                                        users: []
                                    }
                                    toPushCategory = true
                                }

                                category.users.push(user)
                                if(toPushCategory) {
                                    this.users.push(category)
                                }

                            } else {
                                let category = this.users.find(x => x.category == "player")
                                category.users.push(user)
                            }

                            this.users.sort((a, b) => a.category.localeCompare(b.category))

                        }
                        break
                    }

                    case "matches": {
                        this.matches = data.matches
                        break
                    }

                }

            }
        }
        catch(err) {
            this.errorMessage = "Connection lost"
        }

  },
  beforeDestroy() {
    this.ws.close()
  }
}
</script>

<style scoped lang="scss">

.disabled {
    background: #7289DA !important;
    opacity: 0.5;
}

.disabled:hover {
    cursor: not-allowed !important;
    background: #7289DA !important;
}

.disabled:active {
    background: #7289DA !important;
}

.new-duel .backdrop {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    background: #000;
    opacity: 0.5;
}

.new-duel .wizard {
    position: absolute;
    top: calc(50vh - 323px / 2);
    left: calc(50% - 250px / 2);
    background: #36393F;
    width: 250px;
    border-radius: 4px;
    color: #fff;
    border: 1px solid #666;
    z-index: 100;
}

.wizard .headline {
    color: #ccc;
}

.wizard .spacer {
    margin: 15px;
}

.wizard .helper {
    color: #ccc;
    font-size: 13px;
}

.wizard .btn {
    margin: 0;
    width: 85px;
    text-align: center;
    margin-top: 15px;
}

.wizard .cancel {
    margin-left: 10px;
    background: #FF4C4C;
    color: #fff;
}

.wizard .cancel:hover {
    background: #ed3e3e;
}

input, textarea, select {
    border: none;
    background: #484C52;
    padding: 10px;
    border-radius: 4px;
    width: 200px;
    color: #ccc;
    resize: none;
}
input:focus, textarea:focus, select:focus {
    outline: none
}
input:active, textarea:active, select:active {
    outline: none
}

.wizard select {
    width: 220px;
    margin-top: 4px;
}

.wizard .errorMsg {
    color: red;
    font-size: 14px;
    display: block;
    margin-top: 15px;
}

nav {
    text-align: right;
}

ul {
    list-style: none;
}

li {
    display: inline-block;
    padding-right: 10px;
    margin-right: 10px;
}

nav > ul > li:hover {
    cursor: pointer;
    color: #fff;
}

nav > ul > li.no-cursor:hover {
    cursor: default
}

.title {
    position: absolute;
    top: 16px;
    left: 16px;
}

.psa {
    margin: 16px;
    background: #2B2C31;
    padding: 5px;
    min-height: 20px;
    border-radius: 4px;
    font-size: 14px;
    color: #ccc;
}

.psa > span {
    display: inline-block;
    vertical-align: middle;
    margin-left: 4px;
}

a {
    color: #7289DA;
}

.btn {
  display: inline-block;
  background: #7289DA;
  color: #E3E3E5;
  font-size: 14px;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn:hover {
  cursor: pointer;
  background: #677BC4;
}

.btn:active {
  background: #5B6EAE;
}

main {
		width: 100%;
		height: 100vh;
		margin: auto;
		overflow: hidden;
}

.box {
		overflow: auto;
		background: #2B2C31;
		min-height: 20px;
		border-radius: 4px;
		font-size: 14px;
		color: #ccc;
		display: inline-block;
		height: calc(100vh - 140px)
}

.user-list {
	width: 10%;
}

.chat {
	width: calc(35% - 15px);
	margin-left: 15px;
}

.duels {
	width: calc(55% - 15px);
	margin-left: 15px;
}

.spaced {
	margin: 15px;
}

.categories > h3 {
	margin-top: 0;
	margin-bottom: 7px;
	display: inline-block;
	color: #eee;
	font-weight: 400;
	font-size: 16px;
}

.duels > table {
	width: 100%;
	border-collapse: collapse;
}

.duels td {
	border: none;
  text-align: left;
  padding: 15px;
}

.duels tr:nth-child(odd) {
  background-color: #222429;
}

.duels .btn {
	float: right;
}

.save {
    background: #3CA374 !important;
}

.save:hover {
    background: #35966a !important;
}

.user-category {
	margin-bottom: 10px;
	border-bottom: 1px solid #555;
	color: #777;
	padding-bottom: 5px;
	font-weight: 400;
    text-transform: capitalize;
}

.chatbox {
	display: flex;
  flex-direction: column;
  justify-content: space-between;
	height: 100%;
    overflow: hidden;
}

.chatbox input {
  border: none;
  border-radius: 4px;
  margin: 10px;
  width: calc(100% - 40px);
  background: #484C52;
  padding: 10px;
  color: #ccc;
  &:focus {
    outline: none
  }
  &:active {
    outline: none
  }
}

.duels .btn {
	width: 70px;
}

.user-list .user-name {
	margin-bottom: 10px;
}

.user-name-container {
    overflow: hidden;
}

.user-messages {
	margin-left: 20px;
	margin-top: 0px;
	margin-bottom: 15px;
}

.user-messages > div {
	margin: 3px 0;
	color: #e1e1e1;
}

.messages {
    overflow: auto;
    margin-bottom: 0;
    padding-bottom: 0;
}

*::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    border-radius: 10px;
    background-color: #484C52;
  }

*::-webkit-scrollbar {
  width: 6px;
  height: 6px;
  background-color: #484C52;
}

*::-webkit-scrollbar-thumb {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    background-color: #222;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: #000;
  opacity: 0.5;
  z-index: 100;
}

.error p {
  padding: 5px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: #2B2E33 !important;
  border: 1px solid #222428;
}

.error {
  border: 1px solid #666;
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #36393F;
  z-index: 3005;
  left: calc(50% - 300px / 2);
  top: 40vh;
  padding: 10px;
  font-size: 14px;
  color: #ccc;
}

.new-duel-btn {
    margin-left: 7px;
    color: #fff;
    background: #3CA374;
    padding: 3px 5px;
    border-radius: 4px;
    font-size: 12px;
    text-transform: uppercase;
    font-weight: 400;
    position: absolute;
    top: -1px;
}

.new-duel-btn:hover {
    cursor: pointer;
    background: #35966A;
}


</style>
