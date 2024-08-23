<script setup>
import { reactive, computed, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime.js";

import {
  AcceptSubscription,
  CancelSubscription,
  RejectSubscription,
  RequestContact,
  AcceptConferenceInvitation,
  Logout
} from '../../wailsjs/go/main/App.js';

import { PlusIcon } from "@heroicons/vue/24/solid";
import {PowerIcon} from "@heroicons/vue/24/solid";

import Swal from "sweetalert2";

// Estado para las notificaciones
const state = reactive({
  // Notifications: [
  //   { type: "subscription", message: "Subscription request from user1", username: "user1" },
  // ],
  // Errors: [{ type: "error", message: "Error message" }],
  Notifications: [],
  Errors: [],
});

// Cálculo del total de notificaciones
const totalNotifications = computed(() => {
  return state.Notifications.length + state.Errors.length;
});

// Estado para mostrar/ocultar el panel de notificaciones
const showNotificationPanel = ref(false);

// Funciones para manejar los eventos
const onSuccess = async () => {
  EventsOn("success", (message) => {
    state.Notifications.push({
      type: "success",
      message: message,
    });
  });
};

const onError = async () => {
  EventsOn("error", (message) => {
    state.Errors.push({
      type: "error",
      message: message,
    });

    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: message,
      footer: '<a href="">Try logging in again</a>'
    });
  });
};

const onSubscribe = async () => {
  EventsOn("subscription", (username) => {
    state.Notifications.push({
      type: "subscription",
      message: `${username} has requested to subscribe to you.`,
      username: username,
    });
  });
};

const onConferenceInvite = async () => {
  EventsOn("conference-invitation", (room) => {
    state.Notifications.push({
      type: "conference-invitation",
      message: `You have been invited to join ${room}`,
      username: room,  // JID de la sala
    });
  });
};

const onNotification = async () => {
  EventsOn("notification", (message, type) => {
    state.Notifications.push({
      type: type,
      message: message,
      username: ""
    });
  });
};

const onMessage = async () => {
  EventsOn("message", (from) => {
    state.Notifications.push({
      type: "notification",
      message: `New message from ${from}`,
      username: from,
    });
  });
};

// Función para descartar una notificación
const dismissNotification = (index, type) => {
  if (type === "error") {
    state.Errors.splice(index, 1);
  } else {
    state.Notifications.splice(index, 1);
  }
};

const dismissAllNotifications = () => {
  state.Notifications = [];
  state.Errors = [];
};

// Funciones para aceptar o rechazar una suscripción
const acceptSubscription = (index, username, type) => {
  // Aquí iría la lógica para aceptar la suscripción

  if (type === "subscription") {
    console.log("Accepting subscription from", username);
    AcceptSubscription(username);
    dismissNotification(index, "subscription");
  } else if (type === "conference-invitation") {
    console.log("Accepting conference invitation from", username);
    AcceptConferenceInvitation(username);
    dismissNotification(index, "conference-invitation");
  }

};

// Función para rechazar una suscripción
const rejectSubscription = (index, username, type) => {
  // Aquí iría la lógica para rechazar la suscripción
  RejectSubscription(username)
  dismissNotification(index, "subscription");
};

// Funciones para cerrar sesión
const logout = () => {
  Swal.fire({
    title: 'Are you sure you want to log out?',
    showCancelButton: true,
    confirmButtonText: `Yes`,
    cancelButtonText: `No`,
  }).then((result) => {
    if (result.isConfirmed) {
      // Aquí iría la lógica para cerrar sesión
      console.log("Logging out...");
      Logout();
    }
  });
};

// Función para mostrar/ocultar el panel de notificaciones
const toggleNotificationPanel = () => {
  showNotificationPanel.value = !showNotificationPanel.value;
};

// Request panel
const showRequestPanel = ref(false);
const newContact = ref("");  // Estado para almacenar el JID introducido

const sendSubscriptionRequest = () => {
  if (newContact.value) {
    console.log("Sending subscription request to", newContact.value);
    RequestContact(newContact.value);
    newContact.value = "";  // Limpiar el input después de enviar la solicitud
    showRequestPanel.value = false;  // Cerrar el panel después de enviar la solicitud
  }
};


// Inicializar los listeners
onSuccess();
onError();
onSubscribe();
onNotification();
onMessage();
onConferenceInvite();

</script>

<template>
  <div class="flex justify-between items-center w-full p-3 mt-2 mb-2 bg-blue-500 text-white">
    <div class="flex justify-evenly items-center">
      <div class="relative" @click="toggleNotificationPanel">
        <div class="flex justify-center items-center w-8 h-8 rounded-full bg-red-500 text-white text-[clamp(0.5rem,14px,2rem)] cursor-pointer">
          {{ totalNotifications }}
        </div>
        <div v-if="showNotificationPanel" class="absolute top-8 left-0 w-fit max-h-[calc(80vh-100px)] min-w-[300px] min-h-[100px] p-2 bg-white border border-gray-300 rounded-lg shadow-md z-10 overflow-x-hidden overflow-y-auto scrollbar-thin">
          <button class="px-2 py-1 mb-2 bg-blue-500 text-white rounded cursor-pointer" @click="dismissAllNotifications">Dismiss all</button>
          <div v-for="(notification, index) in state.Notifications" :key="index" class="flex justify-between items-center mb-2 p-2 border-b border-gray-300 cursor-default">
            <div class="flex flex-col my-1">
              <p class="mx-4 text-gray-800">{{ notification.message }}</p>
              <div v-if="notification.type === 'subscription' || notification.type === 'conference-invitation'" class="flex justify-center mt-1">
                <button class="mx-1 px-2 py-1 bg-blue-500 text-white rounded cursor-pointer" @click.stop="acceptSubscription(index, notification.username, notification.type)">Accept</button>
                <button class="mx-1 px-2 py-1 bg-red-500 text-white rounded cursor-pointer" @click.stop="rejectSubscription(index, notification.username, notification.type)">Reject</button>
              </div>
            </div>
            <button class="px-2 py-1 bg-indigo-500 text-white rounded cursor-pointer" @click.stop="dismissNotification(index, notification.type)">Dismiss</button>
          </div>
          <div v-for="(error, index) in state.Errors" :key="index" class="flex justify-between items-center mb-2 p-2 border-b border-gray-300 cursor-default">
            <p class="text-gray-800">{{ error.message }}</p>
            <button class="px-2 py-1 bg-red-500 text-white rounded cursor-pointer" @click.stop="dismissNotification(index, 'error')">Dismiss</button>
          </div>
        </div>
      </div>
    </div>
    <div class="flex justify-evenly items-center">
      <PlusIcon class="w-7 h-7 mx-4 cursor-pointer" @click="showRequestPanel = !showRequestPanel" />
      <div v-if="showRequestPanel" class="flex items-center ml-4">
        <input v-model="newContact" type="text" placeholder="Enter JID" class="px-2 py-1 border border-gray-300 rounded mr-2 text-gray-500" />
        <button class="px-2 py-1 bg-blue-500 text-white rounded cursor-pointer" @click="sendSubscriptionRequest">Send Request</button>
      </div>
      <PowerIcon class="w-7 h-7 mx-4 cursor-pointer" @click="logout" />
    </div>
  </div>
</template>

<style scoped>

.scrollbar-thin {
  scrollbar-width: thin;
}

</style>
