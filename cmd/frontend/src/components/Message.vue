<script setup>

import {reactive} from 'vue'

const message = reactive({
  body: ""
})

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

// set the message body from the props
message.body = props.message.body

</script>

<template>
  <div class="message-container">
    <p class="message-body"> {{ message.body }} </p>
  </div>
</template>

<style scoped>

.message-container {
  display: flex;
  justify-content: flex-start;
  margin: 0.5rem;
}

.message-body {
  padding: 0.5rem;
  background-color: #f0f0f0;
  border-radius: 0.5rem;
  color: #333;
}

</style>