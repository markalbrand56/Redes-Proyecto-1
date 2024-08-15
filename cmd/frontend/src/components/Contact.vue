<script setup>

import {reactive} from 'vue'

const props = defineProps({
  contact: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['setCorrespondent'])

function setCorrespondent() {
  console.log("Setting correspondent")
  emit('setCorrespondent', props.contact.jid)
}

</script>

<template>

  <div class="contact-container" @click="setCorrespondent">
    <p class="contact-jid"> {{ contact.jid }} </p>
  </div>

</template>

<style scoped>

.contact-container {
  display: flex;
  justify-content: flex-start;
  margin: 0.5rem;
  cursor: pointer;
}

.contact-jid {
  padding: 0.5rem;
  background-color: #f0f0f0;
  border-radius: 0.5rem;
  color: #333;
}

</style>