<script setup>

import {reactive} from 'vue'

import Message from "./Message.vue";
import {models} from "../../wailsjs/go/models.ts";

const props = defineProps({
  messages: {
    type: Array[models.Message],
    required: true
  }
})

</script>

<template>
  <div class="conversation-container">
    <Message v-for="message in messages" :key="message.timestamp" :message="message" />
  </div>
</template>

<style scoped>

.conversation-container {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  height: 100%;

  gap: 1rem;
  background-color: #383838;
  padding-top: 5rem;
  overflow-y: hidden;

}

</style>