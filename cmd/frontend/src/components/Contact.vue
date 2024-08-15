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

const firstLetter = props.contact.jid.charAt(0).toUpperCase()
</script>

<template>
  <div class="contact-container" @click="setCorrespondent">
    <div class="contact-icon">{{ firstLetter }}</div>
    <p class="contact-jid">{{ contact.jid }}</p>
    <div class="arrow-right">{{ ">" }}</div>
  </div>
</template>

<style scoped>
.contact-container {
  display: flex;
  align-items: center;
  margin: 0.5rem;
  cursor: pointer;

  background-color: #f0f0f0;
  border-radius: 0.5rem;
  padding: 0 0.5rem;
}

.contact-icon {
  width: 40px;
  height: 40px;

  background-color: #007bff;
  color: white;

  display: flex;
  justify-content: center;
  align-items: center;

  border-radius: 50%;
  font-size: 1.25rem;
  margin-right: 0.75rem;
}

.contact-jid {
  color: #333;
  margin: 1rem 0.5rem;
}

.arrow-right {
  margin-left: auto;
  color: #666;
  font-size: 1.15rem;
  font-weight: bold;
}
</style>
