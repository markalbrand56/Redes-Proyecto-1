<script setup>

import {onMounted, reactive, ref} from "vue";

import {CancelSubscription, DeleteConference, GetContacts, SendInvitation,} from '../../wailsjs/go/main/App.js';

const Contacts = reactive({
  contacts: [],
})

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

// Delete conference
const deleteConference = (jid) => {
  console.log("Deleting conference", jid)
  DeleteConference(jid)
}

const getContacts = async () => {
  Contacts.contacts = await GetContacts()

  console.log("Contacts: ", Contacts.contacts)
}

onMounted(() => {
  getContacts()
})

</script>

<template>
  <div class="fixed top-0 left-0 w-full h-full bg-black bg-opacity-50 flex justify-center items-center" @click="closeOptions">
    <div class="bg-white p-5 rounded-lg shadow-lg text-center w-full max-w-sm" @click.stop>
      <!--  Botones de contactos    -->
      <button v-if="!props.isConference" class="w-full my-4 py-2 px-4 rounded-md bg-blue-500 text-white hover:bg-blue-600" @click="removeContact">Remove contact</button>

      <!--  Botones de conferencias    -->
      <button v-if="props.isConference" class="w-full my-4 py-2 px-4 rounded-md bg-blue-500 text-white hover:bg-blue-600" @click="toggleInviting">Invite contact</button>

      <select v-if="props.isConference && inviting" v-model="contactInvite.jid" class="w-4/5 my-4 py-2 rounded-md bg-gray-200 text-gray-700">
        <option value="" disabled selected>Select a contact</option>
        <option v-for="contact in Contacts.contacts" :key="contact" :value="contact" class="text-gray-600">
          {{ contact }}
        </option>
      </select>

      <button v-if="props.isConference && inviting" class="w-4/5 my-4 py-2 rounded-md bg-gray-200 text-gray-700 hover:bg-gray-300" @click="inviteContact(contactInvite.jid)">Send invitation</button>

      <button v-if="props.isConference" class="w-full my-4 py-2 px-4 rounded-md bg-blue-500 text-white hover:bg-blue-600" @click="exitConference(jid)">Exit conference</button>

      <button v-if="props.isConference" class="w-full my-4 py-2 px-4 rounded-md bg-blue-500 text-white hover:bg-blue-600" @click="deleteConference(jid)">Delete conference</button>
    </div>
  </div>
</template>


<style scoped>


</style>