<template>
  <div>
    <Header></Header>

    <div class="settings-row">
      <div class="flex-1">
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

            <button class="btn" type="submit">Change password</button>

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
              Flip all cards that would otherwise be upside down
            </label>
          </div>
        </div>
      </div>
    </div>

    <div class="settings-row">
      <div class="flex-1">
        <h1>Preferences</h1>
        <div class="area">
          <form @submit.prevent="savePreferences">
            <div class="form-control">
              <label for="playmat">Playmat</label>
              <input
                id="playmat"
                v-model="preferences.playmat"
                type="text"
                placeholder="https://i.imgur.com/HDpcrUt.png"
                style="width: calc(100% - 15px)"
              />
            </div>

            <br />

            <button class="btn" type="submit">Update preferences</button>

            <div v-if="preferencesError" style="color: red; margin-top: 10px;">
              {{ preferencesError }}
            </div>

            <div
              v-if="preferencesSuccess"
              style="color: green; margin-top: 10px;"
            >
              {{ preferencesSuccess }}
            </div>
          </form>
        </div>
      </div>

      <div class="flex-1"></div>

      <div class="flex-1"></div>
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

      preferences: {
        playmat: ""
      },
      preferencesError: "",
      preferencesSuccess: "",

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
        this.passwordError = e.response.data.message;
      }
    },
    async savePreferences() {
      this.preferencesError = "";
      this.preferencesSuccess = "";

      try {
        let res = await call({
          path: `/preferences`,
          method: "PUT",
          body: this.preferences
        });

        this.preferencesSuccess = res.data.message;
      } catch (e) {
        console.log(e.response.data.error);
        this.preferencesError = e.response.data.error;
      }
    }
  },

  async mounted() {
    try {
      let res = await call({
        path: `/preferences`,
        method: "GET"
      });

      this.preferences = res.data;
    } catch (e) {
      alert("failed to fetch preferences, try to log out and back in again");
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
.settings-row {
  margin: 0 10px;
  display: flex;
  margin-bottom: 50px;

  .area {
    background: url(/assets/images/overlay_30.png);
    padding: 10px;
    border-radius: 4px;
    margin-top: 10px;
    font-size: 14px;
    height: calc(100% - 45px);
  }

  h1 {
    font-weight: normal;
    font-size: 16px;
    margin: 0;
  }

  input {
    color: #fff;
    padding: 7px;
    border-radius: 3px;
    border: none;
    outline: none;
    background: url(/assets/images/overlay_50.png);
  }

  button {
    padding: 5px 36px;
  }
}

.settings-row > div {
  margin: 5px;
}

.flex-1 {
  flex: 1 1 0px;
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

.btn {
  outline: none;
  border: none;
  display: inline-block;
  background: #5865F2;
  color: #e3e3e5;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn:hover {
  cursor: pointer;
  background: #515de2;
}

.btn:active {
  background: #4c58d3 !important;
}
</style>
