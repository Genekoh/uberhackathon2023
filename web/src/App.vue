<script setup>
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
  <nav>
    <RouterLink to="/" class="nav-item title">
      <h1>Uber Car</h1>
    </RouterLink>

    <div class="left-nav">
      <RouterLink to="/signin" class="account-item btn nav-item">
        <i class="fa-solid fa-user"
          ><p class="txt"><b>Account</b></p></i
        >
      </RouterLink>
      <div class="menu-item nav-item">
        <i class="fa-solid fa-bars"></i>
      </div>
    </div>
  </nav>
  <main class="page-content">
    <RouterView />
  </main>
</template>

<style scoped>
nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: black;
  height: 5rem;
  padding-left: 2rem;
  padding-right: 2rem;
}

.title {
  color: white;
  text-decoration: none;
}

.title h1 {
  font-size: 3rem;
}

.left-nav {
  display: flex;
  align-items: center;
  gap: 2.5rem;
}

.fa-solid {
  color: white;
  font-size: 2.3rem;
  transition: opacity 0.15s;
}

.fa-user {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.txt {
  font-size: 0.8rem;
  margin-top: 1px;
  color: white;
  transition: opacity 0.15s;
}

.btn {
  cursor: pointer;
  border: none;
  background-color: black;
  transition: box-shadow 0.2s;
}

.btn:hover .fa-solid {
  opacity: 0.8;
}
.btn:active .fa-solid {
  color: grey;
  opacity: 0.3;
}

.btn:active .txt {
  color: grey;
  opacity: 0.3;
}
</style>

<style>
* {
  padding: 0px;
  margin: 0px;
  box-sizing: border-box;
  font-family: "Aoboshi One";
}
</style>
