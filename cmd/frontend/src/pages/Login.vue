<script setup>
import { reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Login } from '../../wailsjs/go/main/App.js';
import {EventsOn} from "../../wailsjs/runtime/runtime.js";

import PulseLoader from 'vue-spinner/src/PulseLoader.vue';
import Swal from "sweetalert2";


const user = reactive({
  username: "",
  password: "",
  loggingIn: false
});

const router = useRouter();

async function handleLogin() {
  console.log(user.username, user.password);
  user.loggingIn = true;
  await Login(user.username, user.password);
}

// Event listeners

EventsOn("login", (jid) => {
  user.loggingIn = false;
  console.log("Login successful for jid: ", jid);
  router.push("/chat");
});

EventsOn("login-error", (error) => {
  user.loggingIn = false;
  Swal.fire({
    icon: 'error',
    title: 'Login failed',
    text: error,
  });
  console.error("Login error: ", error);
});


</script>

<template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="handleLogin">
      <div>
        <label for="username">Username:</label>
        <input type="text" v-model="user.username" id="username" required />
      </div>
      <div>
        <label for="password">Password:</label>
        <input type="password" v-model="user.password" id="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
    <pulse-loader v-if="user.loggingIn" color="#007bff" size="10px" class="loader"/>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

form div {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

input {
  padding: 0.5rem;
  border-radius: 0.25rem;
  border: 1px solid #ccc;
}

button {
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  border: none;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.loader {
  margin-top: 1rem;
}
</style>