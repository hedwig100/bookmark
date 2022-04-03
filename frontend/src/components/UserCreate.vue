<template>
  <div class="login">
    <p id="exp">
      You can create your account in this page. If you already have an account,
      go <router-link to="/login" id="to_login">here</router-link>. <br />
      <span id="warn">You cannot change your username after registered.</span>
    </p>
    <p v-show="isError" id="err">
      {{ errMsg }}
    </p>
    <form>
      <div class="in">
        <p>username</p>
        <input v-model="username" />
      </div>
      <div class="in">
        <p>password</p>
        <input v-model="password" />
      </div>
      <button @click="submit">register</button>
    </form>
  </div>
</template>

<script>
import { client } from "../client";

export default {
  data() {
    return {
      username: "",
      password: "",
      isError: false,
      errMsg: "",
    };
  },
  methods: {
    submit() {
      if (this.username === "") {
        this.errMsg = "Username must not be empty.";
        this.isError = true;
        return;
      }
      if (this.password === "") {
        this.errMsg = "Password must not be empty";
        this.isError = true;
        return;
      }
      client
        .post("/users", {
          username: this.username,
          password: this.password,
        })
        .then((response) => {
          console.log(response);
          if (response.status == 201) {
            // created
            console.log("created");
            this.isError = false;
            this.$router.push("/users/" + this.username);
          }
        })
        .catch((error) => {
          if (error.response.data.code == 0) {
            this.errMsg = error.response.data.message;
            this.isError = true;
          }
          if (error.response) {
            console.log(error.response.data);
            console.log(error.response.status);
            console.log(error.response.headers);
          } else if (error.request) {
            console.log(error.request);
          } else {
            console.log("Error", error.message);
          }
        });
    },
  },
};
</script>

<style scoped>
.login {
  display: flex;
  flex-flow: column;
  justify-content: flex-start;
  align-items: center;
  background-color: ghostwhite;
  padding: 0;
  margin: 0;
  width: 100vw;
  height: 100vh;
}
#exp {
  margin-top: 50px;
  padding: 0;
}
#to_login:hover {
  opacity: 0.5;
}
#warn {
  font-weight: bold;
}
#err {
  font-weight: bold;
  color: red;
}
form {
  display: block;
  border: 2px solid black;
  border-radius: 20px;
  margin: 30px;
  padding: 40px;
  width: 500px;
  height: 250px;
  background-color: white;
}
.in {
  display: flex;
  justify-content: space-around;
  margin: 20px;
  height: 40px;
}
.in p {
  font-size: 25px;
  margin: 0;
  padding: 0;
}
input {
  font-size: 1em;
  padding: 5px;
}
button {
  margin: 30px;
  padding: 16px 100px;
  background: black;
  color: white;
  font-weight: bold;
  font-size: 20px;
  border-radius: 10px;
  cursor: pointer;
}
button:hover {
  opacity: 0.5;
}
</style>
