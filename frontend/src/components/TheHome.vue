<template>
  <div class="hello">
    <div class="book_shelf">
      <button
        class="book"
        v-for="read in reads"
        :key="read.readId"
        @click="toBookDetail(el, read)"
      >
        <p>
          <span id="book_name">{{ read.bookName }}</span> {{ read.authorName }}
        </p>
      </button>
    </div>
    <button @click="getBooks">send a request</button>
    <p>{{ msg }}</p>
  </div>
</template>

<script>
import { client } from "../client";

export default {
  data() {
    return {
      msg: "",
      // sample dummy data, actually this data should be response data.
      reads: [
        {
          readId: "19u3fge",
          bookName: "Harry Potter",
          authorName: "J.K. Rowling",
          genres: ["fantasy", "for children"],
          thoughts: "Voldemort scared me a lot.",
          readAt: "2021-10-30",
        },
        {
          readId: "fejoie",
          bookName: "A Christmas Carol",
          authorName: "Charles Dickens",
          genres: ["for children"],
          thoughts: "I want to read at the Christmas.",
          readAt: "2022-2-30",
        },
      ],
    };
  },
  mounted() {
    // this.getBooks();
  },
  methods: {
    getBooks() {
      console.log("submit /users/a/books request.");
      client
        .get("/users/" + this.$route.params.username + "/books")
        .then((resp) => {
          console.log(resp);
          this.msg = resp.data;
        })
        .catch((err) => {
          console.log(err);
          window.alert(err);
        });
    },
    toBookDetail(el, read) {
      console.log(read);
      this.$router.push("/users/books/" + read.readId);
    },
  },
};
</script>

<style scoped>
.hello {
  min-height: 100px;
  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}
.book_shelf {
  width: 80vw;
  height: 70vh;
  border-bottom: 10px solid black;
  border-left: 10px solid black;
  border-right: 10px solid black;
  margin: 30px;
  display: flex;
  align-items: flex-end;
}
.book {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 60px;
  height: 500px;
  border: 1px solid tan;
  background-color: peru;
}
.book:hover {
  opacity: 0.5;
}
.book p {
  writing-mode: vertical-lr;
  white-space: nowrap;
}
#book_name {
  margin: 15px;
  font-size: 18px;
}
</style>
