<script setup>
import { reactive, computed, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime.js";

import {
  AcceptSubscription,
  CancelSubscription,
  RequestContact,
  AcceptConferenceInvitation,
  Logout
} from '../../wailsjs/go/main/App.js';

import { PlusIcon } from "@heroicons/vue/24/solid";
import {PowerIcon} from "@heroicons/vue/24/solid";

import Swal from "sweetalert2";

// Estado para las notificaciones
const state = reactive({
  Notifications: [
      {type: "subscription", message: "John Doe has requested to subscribe to you.", username: "John Doe" },
  ],
  Errors: [{ type: "error", message: "Error message" }],
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
  CancelSubscription(username)
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
  <div class="nav-container">
    <div class="nav-section">
      <div class="notifications" @click="toggleNotificationPanel">
        <div class="notification-count">{{ totalNotifications }}</div>
        <div v-if="showNotificationPanel" class="notification-panel">
          <button class="notification-dismiss-all" @click="dismissAllNotifications">Dismiss all</button>
          <div v-for="(notification, index) in state.Notifications" :key="index" class="notification-item">
            <div class="notification-body">
              <p>{{ notification.message }}</p>

              <div v-if="notification.type === 'subscription' || notification.type === 'conference-invitation'" class="subscription-buttons">
                <button @click.stop="acceptSubscription(index, notification.username, notification.type)">Accept</button>
                <button @click.stop="rejectSubscription(index, notification.username, notification.type)">Reject</button>
              </div>

            </div>
            <button @click.stop="dismissNotification(index, notification.type)">Dismiss</button>
          </div>

          <div v-for="(error, index) in state.Errors" :key="index" class="notification-item">
            <p>{{ error.message }}</p>
            <button @click.stop="dismissNotification(index, 'error')">Dismiss</button>
          </div>

        </div>
      </div>
    </div>
    <div class="nav-section">
      <PlusIcon class="icon" @click="showRequestPanel = !showRequestPanel" />
      <div v-if="showRequestPanel" class="request-panel">
        <input v-model="newContact" type="text" placeholder="Enter JID" class="request-input" />
        <button @click="sendSubscriptionRequest">Send Request</button>
      </div>
      <PowerIcon class="icon" @click="logout" />
    </div>

  </div>
</template>

<style scoped>
.nav-container {
  display: flex;
  justify-content: space-between;
  align-items: center;

  width: calc(100% - 1.5rem);

  padding: 0.75rem;

  background-color: #007bff;
  color: white;
}

.nav-section {
  display: flex;
  justify-content: space-evenly;
  align-items: center;
}

.notifications {
  position: relative;
}

.notification-count {
  display: flex;
  justify-content: center;
  align-items: center;

  width: 32px;
  height: 32px;

  border-radius: 50%;
  background-color: red;
  color: white;
  font-size: clamp(0.5rem, 14px, 2rem);
  cursor: pointer;

}

.notification-panel {
  width: fit-content;
  max-height: calc(80vh - 100px);
  min-width: 300px;
  min-height: 100px;

  position: absolute;
  top: 30px;
  left: 0;

  padding: 10px;

  background-color: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;

  overflow-x: hidden;
  overflow-y: auto;
  scrollbar-width: thin;
}

.notification-dismiss-all {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #ff0000;
  color: white;
  cursor: pointer;
}

.notification-item {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;

  margin-bottom: 10px;
  padding: 10px;
  border-bottom: 1px solid #ccc;

  cursor: default;
}

.notification-body {
  display: flex;
  flex-direction: column;
  margin: 0.25rem 0;
}

.notification-body p {
  margin: 0 1rem;
  color: #1b2636;
}

.notification-item button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #ff0000;
  color: white;
  cursor: pointer;
}

.subscription-buttons {
  display: flex;
  justify-content: center;
  margin-top: 5px;
}

.subscription-buttons button {
  margin: 0 5px;
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

.icon {
  width: 28px;
  height: 28px;

  margin: 0 1rem;

  cursor: pointer;
}

.request-panel {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-left: 1rem;
}

.request-input {
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-right: 0.5rem;
}

.request-panel button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

</style>
