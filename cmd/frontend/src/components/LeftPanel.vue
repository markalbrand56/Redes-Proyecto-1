<script setup>
import {
  GetConferences,
  GetContacts,
  GetCurrentUser,
  ProbeContacts,
  SetStatus,
  UpdateContacts,
  SetStatusMessage
} from '../../wailsjs/go/main/App.js'

import Contact from "./Contact.vue";
import StatusPopup from "./StatusPopUp.vue";

import {computed, onMounted, reactive, ref} from "vue";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";


const emit = defineEmits(["contactClicked", "conferenceClicked", "closePanel"]);
const showPopup = ref(false);

const User = reactive({
  jid: "",
  contacts: [],
  conferences: [],
  status: 0,
});



function togglePopup() {
  console.log("Toggling popup", showPopup.value);
  showPopup.value = !showPopup.value;
}

const getCurrentUser = async () => {
  User.jid = await GetCurrentUser();
  console.log("Current user: ", User.jid);
};

const getContacts = async () => {
  const contacts = await GetContacts();

  console.log("Contacts: ", contacts);

  User.contacts = [];
  contacts.forEach(contact => {
    User.contacts.push({
      jid: contact,
      status: "Disconnected",
      statusMessage: ""
    });
  });

  console.log("User contacts: ", User.contacts);
  ProbeContacts();
};

const getConferences = async () => {
  const conferences = await GetConferences();

  console.log("Conferences: ", conferences);

  User.conferences = [];
  conferences.forEach(conference => {
    User.conferences.push({
      jid: conference.jid,
      name: conference.alias
    });
  });

  console.log("User conferences: ", User.conferences);
};


function handleContactClicked(jid) {
  console.log("Contact clicked: ", jid);

  const statusMessage = User.contacts.find(contact => contact.jid === jid).statusMessage;
  console.log("Status message: ", statusMessage);

  emit("contactClicked", jid, statusMessage);
}

function handleConferenceClicked(jid) {
  console.log("Conference clicked: ", jid);

  const conferenceAlias = User.conferences.find(conference => conference.jid === jid).name;

  emit("conferenceClicked", jid, conferenceAlias);
}

function handleStatusChange(status) {
  console.log("Status changed: ", status);
  User.status = status;

  SetStatus(status);
}

function handleStatusMessageChange(statusMessage) {
  console.log("Status message changed: ", statusMessage);
  User.statusMessage = statusMessage;

  SetStatusMessage(statusMessage);

}

function handleClose() {
  emit("closePanel");
}

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

// *************************************** Event listeners ***************************************
const listenPresenceUpdate = async () => {
  EventsOn("presence", (jid, status, statusMessage) => {
    console.log("Presence update", jid, status, statusMessage)

    User.contacts.forEach((contact) => {
      if (contact.jid === jid) {
        contact.status = status
        contact.statusMessage = statusMessage
      }
    })
  })
}

const listenContacts = async () => {
  EventsOn("contacts", (contacts) => {
    console.log("Contacts", contacts)
    User.contacts = contacts.map(contact => {
      return {
        jid: contact,
        status: "Disconnected",
        statusMessage: ""
      }
    })

    ProbeContacts()
  })
}

const listenConferences = async () => {
  EventsOn("conferences", (conferences) => {
    // conferences is a map conferences[item.Name] = item.Jid
    console.log("Conferences", conferences)

    conferences.forEach(conference => {
      User.conferences.push({
        jid: conference.jid,
        name: conference.alias
      })
    })

  })
}

onMounted(() => {
  // Initial data fetch
  getCurrentUser();
  getContacts();
  getConferences();

  // Event listeners
  listenPresenceUpdate();
  listenContacts();
  listenConferences();
});

</script>

<template>

  <div id="left-panel" class="left-panel fixed top-12 left-16 flex flex-col items-center justify-start min-w-96 w-fit h-[calc(75%-2rem)] m-4 p-4 bg-gray-100 rounded-xl">
    <div id="correspondents" class="correspondents flex flex-col items-center justify-start w-[calc(100%-2rem)] h-4/5 mt-4 overflow-y-auto object-contain scrollbar-thin scrollbar-thumb-black scrollbar-track-gray-600">
      <h2 @click="UpdateContacts" class="cursor-pointer text-lg text-center text-gray-900 border-b border-gray-900 py-2 px-4 bg-gray-100 rounded-xl">Contacts</h2>
      <div @click.stop="handleClose" id="contacts" class="contact-section w-full">
        <Contact v-for="contact in User.contacts" :contact="{jid: contact.jid}" :key="contact" @setCorrespondent="handleContactClicked" :status="contact.status"/>
      </div>
      <h2 @click="UpdateContacts" class="cursor-default text-lg text-center text-gray-900 border-b border-gray-900 py-2 px-4 bg-gray-100 rounded-xl">Group chats</h2>
      <div @click.stop="handleClose" id="conferences" class="contact-section w-full">
        <Contact v-for="(conference) in User.conferences" :contact="{jid: conference.jid}" :alias="conference.name" :key="conference.jid" @setCorrespondent="handleConferenceClicked"/>
      </div>
    </div>
    <div id="current-account" class="current-account flex flex-col items-center justify-center w-[calc(100%-2rem)] h-fit mt-4 p-4 border-2 border-gray-300 bg-white rounded-xl cursor-pointer" @click="togglePopup">
      <div class="flex items-center justify-center">
        <div :class="['status-indicator', statusColor, 'w-3 h-3 mr-2 border border-black rounded-full']"></div>
        <p class="text-lg text-gray-900">{{ User.jid }}</p>
      </div>
      <p class="text-sm text-gray-500" v-if="User.statusMessage !== ''">{{ User.statusMessage }}</p>
    </div>
    <StatusPopup v-if="showPopup" @statusChanged="handleStatusChange" @closePopup="togglePopup" @statusMessageChanged="handleStatusMessageChange" />
  </div>

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