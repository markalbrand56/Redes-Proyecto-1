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
  }
})

message.body = props.message.body
// 2024-08-14T01:26:58.7680233-06:00 format as 08-14 01:26
message.timestamp = props.message.timestamp.slice(5, 16).replace("T", " ")

</script>

<template>
  <div class="message-container">
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

</style>