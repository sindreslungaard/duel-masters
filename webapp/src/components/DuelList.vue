<template>
  <Panel title="Duels" class="duels">
    <LoadingIndicator v-if="isLoading" />

    <Button @click="$emit('newDuel')" text="New Duel" />
    <table>
      <tr v-if="!isLoading && matches.length < 1">
        <td>No matches to show, click the button above to create one.</td>
      </tr>
      <tr v-for="(match, index) in matches" :key="index">
        <td>
          <Username :color="match.color">{{ match.owner }}</Username>
        </td>
        <td>{{ match.name }}</td>
        <td>
          <div
            @click="$router.push('/duel/' + match.id)"
            :class="'btn' + (match.spectate ? '' : ' save')"
          >
            {{ match.spectate ? "Spectate" : "Join match" }}
          </div>
        </td>
      </tr>
    </table>
  </Panel>
</template>

<script>
import Panel from "./Panel";
import Button from "./buttons/Button";
import LoadingIndicator from "./LoadingIndicator";

export default {
  name: "DuelList",
  components: {
    Panel,
    Button,
    LoadingIndicator
  },
  props: {
    matches: {
      type: Array,
      required: true
    },
    isLoading: {
      type: Boolean,
      required: true
    }
  }
};
</script>

<style lang="scss" scoped>
.duels table {
  width: 100%;
  margin-top: 15px;
  border-collapse: collapse;
}

.duels td {
  border: none;
  text-align: left;
  padding: var(--spacing);
}

.duels tr:nth-child(odd) {
  background-color: #222429;
}

.save {
  background: #3ca374 !important;
}

.save:hover {
  background: #35966a !important;
}

.duels .btn {
  width: 70px;
}
</style>
