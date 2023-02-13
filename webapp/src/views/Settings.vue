<template>
  <div>
    <Header></Header>

    <div class="settings">
      <div style="width: 205px">
        <h1>Change Password</h1>
        <div class="area">
          <form @submit.prevent="changePassword">
            <div class="form-control">
              <label for="oldPassword">Old password</label>
              <input
                id="oldPassword"
                v-model="oldPassword"
                type="password"
                placeholder="********"
              />
            </div>

            <br />

            <div class="form-control">
              <label for="newPassword">New password</label>
              <input
                id="newPassword"
                v-model="newPassword"
                type="password"
                placeholder="********"
              />
            </div>

            <br />

            <div class="form-control">
              <label for="newPasswordAgain">New password again</label>
              <input
                id="newPasswordAgain"
                v-model="newPasswordAgain"
                type="password"
                placeholder="********"
              />
            </div>

            <br />

            <button type="submit">Change password</button>

            <div v-if="passwordError" style="color: red; margin-top: 10px;">
              {{ passwordError }}
            </div>

            <div v-if="passwordSuccess" style="color: green; margin-top: 10px;">
              {{ passwordSuccess }}
            </div>
          </form>
        </div>
      </div>

      <div class="flex-1">
        <h1>Muted players</h1>
        <div class="area">
          <div class="section-info">
            {{
              settings.muted.length
                ? ""
                : "You have not muted anyone in the chat 😇"
            }}
          </div>

          <div v-for="player in settings.muted" :key="player">
            <div class="player">
              <span class="name">{{ player }}</span>
              <MuteIcon :player="player" @toggled="refresh()" />
            </div>
          </div>
        </div>
      </div>

      <div class="flex-1">
        <h1>Toggles</h1>
        <div class="area">
          <div class="checkbox-container">
            <input
              type="checkbox"
              id="noUpsideDownCards"
              :checked="settings.noUpsideDownCards"
              @change="setUpsideDownCards($event)"
            />
            <label for="noUpsideDownCards">
              Flip your manazone cards and the opponent's battlezone and
              shieldzone cards so they are readable by you
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Header from "../components/Header";
import MuteIcon from "../components/MuteIcon";
import { getSettings, setSettings, patchSettings } from "../helpers/settings";
import { call } from "../remote";

export default {
  name: "settings",
  components: { Header, MuteIcon },
  data() {
    return {
      oldPassword: "",
      newPassword: "",
      newPasswordAgain: "",
      passwordError: "",
      passwordSuccess: "",

      settings: getSettings()
    };
  },
  methods: {
    refresh() {
      this.settings = getSettings();
    },

    setUpsideDownCards(e) {
      patchSettings({ noUpsideDownCards: e.target.checked });
    },

    async changePassword() {
      this.passwordError = "";
      this.passwordSuccess = "";

      if (this.newPassword !== this.newPasswordAgain) {
        this.passwordError = "New password and password again does not match";
        return;
      }

      if (this.newPassword.length < 6) {
        this.passwordError = "New password must be at least 6 characters long";
        return;
      }

      try {
        let res = await call({
          path: `/auth/reset-password`,
          method: "POST",
          body: {
            oldPassword: this.oldPassword,
            newPassword: this.newPassword
          }
        });

        this.passwordSuccess = res.data.message;
      } catch (e) {
        this.passwordError = e.response.data.error;
      }
    }
  },

  created() {
    addEventListener("storage", this.refresh);
  },

  beforeDestroy() {
    removeEventListener("storage", this.refresh);
  }
};
</script>

<style scoped lang="scss">
.settings {
  margin: 0 10px;
  display: flex;

  .area {
    background: #2b2c31;
    padding: 10px;
    border-radius: 4px;
    margin-top: 10px;
    font-size: 14px;
  }

  h1 {
    font-weight: normal;
    font-size: 16px;
    margin: 0;
  }

  input {
    padding: 7px;
    border-radius: 3px;
    border: none;
    outline: none;
    background: #ccc;
  }

  button {
    width: 100%;
    padding: 5px;
  }
}

.settings > div {
  margin: 5px;
}

.flex-1 {
  flex-grow: 1;
}

.form-control {
  label {
    padding-bottom: 5px;
    display: block;
  }
}

.player {
  display: flex;
  align-items: center;
  margin-bottom: 5px;

  border-bottom: 1px solid #555;
  padding-bottom: 5px;

  .name {
    font-size: 14px;
    margin-right: 10px;
  }
}
</style>
