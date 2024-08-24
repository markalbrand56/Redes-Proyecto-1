<script setup>
import {reactive, ref} from "vue";
import Swal from "sweetalert2";

const emit = defineEmits(['statusChanged', 'closePopup', 'statusMessageChanged']);

const updatingStatus = ref(false);

const StatusMessage = reactive({
  message: ''
});

const updateStatus = (status) => {
  emit('statusChanged', status);
  closePopup();
};

const updateStatusMessage = () => {
  console.log('Updating status message');
  emit('statusMessageChanged', StatusMessage.message);
  closePopup();
};

const closePopup = () => {
  emit('closePopup');
};

function handleSetStatusMessage() {
  updatingStatus.value = !updatingStatus.value;
}

function handleDeleteAccount() {
  Swal.fire({
    title: 'Are you sure?',
    text: "You won't be able to revert this!",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#3085d6',
    cancelButtonColor: '#d33',
    confirmButtonText: 'Yes, delete it!'
  }).then((result) => {
    if (result.isConfirmed) {
      console.log('Deleting account');

      Swal.fire({
        title: 'Account deleted',
        text: 'Your account has been deleted',
        footer: 'You will be logged out',
        icon: 'success',
        }
      );
    }
  });
}

</script>

<template>
  <div class="flex flex-col fixed inset-0 bg-black bg-opacity-50 justify-center items-center" @click="closePopup">
    <div class="flex flex-col bg-white p-5 rounded-lg shadow-md text-center space-y-4 w-fit h-fit" @click.stop>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatus(0)">Online</button>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatus(1)">Away</button>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatus(2)">Busy</button>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatus(3)">Not Available</button>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatus(4)">Offline</button>
      <br>
      <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="handleSetStatusMessage">Set Status</button>
      <div v-if="updatingStatus" class="flex flex-col items-center justify-center">
        <input type="text" class="border border-gray-300 rounded-lg p-2 text-gray-600" placeholder="Status message" v-model="StatusMessage.message">
        <button class="py-2 px-4 bg-blue-500 text-white rounded-lg hover:bg-blue-600" @click="updateStatusMessage">Update</button>
      </div>
      <br>
      <button class="py-2 px-4 bg-red-500 text-white rounded-lg hover:bg-red-600" @click="handleDeleteAccount">Delete Account</button>
    </div>
  </div>
</template>

<style scoped>
</style>
