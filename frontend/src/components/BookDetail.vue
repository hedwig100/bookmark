<template>
  <div class="hello">
    <button @click="getBookDetail">send a request</button>
    <p>{{ msg }}</p>
  </div>
</template>

<script>
import { client } from "../client";

export default {
  name: "BookDetail",
  data() {
    return {
      msg: "",
    };
  },
  mounted() {
    // this.submit();
  },
  methods: {
    getBookDetail() {
      console.log("submit /users/:username/books/:readId request.");
      client
        .get(
          `/users/${this.$route.params.username}/books/${this.$route.params.readId}`
        )
        .then((resp) => {
          console.log(resp);
          this.msg = resp.data;
        })
        .catch((err) => {
          console.log(err.message);
          window.alert(err);
        });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.hello {
  min-height: 100px;
}
</style>
