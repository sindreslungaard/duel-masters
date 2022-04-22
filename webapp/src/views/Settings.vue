<template>
  <div>
    <Header></Header>
    <div class="container">
      <div class="section-info">
        {{
          mutedPlayers.length
            ? "Muted players in chat:"
            : "You have not muted anyone in the chat 😇"
        }}
      </div>

      <ul>
        <li v-for="player in mutedPlayers" :key="player">
          <div class="player">
            <span class="name">{{ player }}</span>
            <MuteIcon :player="player" @toggled="refresh()" />
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import Header from "../components/Header";
import MuteIcon from "../components/MuteIcon";
import { getMutedPlayers } from "../helpers/mute";

export default {
  name: "settings",
  components: { Header, MuteIcon },
  data() {
    return {
      mutedPlayers: getMutedPlayers()
    };
  },
  methods: {
    refresh() {
      this.mutedPlayers = getMutedPlayers();
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
}
</style>
