<script setup>
import { reactive, ref, onMounted } from 'vue'
import { models } from "../../wailsjs/go/models.ts";
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

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

const isUserMessage = (props.message.from.split("/")[1] === props.user.split("@")[0]) || (props.message.from === props.user)
const sender = props.message.from.split("/")[1]

const isImage = props.message.body.startsWith("https://") && (props.message.body.endsWith(".png") || props.message.body.endsWith(".jpg") || props.message.body.endsWith(".jpeg"))
const imageLoaded = ref(false)

function handleImageLoad() {
  console.log("Image loaded")
  imageLoaded.value = true
}

// Para asegurar que la imagen comience a cargar en `onMounted`
onMounted(() => {
  if (isImage) {
    const img = new Image();
    img.src = message.body;
    img.onload = handleImageLoad;
    img.onerror = () => {
      console.error("Error loading image");
      imageLoaded.value = true; // Evita el loader infinito en caso de error
    }
  }
})

</script>

<template>
  <div :class="['message-container', isUserMessage ? 'user-message' : 'other-message']">
    <p v-if="isConference && !isUserMessage" class="message-sender"> {{ sender }} </p>
    <PulseLoader v-if="isImage && !imageLoaded" color="#007bff" size="10px" class="mt-4 min-w-32"/>
    <img v-if="isImage && imageLoaded" :src="message.body" alt="Image">
    <div class="inner-message">
      <p class="message-body" v-if="!isImage"> {{ message.body }}  </p>
      <p class="message-timestamp"> {{ message.timestamp }} </p>
    </div>
  </div>
</template>

<style scoped>
.message-container {
  @apply flex flex-col justify-start my-2;
}

.message-container img {
  @apply max-w-[65%] max-h-[300px] my-2 mr-4 rounded-lg object-contain;
}

.inner-message {
  @apply flex flex-row items-center;
}

.message-body {
  @apply p-2 bg-gray-100 rounded-lg text-gray-800 max-w-[65%] ml-4 my-2;
}

.message-timestamp {
  @apply text-gray-500 text-xs mx-4;
}

.message-sender {
  @apply text-gray-400 text-sm ml-4 border border-gray-400 rounded-lg p-1;
}

.message-container.user-message {
  @apply items-end justify-end;
}

.message-container.other-message {
  @apply items-start justify-start;
}
</style>
