<script setup>

import { ChevronRightIcon } from '@heroicons/vue/24/solid'
import {computed} from "vue";

const props = defineProps({
  contact: {
    type: Object,
    required: true
  },
  status: {
    type: String,
    required: false
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

const statusColor = computed(() => {
  switch (props.status) {
    case 'Online':  //  Online
      return 'green'

    case 'Disconnected':  //  Disconnected / Invisible
      return 'gray'

    case 'Away':  //  Away
      return 'yellow'

    case 'Do Not Disturb':  //  Busy
      return 'red'

    case 'Extended Away':  //  Extended Away
      return 'orange'

    default:
      return 'green'
  }
})

const firstLetter = props.contact.jid.charAt(0).toUpperCase()
const usernameDisplay = props.contact.jid.split('@')[0]
</script>

<template>
  <div class="contact-container" @click="setCorrespondent">
    <div class="contact-icon">{{ firstLetter }}</div>
    <div :class="['status-indicator', statusColor]"></div>
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

.status-indicator {
  width: 12px;
  height: 12px;

  border-radius: 50%;
  border: 1px solid #000000;

  margin-right: 0.5rem;
}

.green {
  background-color: green;
}

.red {
  background-color: red;
}

.yellow {
  background-color: #ceb300;
}

.orange {
  background-color: #ff4c00;
}

.gray {
  background-color: gray;
}

</style>
