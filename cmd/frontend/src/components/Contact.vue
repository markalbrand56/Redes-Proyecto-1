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
  @apply flex items-center my-2 cursor-pointer bg-gray-100 rounded-lg px-2;
}

.contact-icon {
  @apply w-8 h-8 bg-blue-600 text-white flex justify-center items-center rounded-full text-xl mr-5;
}

.contact-jid {
  @apply text-gray-800 my-1;
}

.arrow-right {
  @apply h-4 w-4 ml-auto text-gray-600 text-lg font-bold;
}

.status-indicator {
  @apply w-3 h-3 rounded-full border border-black mr-2;
}

.green {
  @apply bg-green-500;
}

.red {
  @apply bg-red-500;
}

.yellow {
  @apply bg-yellow-500;
}

.orange {
  @apply bg-orange-500;
}

.gray {
  @apply bg-gray-500;
}
</style>

