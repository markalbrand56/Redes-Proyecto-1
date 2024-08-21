<script setup>

import { ChevronRightIcon } from '@heroicons/vue/24/solid'

const props = defineProps({
  contact: {
    type: Object,
    required: true
  },
  alias: {
    type: String,
    required: false
  }
})

const emit = defineEmits(['setCorrespondent'])

function setCorrespondent() {
  console.log("Setting correspondent")
  emit('setCorrespondent', props.contact.jid)
}

const firstLetter = props.contact.jid.charAt(0).toUpperCase()
const usernameDisplay = props.contact.jid.split('@')[0]
</script>

<template>
  <div class="contact-container" @click="setCorrespondent">
    <div class="contact-icon">{{ firstLetter }}</div>
<!--    <p class="contact-jid">{{ usernameDisplay  }}</p>-->
<!--    si hay alias, se puestra en vez del jid-->
    <p class="contact-jid">{{ alias ? alias : usernameDisplay }}</p>
    <ChevronRightIcon class="arrow-right" />
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
  width: 30px;
  height: 30px;

  background-color: #007bff;
  color: white;

  display: flex;
  justify-content: center;
  align-items: center;

  border-radius: 50%;
  font-size: 1.25rem;
  margin: 0.5rem 1.25rem 0.5rem 0.5rem;
}

.contact-jid {
  color: #333;
  margin: 0.25rem;
}

.arrow-right {
  height: 16px;
  width: 16px;
  margin-left: auto;
  color: #666;
  font-size: 1.15rem;
  font-weight: bold;
}
</style>
