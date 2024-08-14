<script setup>
import {reactive} from 'vue'
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
import Conversation from "../components/Conversation.vue";

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
  contact: ""
})

const messages = reactive({
  messages: []
})

function setCorrespondent() {
  console.log("Setting correspondent")
  data.resultText = "Setting correspondent"
  SetCorrespondent("alb21005@alumchat.lol")
}

function sendMessage() {
  SendMessage(data.name)
}

function getContacts() {
  console.log("Getting contacts")
  data.resultText = "Getting contacts"
  UpdateContacts()
}

function addContact() {
  console.log("Adding contact")
  data.resultText = "Adding contact"
  RequestContact(data.contact)
}

function cancelSubscription() {
  console.log("Cancelling subscription")
  data.resultText = "Cancelling subscription"
  CancelSubscription(data.contact)
}

function updateStatus(status) {
  console.log("Updating status")
  data.resultText = "Updating status"
  SetStatus(status)
}

const getMessages = async () => {
    const msg = await GetMessages("alb21005@alumchat.lol/gajim.0O3D5ZZ0")
    console.log(msg)

    messages.messages = msg
}

function getArchive() {
  GetArchive(data.name)
}

// Event listeners

const receiveMessages = async () => {
    EventsOn("message", (message, from) => {
      data.resultText = "Message from " + from + ": " + message.body
      console.log("Message from " + from + ": " + message.body + " at " + message.timestamp)
  })
}

const updateContacts = async () => {
    EventsOn("contacts", (contacts) => {
      // contacts is an array of strings
      data.resultText = "Contacts: " + contacts.join(", ")
  })
}

const successEvent = async () => {
    EventsOn("success", (message) => {
      data.resultText = message
  })
}

const subRequest = async () => {
    EventsOn("subscription-request", (user) => {
      data.resultText = "Subscription request from " + user
      AcceptSubscription(user)
  })


}

receiveMessages()
updateContacts()
successEvent()
subRequest()

getMessages()

// ************************************************************ //
</script>

<template>
  <main>
    <h1>Chat</h1>
    <Conversation :messages="messages.messages" />
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="sendMessage">Send</button>
      <button class="btn" @click="setCorrespondent">Set</button>
      <button class="btn" @click="getContacts">Get</button>
      <button class="btn" @click="getArchive">GetA</button>
    </div>
    <div id="contacts" class="input-box">
      <input id="contact" v-model="data.contact" autocomplete="off" class="input" type="text"/>
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

  </main>
</template>

<style scoped>
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
