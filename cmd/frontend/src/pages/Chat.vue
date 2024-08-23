<script setup>
import {reactive, onMounted, nextTick, ref, computed} from 'vue'
import {
  SendMessage,
  SendConferenceMessage,
  UpdateContacts,
  ProbeContacts,
  AcceptSubscription,
  SetStatus,
  GetMessages,
  GetMessagesConference,
  GetCurrentUser,
} from '../../wailsjs/go/main/App.js'

import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {models} from "../../wailsjs/go/models.ts";

import Conversation from "../components/Conversation.vue";
import Contact from "../components/Contact.vue";
import StatusPopup from "../components/StatusPopUp.vue";
import Nav from "../components/Nav.vue";
import Options from "../components/Options.vue";

import {CogIcon} from "@heroicons/vue/24/solid";
import {Bars3Icon} from "@heroicons/vue/24/solid";
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
  contacts: [],
  conferences: {},
  status: 0,
  statusColor: 'green'
})

const Messages = reactive({
  messages: []
})

const showLeftPanel = ref(false)

const Debug = reactive({
  resultText: "Please enter your name below 👇",
})

const messageSectionRef = ref(null)
const showPopup = ref(false);
const showOptions = ref(false);

function scrollToBottom() {
  nextTick(() => {
    if (messageSectionRef.value) {
      messageSectionRef.value.scrollTop = messageSectionRef.value.scrollHeight
    }
  })
}

const togglePopup = () => {
  console.log('Toggling popup');
  showPopup.value = !showPopup.value;
};

const toggleOptions = () => {
  console.log('Toggling options');
  showOptions.value = !showOptions.value;
};

const handleStatusChange = (status) => {
  console.log('Nuevo estado:', status);
  updateStatus(status);
};


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
  Debug.resultText = "Getting contacts"
  UpdateContacts()
}

function updateStatus(status) {
  console.log("Updating status")
  Debug.resultText = "Updating status"
  SetStatus(status)
  User.status = status
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
  showLeftPanel.value = !showLeftPanel.value
}

function handleContactClicked(jid) {
  console.log("Contact clicked", jid)
  Message.jid = jid  // Set the current correspondent on the frontend
  Message.isConference = false
  Message.body = ""
  Message.statusMessage = User.contacts.find((contact) => contact.jid === jid).statusMessage

  Debug.resultText = "Setting correspondent to " + jid + " " + Message.statusMessage

  getMessages()  // Get the messages for the current correspondent
}

function handleConferenceClicked(jid) {
  console.log("Conference clicked", jid)

  Message.jid = jid  // Set the current correspondent on the frontend
  Message.isConference = true
  Debug.resultText = "Setting correspondent to " + jid

  getConferenceMessages(jid)  // Get the messages for the current correspondent
}

// Event listeners

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

const listenContacts = async () => {
  EventsOn("contacts", (contacts) => {
    // contacts is an array of strings
    Debug.resultText = "Contacts: " + contacts.join(", ")
    console.log("Contacts", contacts)
    // map each contact to a Contact object
    User.contacts = contacts.map((contact) => {
      return {jid: contact, status: "Disconnected"}
    })

    ProbeContacts()
  })
}

