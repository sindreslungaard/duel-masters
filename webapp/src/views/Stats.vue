<template>
  <div>
    <Header></Header>

    <div class="stats-row">
      <div class="flex-1">
        <h1>Your Game Statistics</h1>
        <div class="area">
          <div class="stat-item">
            <span class="stat-label">Total Games Played:</span>
            <span class="stat-value">{{ stats.total_games_played }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Games Won:</span>
            <span class="stat-value">{{ stats.games_won }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Games Lost:</span>
            <span class="stat-value">{{ stats.games_lost }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Win Rate:</span>
            <span class="stat-value">{{ stats.win_rate }}%</span>
          </div>
          <div v-if="statsError" class="error-message">
            {{ statsError }}
          </div>
        </div>
      </div>

      <div class="flex-1"></div>

      <div class="flex-1"></div>
    </div>
  </div>
</template>

<script>
import Header from "../components/Header.vue";
import { call } from "../remote";

export default {
  name: "Stats",
  components: { Header },
  data() {
    return {
      stats: {
        total_games_played: 0,
        games_won: 0,
        games_lost: 0,
        win_rate: 0,
      },
      statsError: "",
    };
  },
  async mounted() {
    try {
      const res = await call({
        path: `/stats`,
        method: "GET",
      });

      this.stats = res.data;
    } catch (e) {
      console.error("Error fetching stats:", e);
      this.statsError =
        e.response?.data?.message ||
        e.message ||
        "Failed to fetch stats, please try again.";
    }
  },
};
</script>

<style scoped lang="scss">
.stats-row {
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

  .stat-item {
    display: flex;
    justify-content: space-between;
    padding: 5px 0;
    border-bottom: 1px solid #555;

    .stat-label {
      font-weight: bold;
    }

    .stat-value {
      font-size: 14px;
    }
  }

  .error-message {
    color: red;
    margin-top: 10px;
  }
}

.stats-row > div {
  margin: 5px;
}

.flex-1 {
  flex: 1 1 0px;
}
</style>
