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
  <div class="flex flex-col items-center justify-center h-full">
    <h2 class="text-4xl text-blue-500 mb-8 font-bold">Login</h2>
    <form @submit.prevent="handleLogin" class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <label for="username" class="text-xl">Username:</label>
        <input type="text" v-model="user.username" id="username" class="p-2 border border-gray-300 rounded text-gray-600" required />
      </div>
      <div class="flex flex-col gap-2">
        <label for="password" class="text-xl">Password:</label>
        <input type="password" v-model="user.password" id="password" class="p-2 border border-gray-300 rounded text-black" required />
      </div>
      <button type="submit" class="p-2 bg-blue-500 text-white rounded hover:bg-blue-600 m-2">Login</button>
    </form>
    <pulse-loader v-if="user.loggingIn" color="#007bff" size="10px" class="mt-4"/>
    <button @click="router.push('/register')" class="text-blue-500 hover:underline">Register</button>
  </div>
</template>