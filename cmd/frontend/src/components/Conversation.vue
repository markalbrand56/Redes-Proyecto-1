<script setup>

import {reactive} from 'vue'

import Message from "./Message.vue";
import {models} from "../../wailsjs/go/models.ts";

const props = defineProps({
  user: {
    type: String,
    required: true
  },
  messages: {
    type: Array[models.Message],
    required: true
  },
  isConference: {
    type: Boolean,
    required: true
  }
})

</script>

<template>
  <div class="conversation-container">
    <Message v-for="message in props.messages" :key="message.timestamp" :message="message"  :user="props.user" :is-conference="props.isConference"/>
  </div>
</template>

<style scoped>

.conversation-container {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;

  height: fit-content;
  min-height: 100%;

  gap: 1rem;
  background-color: #383838;
  padding-top: 5rem;
  overflow-y: hidden;

}

</style>