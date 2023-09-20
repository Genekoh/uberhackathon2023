<template>
  <form @submit.prevent="submitForm">
    <label for="email">Email: </label>
    <input type="text" name="email" v-model="email" />
    <label for="password">Password: </label>
    <input type="password" name="password" v-model="password" />
    <button type="submit">Sign In</button>
  </form>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      email: "",
      password: "",
    };
  },
  methods: {
    submitForm() {
      axios
        .post(
          `http://127.0.0.1:8080/accounts/signin`,
          {
            email: this.email,
            password: this.password,
          },
          {
            withCredentials: true,
          }
        )
        .then((res) => {
          console.log(res);
        });

      this.email = "";
      this.password = "";
    },
  },
};
</script>
