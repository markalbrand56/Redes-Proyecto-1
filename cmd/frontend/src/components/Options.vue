<script setup>

const props = defineProps({
  jid: {
    type: String,
    required: true
  },
  isConference: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['remove-contact', 'exit-conference', 'invite-contact', 'close-options'])

const closeOptions = () => {
  emit('close-options')
}

</script>

<template>

  <div class="options-container" @click="closeOptions">
    <div class="options-content" @click.stop>
      <button v-if="!props.isConference" class="btn" @click="emit('remove-contact', jid)">Remove contact</button>

      <button v-if="props.isConference" class="btn" @click="emit('invite-contact', jid)">Invite to conference</button>

      <button v-if="props.isConference" class="btn" @click="emit('exit-conference', jid)">Exit conference</button>

    </div>
  </div>


</template>

<style scoped>

.options-container {
  position: fixed;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;

  min-height: 40%;

  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.options-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  width: fit-content;

  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.btn {
  display: block;
  width: 100%;

  margin: 1rem;
  padding: 1rem 2rem;

  border: none;
  border-radius: 4px;
  background-color: #007bff;
  color: white;
}

</style>