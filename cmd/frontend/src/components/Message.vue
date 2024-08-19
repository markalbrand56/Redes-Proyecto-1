<script setup>

import {reactive} from 'vue'
import {models} from "../../wailsjs/go/models.ts";

const message = reactive({
  body: "",
  timestamp: ""
})

const props = defineProps({
  message: {
    type: models.Message,
    required: true
  },
  user: {
    type: String,
    required: true
  }
})

message.body = props.message.body
message.timestamp = props.message.timestamp.slice(5, 16).replace("T", " ")

const isUserMessage = props.message.from === props.user

</script>

<template>
  <div
      :class="['message-container', isUserMessage ? 'user-message' : 'other-message']">
    <p class="message-body"> {{ message.body }}  </p>
    <p class="message-timestamp"> {{ message.timestamp }} </p>
  </div>
</template>


<style scoped>

.message-container {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;

  margin: 0.5rem;
}

.message-body {
  padding: 0.5rem;
  background-color: #f0f0f0;
  border-radius: 0.5rem;
  color: #333;
}

.message-timestamp {
  color: #666;
  font-size: 12px;
  margin: 0 1rem;
}

.message-container.user-message {
  justify-content: flex-end;
}

.message-container.other-message {
  justify-content: flex-start;
}

</style>