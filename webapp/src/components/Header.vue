<template>
  <div>
    <span class="title">Welcome, {{ username }}!</span>
    <nav>
      <ul>
        <li @click="$router.push('overview')" class="">Overview</li>
        <li class="no-cursor">|</li>
        <li @click="$router.push('decks')">Decks</li>
        <li class="no-cursor">|</li>
        <li @click="$router.push('settings')">Settings</li>
        <li class="no-cursor">|</li>
        <li @click="$router.push('logout')">Logout</li>
        <li class="no-cursor">|</li>
        <li>
          <a
            target="_blank"
            href="https://github.com/sindreslungaard/duel-masters"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="github-icon"
              width="24"
              height="24"
              viewBox="0 0 24 24"
            >
              <path
                fill="#fff"
                d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
              /></svg
          ></a>
        </li>
      </ul>
    </nav>

    <div class="psa">
      <span
        >Join our
        <a target="_blank" href="https://discord.gg/R3gJrX2">discord</a> to
        contribute with feedback and suggestions. The project is open source and
        available on
        <a
          target="_blank"
          href="https://github.com/sindreslungaard/duel-masters"
          >github</a
        >
        if you're interested in checking it out.
      </span>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";

export default {
  name: "decks",
  computed: {
    username: () => localStorage.getItem("username")
  },
  data() {
    return {};
  },
  methods: {}
};
</script>

<style scoped>
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
  background: #36393f;
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
  background: #ff4c4c;
  color: #fff;
}

.wizard .cancel:hover {
  background: #ed3e3e;
}

input,
textarea,
select {
  border: none;
  background: #484c52;
  padding: 10px;
  border-radius: 4px;
  width: 200px;
  color: #ccc;
  resize: none;
}
input:focus,
textarea:focus,
select:focus {
  outline: none;
}
input:active,
textarea:active,
select:active {
  outline: none;
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
  cursor: default;
}

.title {
  position: absolute;
  top: 16px;
  left: 16px;
}

.psa {
  margin: 16px;
  background: #2b2c31;
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
  color: #7289da;
}

.btn {
  display: inline-block;
  background: #7289da;
  color: #e3e3e5;
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
  background: #677bc4;
}

.btn:active {
  background: #5b6eae !important;
}

.patreon-btn {
  width: 125px;
  margin-bottom: -10px;
  border-radius: 4px;
  margin-right: -3px;
  opacity: 0.9;
}

.patreon-btn:hover {
  opacity: 1;
}

.github-icon {
  margin-bottom: -8px;
  opacity: 0.8;
}

.github-icon:hover {
  opacity: 0.9;
}
</style>
