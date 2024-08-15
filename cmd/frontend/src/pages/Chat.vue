<script setup>
import {reactive, onMounted, nextTick, ref} from 'vue'
import {
  SendMessage,
  SetCorrespondent,
  UpdateContacts,
  RequestContact,
  AcceptSubscription,
  CancelSubscription,
  SetStatus,
  GetMessages,
  GetArchive
} from '../../wailsjs/go/main/App.js'

import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {models} from "../../wailsjs/go/models.ts";

import Conversation from "../components/Conversation.vue";
import Contact from "../components/Contact.vue";

const Message = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
  contact: "",
  body: ""
})

const Correspondent = reactive({
  jid: "",
  name: ""
})

const Contacts = reactive({
  contacts: []
})

const Messages = reactive({
  messages: []
})

const messageSectionRef = ref(null)

function scrollToBottom() {
  nextTick(() => {
    if (messageSectionRef.value) {
      messageSectionRef.value.scrollTop = messageSectionRef.value.scrollHeight
    }
  })
}

function setCorrespondent(jid) {
  console.log("Setting correspondent")
  Message.resultText = "Setting correspondent to " + jid
  SetCorrespondent(jid)
}

function sendMessage() {
  if (Message.body === "") {
    Message.resultText = "Please enter a message to send"
    return
  }

  if (Correspondent.jid === "") {
    Message.resultText = "Please select a contact to send a message to"
    return
  }

  SendMessage(Message.body)
  Message.body = ""
  getMessages()
}

function getContacts() {
  console.log("Getting contacts")
  Message.resultText = "Getting contacts"
  UpdateContacts()
}

function addContact() {
  console.log("Adding contact")
  Message.resultText = "Adding contact"
  RequestContact(Message.contact)
}

function cancelSubscription() {
  console.log("Cancelling subscription")
  Message.resultText = "Cancelling subscription"
  CancelSubscription(Message.contact)
}

function updateStatus(status) {
  console.log("Updating status")
  Message.resultText = "Updating status"
  SetStatus(status)
}

function getMessages() {
  console.log("Getting messages")
  GetMessages(Correspondent.jid).then((messages) => {
    console.log("Messages", messages)
    if (messages) {
      Messages.messages = messages.map((message) => {
        return new models.Message(message)
      })
      scrollToBottom()
    } else {
      Messages.messages = []
    }
  })
}

function getArchive(jid) {
  console.log("Getting archive")
  GetArchive(jid)
}

function handleContactClicked(jid) {
  console.log("Contact clicked", jid)
  setCorrespondent(jid)  // Set the current correspondent on the backend
  Correspondent.jid = jid  // Set the current correspondent on the frontend

  getArchive(jid)  // Get the messages for the current correspondent
}

// Event listeners

const receiveMessages = async () => {
  EventsOn("message", (from) => {
    console.log("Message", from)
    Message.resultText = "Message from " + from
    if (from === Correspondent.jid) {
      console.log("Updating current conversation")
      getArchive(from)
    }
  })
}

const updateContacts = async () => {
  EventsOn("contacts", (contacts) => {
    // contacts is an array of strings
    Message.resultText = "Contacts: " + contacts.join(", ")
    console.log("Contacts", contacts)
    Contacts.contacts = contacts
  })
}

const successEvent = async () => {
  EventsOn("success", (message) => {
    Message.resultText = message
  })
}

const subRequest = async () => {
  EventsOn("subscription-request", (user) => {
    Message.resultText = "Subscription request from " + user
    AcceptSubscription(user)
  })
}

const updateMessages = async () => {
  EventsOn("update-messages", (jid) => {
    console.log("Updating messages for", jid)
    getMessages()
  })
}

receiveMessages()
updateContacts()
successEvent()
subRequest()
updateMessages()

getContacts()

onMounted(() => {
  scrollToBottom()
})

</script>

<template>
  <main>
    <h1>Chat</h1>
    <div id="display" class="display">

      <div id="left-panel" class="left-panel">
        <div id="contacts" class="contact-section">
          <Contact v-for="contact in Contacts.contacts" :contact="{jid: contact}" :key="contact" @setCorrespondent="handleContactClicked" />
        </div>
      </div>

      <div id="current-chat" class="current-chat">

        <div id="current-contact" class="current-contact">
          <p>{{ Correspondent.jid }}</p>
        </div>

        <div id="messages" class="message-section" ref="messageSectionRef">
          <Conversation :messages="Messages.messages" />
        </div>

        <div id="message-input" class="message-input">
          <input id="message" v-model="Message.body" autocomplete="off" class="input" type="text"/>
          <button class="btn" @click="sendMessage">Send</button>
        </div>

      </div>

    </div>

    <div id="debug" class="debug">

      <div id="result" class="result">{{ Message.resultText }}</div>

      <div id="input" class="input-box">
        <input id="name" v-model="Message.name" autocomplete="off" class="input" type="text"/>
        <button class="btn" @click="sendMessage">Send</button>
        <button class="btn" @click="setCorrespondent">Set</button>
        <button class="btn" @click="getContacts">Get</button>
      </div>

      <div id="contacts-debug" class="input-box">
        <input id="contact" v-model="Message.contact" autocomplete="off" class="input" type="text"/>
        <button class="btn" @click="addContact">Add</button>
        <button class="btn" @click="cancelSubscription">Remove</button>
      </div>

      <div id="status" class="input-box">
        <button class="btn" @click="updateStatus(0)">Online</button>
        <button class="btn" @click="updateStatus(1)">Away</button>
        <button class="btn" @click="updateStatus(2)">Busy</button>
        <button class="btn" @click="updateStatus(3)">NA</button>
        <button class="btn" @click="updateStatus(4)">Offline</button>

      </div>

    </div>


  </main>
</template>

<style scoped>

main {
  display: flex;
  flex-direction: column;

  justify-content: center;
  align-items: center;
  height: 100%;
}

main h1 {
  margin: 1rem;
}

.display {
  display: flex;
  justify-content: space-between;
  width: 100%;
  height: 60%;
}

.left-panel {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;

  width: 20%;
  height: 100%;
  margin: 1rem;
  border: 1px dashed green;
}

.contact-section {
  width: 100%;
  height: 100%;
  border: 1px dashed red;
}

.current-chat {
  width: 75%;
  margin: 1rem;
  border: 1px dashed blue;
}

.current-contact {
  display: flex;
  justify-content: center;
  align-items: center;

  min-height: 10%;
  height: fit-content;
  margin: 1rem;
}

.current-contact p {
  margin: 0.5rem;
  font-size: 18px;
}

.message-section {
  height: 70%;

  margin: 2rem;
  border: 2px solid #000000;

  overflow-y: scroll;
  scrollbar-width: thin;
  scrollbar-color: #000000 #464646;

}

.message-input {
  display: flex;
  justify-content: center;
  align-items: center;

  margin: 2rem;
}

.message-input .input {
  width: 80%;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.message-input .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.message-input .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.message-input .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box {
  padding: 1em;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>