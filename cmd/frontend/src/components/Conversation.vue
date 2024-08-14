<script setup>

import {reactive} from 'vue'

import Message from "./Message.vue";

const conversation = reactive({
  messages: []
})

const props = defineProps({
  messages: {
    type: Array,
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
  gap: 1rem;
  padding: 1rem;
  border: 1px solid #ccc;
  border-radius: 0.5rem;
  background-color: #383838;
}

</style>