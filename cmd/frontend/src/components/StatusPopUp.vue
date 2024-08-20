<script setup>

const emit = defineEmits(['statusChanged', 'closePopup']);

const updateStatus = (status) => {
  emit('statusChanged', status);
  closePopup();
};

const closePopup = () => {
  emit('closePopup');
};

</script>

<template>
  <div class="popup-overlay" @click="closePopup">
    <div class="popup-content" @click.stop>
      <button class="btn" @click="updateStatus(0)">Online</button>
      <button class="btn" @click="updateStatus(1)">Away</button>
      <button class="btn" @click="updateStatus(2)">Busy</button>
      <button class="btn" @click="updateStatus(3)">Not Available</button>
      <button class="btn" @click="updateStatus(4)">Offline</button>
    </div>
  </div>
</template>

<style scoped>
.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.popup-content {
  width: fit-content;

  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.btn {
  display: block;
  width: 100%;

  margin-bottom: 1rem;
  padding: 1rem 2rem;

  border: none;
  border-radius: 4px;
  background-color: #007bff;
  color: white;
  font-size: 16px;
  cursor: pointer;
}
</style>