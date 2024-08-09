<script setup>
import {reactive} from 'vue'
import {SendMessage, SetCorrespondent, UpdateContacts, RequestContact} from '../../wailsjs/go/main/App'
import {EventsOn} from "../../wailsjs/runtime/runtime.js";

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
  contact: ""
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

const receiveMessages = async () => {
    EventsOn("message", (message, from) => {
      data.resultText = "Message from " + from + ": " + message
  })
}

const updateContacts = async () => {
    EventsOn("contacts", (contacts) => {
      data.resultText = "Contacts: " + contacts
  })
}

const successEvent = async () => {
    EventsOn("success", (message) => {
      data.resultText = message
  })
}

const addContact = async () => {
  console.log("Adding contact")
  data.resultText = "Adding contact"
  RequestContact(data.contact)
}

receiveMessages()
updateContacts()
successEvent()

</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="sendMessage">Send</button>
      <button class="btn" @click="setCorrespondent">Set</button>
      <button class="btn" @click="getContacts">Get</button>
    </div>
    <div id="contacts" class="input-box">
      <input id="contact" v-model="data.contact" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="addContact">Add</button>
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
