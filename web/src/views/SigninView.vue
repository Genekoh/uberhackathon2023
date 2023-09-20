<template>
  <form @submit.prevent="submitForm">
    <div class="form-container">
      <label for="email" class="form-label">Email: </label>
      <input type="text" name="email" v-model="email" />
    </div>
    <div class="form-container">
      <label for="password" class="form-label">Password: </label>
      <input type="password" name="password" v-model="password" />
    </div>
    <button type="submit" class="submit-btn">Sign In</button>
  </form>

  <RouterLink to="/signup" class="">
    <h1>Signup</h1>
  </RouterLink>
</template>

<script>
import { RouterLink } from "vue-router";
import axios from "axios";
import router from "../router";

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
      router.push({ name: "home" });
    },
  },
};
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-radius: 2rem;
  color: white;
  margin: 4rem;
  padding: 2rem 2rem 3rem 2rem;
  /* height: 60vh; */
  width: 60%;
  background-color: black;
}
.form-label {
  font-size: 2rem;
}

.submit-btn {
  font-size: 1.8rem;
}
</style>
