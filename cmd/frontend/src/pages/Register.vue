<script setup>

import { reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Register } from '../../wailsjs/go/main/App.js';
import {EventsOn} from "../../wailsjs/runtime/runtime.js";

import PulseLoader from 'vue-spinner/src/PulseLoader.vue';
import Swal from "sweetalert2";
import logo from "../assets/images/logo.png";

const newUser = reactive({
  email: "",
  password: "",
  registering: false
});

const router = useRouter();

async function handleRegister() {
  console.log(newUser.email, newUser.password);
  newUser.registering = true;
  const res = await Register(newUser.email, newUser.password);  // Boolean

  if (res) {

   // show pop with button to go to login
    Swal.fire({
      icon: 'success',
      title: 'Register successful',
      text: 'You can now login',
      showCancelButton: true,
      confirmButtonText: 'Login',
      cancelButtonText: 'Close',
    }).then((result) => {
      if (result.isConfirmed) {
        router.push("/");
      }
    });

  } else {

    // Error
    newUser.registering = false;
    Swal.fire({
      icon: 'error',
      title: 'Register failed',
      text: 'Your registration failed, try with another username',
    });

  }
}

</script>

<template>

  <div class="flex flex-col items-center justify-center h-fit">
    <img :src="logo" alt="logo" class="w-96 mt-4"/>
    <h2 class="text-4xl text-blue-500 mb-8 font-bold">Register</h2>
    <form @submit.prevent="handleRegister" class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <label for="email" class="text-xl">Email:</label>
        <input type="email" v-model="newUser.email" id="email" class="p-2 border border-gray-300 rounded text-black" required />
      </div>
      <div class="flex flex-col gap-2">
        <label for="password" class="text-xl">Password:</label>
        <input type="password" v-model="newUser.password" id="password" class="p-2 border border-gray-300 rounded text-black" required />
      </div>
      <button type="submit" class="p-2 bg-blue-500 text-white rounded hover:bg-blue-600 m-2">Register</button>
    </form>
    <button @click="router.push('/')"  class="text-blue-500 hover:underline mt-2">Return </button>
    <pulse-loader v-if="newUser.registering" color="#007bff" size="10px" class="mt-4"/>
  </div>
  
</template>

<style scoped>

</style>