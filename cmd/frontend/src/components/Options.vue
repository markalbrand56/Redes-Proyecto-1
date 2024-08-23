<script setup>

import {reactive, ref} from "vue";

import {
  CancelSubscription,
  SendInvitation,
} from '../../wailsjs/go/main/App.js';

const contactInvite = reactive({
  jid: "",
})

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

const inviting = ref(false)

const toggleInviting = () => {
  inviting.value = !inviting.value
}

// Invite contact to conference
const inviteContact = (jid) => {
  console.log("Inviting contact", jid, "to conference", props.jid)

  inviting.value = false

  SendInvitation(props.jid, jid);
}

// Remove contact from roster
const removeContact = () => {
  console.log("Removing contact", props.jid)
  // emit('remove-contact', jid)

  CancelSubscription(props.jid)
}

// Exit conference
const exitConference = (jid) => {
  console.log("Exiting conference", jid)
}

</script>

<template>

  <div class="options-container" @click="closeOptions">
    <div class="options-content" @click.stop>
      <!--  Botones de contactos    -->
      <button v-if="!props.isConference" class="btn" @click="removeContact">Remove contact</button>

      <!--  Botones de conferencias    -->
      <button v-if="props.isConference" class="btn" @click="toggleInviting">Invite contact</button>
      <input v-if="props.isConference && inviting" v-model="contactInvite.jid" type="text" placeholder="Enter JID" class="request-input" />
      <button v-if="props.isConference && inviting" class="btn-secondary" @click="inviteContact(contactInvite.jid)">Send invitation</button>

      <button v-if="props.isConference" class="btn" @click="exitConference(jid)">Exit conference</button>

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
  cursor: pointer;
}

.request-input {
  width: 80%;
  min-width: fit-content;

  margin: 1rem;
  padding: 0.5rem;

  border: none;
  border-radius: 4px;
  background-color: #f0f0f0;
  color: #333;
}

.btn-secondary {
  display: block;
  width: 80%;
  min-width: fit-content;

  margin: 1rem;
  padding: 0.5rem;

  border: none;
  border-radius: 4px;
  background-color: #f0f0f0;
  color: #333;
  cursor: pointer;
}

</style>