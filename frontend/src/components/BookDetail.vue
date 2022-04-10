<template>
  <div class="hello">
    <button @click="submit">send a request</button>
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
    submit() {
      console.log("submit /hello request.");
      client
        .get("/hello")
        .then((resp) => {
          console.log(resp);
          this.msg = resp.data;
        })
        .catch((err) => {
          console.log(err);
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
