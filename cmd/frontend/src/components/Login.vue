<script setup>
import { reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Login } from '../../wailsjs/go/main/App';
import {EventsOn} from "../../wailsjs/runtime/runtime.js";


const user = reactive({
  username: "",
  password: ""
});

const router = useRouter();

async function handleLogin() {
  console.log(user.username, user.password);
  await Login(user.username, user.password);
}

// Event listeners

EventsOn("success", () => {
  console.log("Login successful");
  router.push("/chat");
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
</style>