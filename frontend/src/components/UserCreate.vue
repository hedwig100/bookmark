<template>
  <div class="login">
    <p id="warn">
      <span>You cannot change your username after registered.</span>
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
export default {
  data() {
    return {
      username: "",
      password: "",
    };
  },
  methods: {
    submit() {
      if (this.username === "") {
        window.alert("Username must not be empty.");
        return;
      }
      if (this.password === "") {
        window.alert("Password must not be empty.");
        return;
      }
      // TODO axios /users POST
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
#warn {
  margin-top: 50px;
  padding: 0;
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
