<template>
  <div>

      <div v-show="wizardVisible" class="new-duel">
          <div class="backdrop"></div>
          <div class="wizard">
              <div class="spacer">
                <span class="headline">Create a new duel</span>
                <br><br>
                <form>
                    <input v-model="wizard.name" type="text" placeholder="Name">
                    <br><br>
                    <textarea v-model="wizard.description" rows="3" cols="50" placeholder="Description"></textarea>
                    <br><br>
                    <span class="helper">Visibility</span>
                    <select v-model="wizard.visibility">
                        <option value="public">Public</option>
                        <option value="private">Private</option>
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

      <Header></Header>

      <div class="main">
          <div @click="toggleWizard()" class="btn">NEW DUEL</div>
      </div>   

  </div>
</template>

<script>
import { call } from '../remote'
import Header from '../components/Header.vue'

export default {
  name: 'overview',
  components: {
      Header
  },
  computed: {
      username: () => localStorage.getItem('username')
  },
  data() {
      return {
          wizardVisible: false,
          wizardError: "",
          wizard: {
              name: "",
              description: "",
              visibility: "public"
          }
      }
  },
  methods: {
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
          try {
            let res = await call({
                path: "/match",
                method: "POST",
                body: this.wizard
            })

            this.$router.push({ path: '/duel/' + res.data.matchUid + "/" +  res.data.inviteId})

          } catch(e) {
              try {
                  console.log(e)
                  this.wizardError = e.response.data.message
              } catch(err) {
                  console.log(err)
                  this.wizardError = "Unable to communicate with the server. Please try again later."
              }
          }

      }
  }
}
</script>

<style scoped>

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

.main {
    margin: 0 15px;
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

</style>
