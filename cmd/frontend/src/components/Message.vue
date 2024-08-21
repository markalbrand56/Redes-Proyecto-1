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
  },
  isConference: {
    type: Boolean,
    required: true
  }
})

message.body = props.message.body
message.timestamp = props.message.timestamp.slice(5, 16).replace("T", " ")

const isUserMessage = (props.message.from.split("/")[1] === props.user.split("@")[0] )  || (props.message.from === props.user)
const sender = props.message.from.split("/")[1]

const isImage = props.message.body.startsWith("https://") && (props.message.body.endsWith(".png") || props.message.body.endsWith(".jpg") || props.message.body.endsWith(".jpeg"))

</script>

<template>
  <div :class="['message-container', isUserMessage ? 'user-message' : 'other-message']">
    <p v-if="isConference && !isUserMessage" class="message-sender"> {{ sender }} </p>
    <img v-if="isImage" :src="message.body" alt="Image">
    <div class="inner-message">
      <p class="message-body" v-if="!isImage"> {{ message.body }}  </p>
      <p class="message-timestamp"> {{ message.timestamp }} </p>
    </div>
  </div>
</template>


<style scoped>

.message-container {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;

  margin: 0.5rem;
}

.message-container img {
  max-width: 65%;

  margin: 0.5rem;

  border-radius: 0.5rem;
  object-fit: contain;
}

.inner-message {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.message-body {
  padding: 0.5rem;
  background-color: #f0f0f0;
  border-radius: 0.5rem;
  color: #333;
  max-width: 65%;

}


.message-timestamp {
  color: #666;
  font-size: 12px;
  margin: 0 1rem;
}

.message-sender {
  color: #d5d5d5;
  font-size: 12px;
  margin: 0 1rem;

  border: 1px solid #d5d5d5;
  border-radius: 0.5rem;
  padding: 5px;
}

.message-container.user-message {
  align-items: flex-end;
  justify-content: flex-end;
}

.message-container.other-message {
  align-items: flex-start;
  justify-content: flex-start;
}

</style>