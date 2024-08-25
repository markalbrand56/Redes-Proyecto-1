<script setup>
import {reactive, onMounted, nextTick, ref} from 'vue'
import {
  SendMessage,
  SendFileMessage,
  SendConferenceMessage,
  SendConferenceFileMessage,

  SetStatus,
  GetMessages,
  GetMessagesConference,
  UpdateContacts,
  AcceptSubscription,
  GetCurrentUser,
} from '../../wailsjs/go/main/App.js'

import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {models} from "../../wailsjs/go/models.ts";

import Conversation from "../components/Conversation.vue";
import Nav from "../components/Nav.vue";
import Options from "../components/Options.vue";
import FileInput from "../components/FileInput.vue";
import LeftPanel from "../components/LeftPanel.vue";

import {CogIcon, Bars3Icon, PaperAirplaneIcon} from "@heroicons/vue/24/solid";
import {useRouter} from "vue-router";

const router = useRouter();

const Message = reactive({
  jid: "",
  body: "",
  isConference: false,
  statusMessage: ""
})

const User = reactive({
  jid: "",
})

const Messages = reactive({
  messages: []
})

const showLeftPanel = ref(false)

const Debug = reactive({
  resultText: "Please enter your name below 👇",
})

const messageSectionRef = ref(null)  // Reference to the message section
const showOptions = ref(false);  // Show options for current conversation

function scrollToBottom() {
  nextTick(() => {
    if (messageSectionRef.value) {
      messageSectionRef.value.scrollTop = messageSectionRef.value.scrollHeight
    }
  })
}

const toggleOptions = () => {
  console.log('Toggling options');
  showOptions.value = !showOptions.value;
};


function sendMessage() {
  if (Message.body === "") {
    return
  }

  if (Message.jid === "") {
    return
  }

  if (Message.isConference) {
    SendConferenceMessage(Message.body, Message.jid)
  } else {
    SendMessage(Message.body, Message.jid, User.jid)
  }

  // Body, to, from
  Message.body = ""
}

function getContacts() {
  console.log("Getting contacts")
  UpdateContacts()
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

function getConferenceMessages(jid) {
  console.log("Getting conference messages")
  GetMessagesConference(jid).then((messages) => {
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

const handleToggleLeftPanel = () => {
  console.log("Toggling left panel", showLeftPanel.value)
  showLeftPanel.value = !showLeftPanel.value
}

function handleContactClicked(jid, statusMessage) {
  console.log("Contact clicked", jid)

  Message.jid = jid  // Set the current correspondent on the frontend
  Message.isConference = false
  Message.body = ""
  Message.statusMessage = statusMessage

  Debug.resultText = "Setting correspondent to " + jid + " " + Message.statusMessage

  getMessages()  // Get the messages for the current correspondent
}

function handleConferenceClicked(jid, alias) {
  console.log("Conference clicked", jid)

  Message.jid = jid  // Set the current correspondent on the frontend
  Message.statusMessage = alias
  Message.isConference = true
  Debug.resultText = "Setting correspondent to " + jid

  getConferenceMessages(jid)  // Get the messages for the current correspondent
}

function handleFileUploaded(url) {
  console.log('File uploaded:', url);

  console.log("Sending file", Message.body)
  // file, to, from
  if (Message.isConference) {
    SendConferenceFileMessage(url, Message.jid)
  } else {
    SendFileMessage(url, Message.jid, User.jid)
  }
}

// *************************************** Event listeners ***************************************

const listenMessages = async () => {
  EventsOn("message", (from) => {
    console.log("Message", from)

    Debug.resultText = "Message from " + from

    if (from === Message.jid) {
      console.log("Updating current conversation")
      // getArchive(from)

      if (Message.isConference){
        getConferenceMessages(from)
      } else {
        getMessages()
      }
    }
  })
}

const listenSuccess = async () => {
  EventsOn("success", (message) => {
    Debug.resultText = message

    if (message === "Message sent") {
      getMessages()
    } else if (message === "Subscription accepted"){
      getContacts()
    }
  })
}

const listenSubRequest = async () => {
  EventsOn("subscription-request", (user) => {
    Debug.resultText = "Subscription request from " + user
    AcceptSubscription(user)
  })
}

const listenUpdateMessages = async () => {
  EventsOn("update-messages", (jid) => {
    console.log("Updating messages for", jid)
    getMessages()
  })
}

const listenLogout = async () => {
  EventsOn("logout", () => {
    console.log("Logout")
    Debug.resultText = "Logout"
    router.push("/")
  })
}

GetCurrentUser().then((user) => {
  User.jid = user
  console.log("User", user)
  Debug.resultText = "User: " + user
})

listenMessages()
listenSuccess()
listenSubRequest()
listenUpdateMessages()
listenLogout()

getContacts()


onMounted(() => {
  scrollToBottom()
})

</script>

<template>
  <main class="flex flex-col items-center justify-start h-full">

    <Nav />

    <div id="display" class="flex justify-between w-full h-[90%]">

      <Bars3Icon class="w-12 h-12 m-3 ml-4 cursor-pointer" @click="handleToggleLeftPanel" />

      <LeftPanel v-if="showLeftPanel" @contactClicked="handleContactClicked" @conferenceClicked="handleConferenceClicked" @closePanel="handleToggleLeftPanel" />

      <div id="current-chat" class="current-chat w-11/12 mx-8 mt-2">

        <div id="top-bar" class="top-bar flex items-center justify-center h-fit my-4">
          <div id="current-contact" class="current-contact flex flex-col items-center justify-center h-full w-full p-4 bg-white rounded-xl">
            <p class="current-contact-jid text-lg text-gray-900">{{ Message.jid }}</p>
            <p v-if="Message.statusMessage" class="current-contact-status-message text-sm text-gray-500">{{ Message.statusMessage }}</p>
          </div>
          <CogIcon class="dots w-11 h-11 ml-4 cursor-pointer" @click="toggleOptions" />
          <Options :is-conference="Message.isConference" :jid="Message.jid" v-if="showOptions" @close-options="toggleOptions" />
        </div>

        <div id="messages" class="message-section h-3/4 my-8 border-2 rounded-md border-black overflow-y-auto scrollbar-thin scrollbar-thumb-black scrollbar-track-gray-600" ref="messageSectionRef">
          <Conversation :messages="Messages.messages" :user="User.jid" :is-conference="Message.isConference"/>
        </div>

        <div v-if="Message.jid" id="message-input" class="message-input flex items-center justify-center mt-8">
          <input id="message" v-model="Message.body" autocomplete="off" class="input w-4/5 h-8 px-2 rounded-md border-none bg-gray-200 focus:bg-white text-black" type="text"/>
          <PaperAirplaneIcon @click="sendMessage" class="w-16 h-8 ml-4 p-1 rounded-md cursor-pointer bg-blue-500" />
          <FileInput  @fileUploaded="handleFileUploaded" />
        </div>
      </div>
    </div>

  </main>
</template>

<style scoped>

.message-section {
  scrollbar-width: thin;
  scrollbar-color: black gray;
}
</style>
