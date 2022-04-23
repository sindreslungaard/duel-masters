<template>
  <div class="login-page">
    <div class="form">
      <form @submit.prevent="submit()" class="login-form">
        <p class="title">Sign in to Shobu.io</p>
        <input v-model="username" type="text" placeholder="Username" />
        <input v-model="password" type="password" placeholder="Password" />
        <button>Sign in</button>
        <p v-if="errorMsg" class="error">{{ errorMsg }}</p>
        <p class="message">
          Not registered?
          <router-link
            :to="{
              path: '/register',
              query: { redirect_to: $route.query.redirect_to }
            }"
            >Create an account</router-link
          >
        </p>
      </form>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";

export default {
  name: "login",
  data() {
    return {
      username: "",
      password: "",
      errorMsg: "",
      redirectTo: null
    };
  },
  methods: {
    async submit() {
      try {
        let res = await call({
          path: "/auth/signin",
          method: "POST",
          body: {
            username: this.username,
            password: this.password
          }
        });

        localStorage.setItem("email", res.data.user.email);
        localStorage.setItem("username", res.data.user.username);
        localStorage.setItem("uid", res.data.user.uid);
        localStorage.setItem("permissions", res.data.user.permissions);
        localStorage.setItem("token", res.data.token);

        if (this.redirectTo) {
          this.$router.push(this.redirectTo);
        } else {
          this.$router.push("overview");
        }
      } catch (e) {
        try {
          if (e.response.status == 401) {
            this.errorMsg = "Wrong username or password";
          } else if (e.response.status == 403) {
            this.errorMsg = "You have been banned";
          }
        } catch (err) {
          this.errorMsg =
            "An unexpected error occured. Please try again later.";
        }
      }
    }
  },
  created() {
    console.log(this.$route);
    this.redirectTo = this.$route.query.redirect_to;
  }
};
</script>

<style scoped>
.login-page {
  color: #333;
}

.title {
  margin: 0;
  margin-bottom: 20px;
  padding: 0;
  text-align: left;
}

.error {
  font-size: 14px;
  color: red;
  margin: 0;
  margin-top: 20px;
}

.login-page {
  width: 360px;
  padding: 30vh 0 0;
  margin: auto;
}
.form {
  position: relative;
  z-index: 1;
  background: #ffffff;
  max-width: 360px;
  margin: 0 auto 100px;
  padding: 45px;
  padding-bottom: 35px;
  text-align: center;
  box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
.form input {
  font-family: "Roboto", sans-serif;
  outline: 0;
  background: #f2f2f2;
  width: 100%;
  border: 0;
  margin: 0 0 15px;
  padding: 15px;
  box-sizing: border-box;
  font-size: 14px;
}
.form button {
  font-family: "Roboto", sans-serif;
  text-transform: uppercase;
  outline: 0;
  background: #e22a38;
  width: 100%;
  border: 0;
  padding: 15px;
  color: #ffffff;
  font-size: 14px;
  -webkit-transition: all 0.3 ease;
  transition: all 0.3 ease;
  cursor: pointer;
}
.form button:hover,
.form button:active,
.form button:focus {
  background: #db2533;
}
.form .message {
  margin: 25px 0 0;
  margin-bottom: 0;
  padding-bottom: 0;
  color: #b3b3b3;
  font-size: 12px;
}
.form .message a {
  color: #e22a38;
  text-decoration: none;
}

.container {
  position: relative;
  z-index: 1;
  max-width: 300px;
  margin: 0 auto;
}
.container:before,
.container:after {
  content: "";
  display: block;
  clear: both;
}
.container .info {
  margin: 50px auto;
  text-align: center;
}
.container .info h1 {
  margin: 0 0 15px;
  padding: 0;
  font-size: 36px;
  font-weight: 300;
  color: #1a1a1a;
}
.container .info span {
  color: #4d4d4d;
  font-size: 12px;
}
.container .info span a {
  color: #000000;
  text-decoration: none;
}
.container .info span .fa {
  color: #ef3b3a;
}
</style>
