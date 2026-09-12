<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'

const bookings = ref([])
const loading = ref(true)
const loadError = ref('')
const actionError = ref('')

function formatRange(start, end) {
  const opts = { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' }
  const s = new Date(start).toLocaleString('it-IT', opts)
  const e = new Date(end).toLocaleString('it-IT', { hour: '2-digit', minute: '2-digit' })
  return `${s} – ${e}`
}

async function loadBookings() {
  loading.value = true
  loadError.value = ''
  try {
    const { data } = await client.get('/bookings/me')
    bookings.value = data
  } catch (err) {
    loadError.value = err.response?.data?.error || 'Impossibile caricare le prenotazioni.'
  } finally {
    loading.value = false
  }
}

async function cancelBooking(id) {
  actionError.value = ''
  try {
    await client.patch(`/bookings/${id}/cancel`)
    await loadBookings()
  } catch (err) {
    actionError.value = err.response?.data?.error || 'Impossibile annullare la prenotazione.'
  }
}

onMounted(loadBookings)
</script>

<template>
  <div>
    <header class="page-header">
      <h1>Le mie prenotazioni</h1>
      <p>Lo storico e lo stato delle tue prenotazioni presso Pergola.</p>
    </header>

    <div v-if="loadError" class="error-banner">{{ loadError }}</div>
    <div v-if="actionError" class="error-banner">{{ actionError }}</div>

    <p v-if="loading" class="muted">Caricamento prenotazioni…</p>
    <p v-else-if="!bookings.length" class="muted">Non hai ancora nessuna prenotazione. Vai su "Spazi" per iniziare.</p>

    <ul v-else class="booking-list">
      <li v-for="b in bookings" :key="b.id" class="booking-row">
        <div>
          <h3>{{ b.space?.name }}</h3>
          <span class="muted">{{ formatRange(b.start_time, b.end_time) }}</span>
        </div>
        <div class="booking-actions">
          <span class="badge" :class="b.status === 'confirmed' ? 'badge-confirmed' : 'badge-cancelled'">
            {{ b.status === 'confirmed' ? 'confermata' : 'annullata' }}
          </span>
          <button
            v-if="b.status === 'confirmed'"
            class="btn btn-quiet"
            @click="cancelBooking(b.id)"
          >
            Annulla
          </button>
        </div>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.page-header {
  margin-bottom: 2.25rem;
}

.page-header h1 {
  font-size: 1.9rem;
  margin-bottom: 0.5rem;
}

.page-header p {
  color: var(--muted);
  margin: 0;
}

.muted {
  color: var(--muted);
  font-size: 0.88rem;
}

.booking-list {
  list-style: none;
  margin: 0;
  padding: 0;
  border-top: 1px solid var(--hairline);
}

.booking-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.15rem 0.25rem;
  border-bottom: 1px solid var(--hairline);
}

.booking-row h3 {
  font-size: 1.05rem;
  margin-bottom: 0.2rem;
}

.booking-actions {
  display: flex;
  align-items: center;
  gap: 0.9rem;
}
</style>
