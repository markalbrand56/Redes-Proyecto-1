<script setup>
import {reactive, onMounted, nextTick, ref} from 'vue'
import {
  SendMessage,
  UpdateContacts,
  RequestContact,
  AcceptSubscription,
  CancelSubscription,
  SetStatus,
  GetMessages,
  GetMessagesConference,
  GetArchive
} from '../../wailsjs/go/main/App.js'

import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {models} from "../../wailsjs/go/models.ts";

import Conversation from "../components/Conversation.vue";
import Contact from "../components/Contact.vue";

const Message = reactive({
  jid: "",
  body: "",
  isConference: false
})

const User = reactive({
  jid: "alb21004@alumchat.lol",
  contacts: [],
  conferences: {},
  status: 0
})

const Messages = reactive({
  messages: []
})

const Debug = reactive({
  resultText: "Please enter your name below 👇",
})

const messageSectionRef = ref(null)

function scrollToBottom() {
  nextTick(() => {
    if (messageSectionRef.value) {
      messageSectionRef.value.scrollTop = messageSectionRef.value.scrollHeight
    }
  })
}

function sendMessage() {
  // TODO Implementar nuevos mensajes para el usuario cuando falte algo en el formulario
  if (Message.body === "") {
    Debug.resultText = "Please enter a message to send"
    return
  }

  if (Message.jid === "") {
    Debug.resultText = "Please select a contact to send a message to"
    return
  }

  // Body, to, from
  SendMessage(Message.body, Message.jid, User.jid)
  Message.body = ""
  getMessages()
}

function getContacts() {
  console.log("Getting contacts")
  Debug.resultText = "Getting contacts"
  UpdateContacts()
}

function addContact() {
  console.log("Adding contact")
  Debug.resultText = "Adding contact"
  RequestContact(Message.jid)
}

function cancelSubscription() {
  console.log("Cancelling subscription")
  Debug.resultText = "Cancelling subscription"
  CancelSubscription(Message.jid)
}

function updateStatus(status) {
  console.log("Updating status")
  Debug.resultText = "Updating status"
  SetStatus(status)
}

function getMessages() {
  console.log("Getting messages")
  GetMessages(Message.jid).then((messages) => {
    if (messages.length > 0) {

      Messages.messages = messages.map((message) => {
        return new models.Message(message)
      })

      scrollToBottom()
      // order messages by timestamp
      Messages.messages.sort((a, b) => {
        return new Date(a.timestamp) - new Date(b.timestamp)
      })

      console.log("Messages", Messages.messages)
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
  Message.jid = jid  // Set the current correspondent on the frontend
  Message.isConference = false
  Debug.resultText = "Setting correspondent to " + jid

  getArchive(jid)  // Get the messages for the current correspondent
}

function handleConferenceClicked(jid) {
  console.log("Conference clicked", jid)

  Message.jid = jid  // Set the current correspondent on the frontend
  Message.isConference = true
  Debug.resultText = "Setting correspondent to " + jid

  GetMessagesConference(jid).then((messages) => {
    if (messages.length > 0) {

      console.log("Messages received from conference", messages)

      Messages.messages = messages.map((message) => {
        return new models.Message(message)
      })

      scrollToBottom()
      // order messages by timestamp
      Messages.messages.sort((a, b) => {
        return new Date(a.timestamp) - new Date(b.timestamp)
      })

    } else {
      Messages.messages = []
    }
  })
}

// Event listeners

const receiveMessages = async () => {
  EventsOn("message", (from) => {
    console.log("Message", from)
    Debug.resultText = "Message from " + from
    if (from === Message.jid) {
      console.log("Updating current conversation")
      getArchive(from)
    }
  })
}

const updateContacts = async () => {
  EventsOn("contacts", (contacts) => {
    // contacts is an array of strings
    Debug.resultText = "Contacts: " + contacts.join(", ")
    console.log("Contacts", contacts)
    User.contacts = contacts
  })
}

const updateConferences = async () => {
  EventsOn("conferences", (conferences) => {
    // conferences is a map conferences[item.Name] = item.Jid
    console.log("Conferences", conferences)
    Debug.resultText = "Conferences: " + Object.keys(conferences).join(", ")

    User.conferences = conferences

  })
}

const successEvent = async () => {
  EventsOn("success", (message) => {
    Debug.resultText = message
  })
}

const subRequest = async () => {
  EventsOn("subscription-request", (user) => {
    Debug.resultText = "Subscription request from " + user
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
updateConferences()

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
        <div id="correspondents" class="correspondents">
          <h2>Contacts</h2>
          <div id="contacts" class="contact-section">
            <Contact v-for="contact in User.contacts" :contact="{jid: contact}" :key="contact" @setCorrespondent="handleContactClicked" />
          </div>

          <h2>Group chats</h2>
          <div id="conferences" class="contact-section">
            <Contact v-for="(jid, name) in User.conferences" :contact="{jid: jid}" :key="jid" @setCorrespondent="handleConferenceClicked" />
          </div>
        </div>

        <div id="current-account" class="current-account">
          <p>{{ User.jid }}</p>
        </div>
      </div>

      <div id="current-chat" class="current-chat">

        <div id="current-contact" class="current-contact">
          <p>{{ Message.jid }}</p>
        </div>

        <div id="messages" class="message-section" ref="messageSectionRef">
          <Conversation :messages="Messages.messages"  :user="User.jid" :is-conference="Message.isConference"/>
        </div>

        <div id="message-input" class="message-input">
          <input id="message" v-model="Message.body" autocomplete="off" class="input" type="text"/>
          <button class="btn" @click="sendMessage">Send</button>
        </div>

      </div>

    </div>

    <div id="debug" class="debug">

      <div id="result" class="result">{{ Debug.resultText }}</div>

      <div id="contacts-debug" class="input-box">
        <input id="contact" v-model="User.jid" autocomplete="off" class="input" type="text"/>
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

.correspondents {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;

  width: calc(100% - 2rem);
  height: 80%;
  margin: 1rem;
  border: 1px dashed green;

  box-sizing: border-box;

  overflow-x: hidden;
  overflow-y: scroll;
  scrollbar-width: thin;
  scrollbar-color: #000000 #464646;
}

.correspondents::-webkit-scrollbar {
  width: 10px;
}

.correspondents::-webkit-scrollbar-thumb {
  background-color: #000000;
}

.correspondents::-webkit-scrollbar-track {
  background-color: #464646;
}

.correspondents h2 {
  margin: 0.25rem;

  font-size: min(1.25rem, 2vw);
  color: #1b2636;

  border-bottom: 1px solid #1b2636;
  padding: 0.5rem;
  width: 70%;
  text-align: center;
  background-color: #f0f0f0;
  border-radius: 0.75rem;
}

.contact-section {
  width: 100%;
  border: 1px dashed red;
}

.current-account {
  display: flex;
  justify-content: center;
  align-items: center;

  min-height: 10%;
  width: 100%;
  height: fit-content;
  margin: 1rem;

  border: 1px dashed blue;

  background-color: #f0f0f0;
  border-radius: 0.75rem;
  cursor: pointer;

}

.current-account p {
  margin: 0.5rem;
  font-size: 18px;
  color: #1b2636;
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