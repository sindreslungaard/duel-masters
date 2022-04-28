<template>
  <div>
    <Header></Header>
    <div class="container">
      <div class="section-info">
        {{
          settings.muted.length
            ? "Muted players in chat:"
            : "You have not muted anyone in the chat 😇"
        }}
      </div>

      <ul>
        <li v-for="player in settings.muted" :key="player">
          <div class="player">
            <span class="name">{{ player }}</span>
            <MuteIcon :player="player" @toggled="refresh()" />
          </div>
        </li>
      </ul>

      <div class="section-info">
        No upside down cards
      </div>

      <div class="checkbox-container">
        <input
          type="checkbox"
          id="noUpsideDownCards"
          :checked="settings.noUpsideDownCards"
          @change="setUpsideDownCards($event)"
        />
        <label for="noUpsideDownCards">
          Flip your manazone cards and the opponent's battlezone and shieldzone cards so they are
          readable by you
        </label>
      </div>
    </div>
  </div>
</template>

<script>
import Header from "../components/Header";
import MuteIcon from "../components/MuteIcon";
import { getSettings, setSettings, patchSettings } from "../helpers/settings";

export default {
  name: "settings",
  components: { Header, MuteIcon },
  data() {
    return {
      settings: getSettings()
    };
  },
  methods: {
    refresh() {
      this.settings = getSettings();
    },

    setUpsideDownCards(e) {
      patchSettings({ noUpsideDownCards: e.target.checked });
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
.container {
  margin: 0 15px;

  .section-info {
    padding-left: 10px;
    padding-top: 16px;
    font-size: 16px;
  }

  .player {
    display: flex;
    align-items: center;
    margin-bottom: 5px;

    .name {
      font-size: 14px;
      margin-right: 10px;
    }
  }

  .checkbox-container {
    font-size: 14px;
    margin-left: 20px;
    margin-top: 15px;

    input {
      margin-left: 0;
    }
  }
}
</style>
