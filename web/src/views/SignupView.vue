<template>
  <form @submit.prevent="submitForm">
    <div class="form-container">
      <label for="username" class="form-label">Username: </label>
      <input type="text" name="username" v-model="username" />
    </div>
    <div class="form-container">
      <label for="firstname" class="form-label">First Name: </label>
      <input type="text" name="firstname" v-model="firstname" />
    </div>
    <div class="form-container">
      <label for="lastname" class="form-label">Last Name: </label>
      <input type="text" name="lastname" v-model="lastname" />
    </div>
    <div class="form-container">
      <label for="email" class="form-label">Email: </label>
      <input type="text" name="email" v-model="email" />
    </div>
    <div class="form-container">
      <label for="salary" class="form-label">Salary: </label>
      <input type="number" name="salary" v-model="salary" />
    </div>
    <div class="form-container">
      <label for="password" class="form-label">Password: </label>
      <input type="password" name="password" v-model="password" />
    </div>
    <button type="submit" class="submit-btn">Sign Up</button>
  </form>
</template>

<script>
import axios from "axios";
import router from "../router";

export default {
  data() {
    return {
      username: "",
      firstname: "",
      lastname: "",
      email: "",
      password: "",
      salary: 0,
    };
  },
  methods: {
    submitForm() {
      axios
        .post(
          `http://127.0.0.1:8080/accounts/signup`,
          {
            username: this.username,
            firstname: this.firstname,
            lastname: this.lastname,
            email: this.email,
            password: this.password,
            salary: this.salary,
          },
          {
            withCredentials: true,
          }
        )
        .then((res) => {
          console.log(res);
        });

      this.username = "";
      this.firstname = "";
      this.lastname = "";
      this.email = "";
      this.password = "";
      this.salary = 0;
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