const listenConferences = async () => {
  EventsOn("conferences", (conferences) => {
    // conferences is a map conferences[item.Name] = item.Jid
    console.log("Conferences", conferences)
    Debug.resultText = "Conferences: " + Object.keys(conferences).join(", ")

    User.conferences = conferences

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

const listenPresenceUpdate = async () => {
  EventsOn("presence", (jid, status, statusMessage) => {
    console.log("Presence update", jid, status, statusMessage)

    User.contacts.forEach((contact) => {
      if (contact.jid === jid) {
        contact.status = status
        contact.statusMessage = statusMessage
      }
    })

    if (Message.jid === jid) {
      Message.statusMessage = statusMessage
    }
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

const statusColor = computed(() => {
  switch (User.status) {
    case 0:  //  Online
      return 'green'

    case 4:  //  Disconnected / Invisible
      return 'gray'

    case 1:  //  Away
      return 'yellow'

    case 2:  //  Busy
      return 'red'

    case 3:  //  Extended Away
      return 'orange'

    default:
      return 'green'
  }
})

listenMessages()
listenContacts()
listenSuccess()
listenSubRequest()
listenUpdateMessages()
listenConferences()
listenPresenceUpdate()
listenLogout()

getContacts()


onMounted(() => {
  scrollToBottom()
})

</script>

<template>
  <main class="flex flex-col items-center justify-center h-full">
    <h1 class="cursor-default my-4">Chat</h1>

    <Nav />

    <div id="display" class="flex justify-between w-full h-4/5">

      <Bars3Icon class="w-12 h-12 m-3 ml-4 cursor-pointer" @click="handleToggleLeftPanel" />

      <div id="left-panel" v-if="showLeftPanel" class="left-panel fixed top-12 left-16 flex flex-col items-center justify-start min-w-96 w-fit h-[calc(75%-2rem)] m-4 p-4 bg-gray-100 rounded-xl">
        <div id="correspondents" class="correspondents flex flex-col items-center justify-start w-[calc(100%-2rem)] h-4/5 mt-4 overflow-y-auto scrollbar-thin scrollbar-thumb-black scrollbar-track-gray-600">
          <h2 @click="getContacts" class="cursor-pointer text-lg text-center text-gray-900 border-b border-gray-900 py-2 px-4 bg-gray-100 rounded-xl">Contacts</h2>
          <div id="contacts" class="contact-section w-full">
            <Contact v-for="contact in User.contacts" :contact="{jid: contact.jid}" :key="contact" @setCorrespondent="handleContactClicked" :status="contact.status" @click="handleToggleLeftPanel"/>
          </div>
          <h2 class="cursor-default text-lg text-center text-gray-900 border-b border-gray-900 py-2 px-4 bg-gray-100 rounded-xl">Group chats</h2>
          <div id="conferences" class="contact-section w-full">
            <Contact v-for="(jid, name) in User.conferences" :contact="{jid: jid}" :alias="name" :key="jid" @setCorrespondent="handleConferenceClicked" @click="handleToggleLeftPanel"/>
          </div>
        </div>
        <div id="current-account" class="current-account flex items-center justify-center w-[calc(100%-2rem)] h-fit my-4 p-4 border-2 border-gray-300 bg-white rounded-xl cursor-pointer" @click="togglePopup">
          <div :class="['status-indicator', statusColor, 'w-3 h-3 mr-2 border border-black rounded-full']"></div>
          <p class="text-lg text-gray-900">{{ User.jid }}</p>
        </div>
        <StatusPopup v-if="showPopup" @statusChanged="handleStatusChange" @closePopup="togglePopup" />
      </div>

      <div id="current-chat" class="current-chat w-11/12 mx-8 my-2">

        <div id="top-bar" class="top-bar flex items-center justify-center h-fit my-4">
          <div id="current-contact" class="current-contact flex flex-col items-center justify-center h-full w-full p-4 bg-white rounded-xl">
            <p class="current-contact-jid text-lg text-gray-900">{{ Message.jid }}</p>
            <p v-if="Message.statusMessage" class="current-contact-status-message text-sm text-gray-500">{{ Message.statusMessage }}</p>
          </div>
          <CogIcon class="dots w-11 h-11 ml-4 cursor-pointer" @click="toggleOptions" />
          <Options :is-conference="Message.isConference" :jid="Message.jid" v-if="showOptions" @close-options="toggleOptions" />
        </div>

        <div id="messages" class="message-section h-[70%] my-8 border-2 rounded-md border-black overflow-y-auto scrollbar-thin scrollbar-thumb-black scrollbar-track-gray-600" ref="messageSectionRef">
          <Conversation :messages="Messages.messages" :user="User.jid" :is-conference="Message.isConference"/>
        </div>

        <div id="message-input" class="message-input flex items-center justify-center my-8">
          <input id="message" v-model="Message.body" autocomplete="off" class="input w-4/5 h-8 px-2 rounded-md border-none bg-gray-200 focus:bg-white text-black" type="text"/>
          <button @click="sendMessage" class="btn w-16 h-8 ml-4 rounded-md cursor-pointer bg-blue-500">Send</button>
        </div>
      </div>
    </div>

  </main>
</template>

<style scoped>
.status-indicator.green {
  background-color: green;
}

.status-indicator.red {
  background-color: red;
}

.status-indicator.yellow {
  background-color: #cebd00;
}

.status-indicator.orange {
  background-color: #e57c03;
}

.status-indicator.gray {
  background-color: gray;
}
</style>
