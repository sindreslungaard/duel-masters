<template>
  <div>
    <span class="title">Welcome, {{ username }}!</span>
    <nav>
      <ul>
        <li @click="$router.push('overview')" class="">Overview</li>
        <li class="no-cursor">|</li>
        <li @click="$router.push('decks')">Decks</li>
        <li class="no-cursor">|</li>
        <li @click="$router.push('logout')">Logout</li>
        <li class="no-cursor">|</li>
        <li>
          <a target="_blank" href="https://patreon.com/shobu"
            ><img class="patreon-btn" src="/assets/images/patreon.png"
          /></a>
        </li>
      </ul>
    </nav>

    <div class="psa">
      <span
        >Join our
        <a target="_blank" href="https://discord.gg/R3gJrX2">discord</a> to
        contribute with feedback and suggestions. Consider supporting us on
        <a target="_blank" style="color: #E66351" href="https://patreon.com/shobu">
        patreon</a> if you like what we're doing ❤️
        </span
      >
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
</style>
