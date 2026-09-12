<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'

const spaces = ref([])
const loading = ref(true)
const loadError = ref('')

const activeSpaceId = ref(null)
const startTime = ref('')
const endTime = ref('')
const bookingError = ref('')
const bookingSuccess = ref('')
const submitting = ref(false)

const typeLabels = { desk: 'Postazione', room: 'Sala riunioni', booth: 'Phone booth' }

const surprising = ref(false)
const surpriseError = ref('')

async function surpriseMe() {
  surpriseError.value = ''
  bookingSuccess.value = ''
  surprising.value = true
  try {
    const { data } = await client.get('/spaces/random')
    activeSpaceId.value = data.id
    bookingError.value = ''
    startTime.value = ''
    endTime.value = ''
    // porta la card in vista, se serve
    requestAnimationFrame(() => {
      document.getElementById(`space-${data.id}`)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
  } catch (err) {
    surpriseError.value = err.response?.data?.error || 'Nessuno spazio libero al momento.'
  } finally {
    surprising.value = false
  }
}

async function loadSpaces() {
  loading.value = true
  loadError.value = ''
  try {
    const { data } = await client.get('/spaces')
    spaces.value = data
  } catch (err) {
    loadError.value = err.response?.data?.error || 'Impossibile caricare gli spazi.'
  } finally {
    loading.value = false
  }
}

function openBooking(space) {
  activeSpaceId.value = activeSpaceId.value === space.id ? null : space.id
  bookingError.value = ''
  bookingSuccess.value = ''
  startTime.value = ''
  endTime.value = ''
}

async function submitBooking(spaceId) {
  bookingError.value = ''
  bookingSuccess.value = ''

  if (!startTime.value || !endTime.value) {
    bookingError.value = 'Indica inizio e fine della prenotazione.'
    return
  }

  submitting.value = true
  try {
    await client.post('/bookings', {
      space_id: spaceId,
      start_time: new Date(startTime.value).toISOString(),
      end_time: new Date(endTime.value).toISOString()
    })
    bookingSuccess.value = 'Prenotazione confermata.'
    activeSpaceId.value = null
  } catch (err) {
    bookingError.value = err.response?.data?.error || 'Impossibile completare la prenotazione.'
  } finally {
    submitting.value = false
  }
}

onMounted(loadSpaces)
</script>

<template>
  <div>
    <header class="page-header">
      <div class="page-header-top">
        <div>
          <h1>Trova il tuo angolo di oggi</h1>
          <p>Scegli una postazione, una sala riunioni o un phone booth e prenota una fascia oraria.</p>
        </div>
        <button class="btn btn-quiet" :disabled="surprising" @click="surpriseMe">
          {{ surprising ? 'Cerco…' : '🎲 Sorprendimi' }}
        </button>
      </div>
    </header>

    <div v-if="bookingSuccess" class="success-banner">{{ bookingSuccess }}</div>
    <div v-if="surpriseError" class="error-banner">{{ surpriseError }}</div>
    <div v-if="loadError" class="error-banner">{{ loadError }}</div>

    <p v-if="loading" class="muted">Caricamento spazi…</p>

    <ul v-else class="space-list">
      <li :id="`space-${space.id}`" v-for="space in spaces" :key="space.id" class="space-row">
        <div class="space-row-main" @click="openBooking(space)">
          <div>
            <h3>{{ space.name }}</h3>
            <span class="muted">{{ space.location }}</span>
            <p v-if="space.vibe" class="vibe">{{ space.vibe }}</p>
          </div>
          <div class="space-meta">
            <span class="badge">{{ typeLabels[space.type] || space.type }}</span>
            <span class="muted">{{ space.capacity }} {{ space.capacity === 1 ? 'persona' : 'persone' }}</span>
          </div>
        </div>

        <div v-if="activeSpaceId === space.id" class="booking-panel">
          <div v-if="bookingError" class="error-banner">{{ bookingError }}</div>
          <div class="booking-fields">
            <div class="field">
              <label>Inizio</label>
              <input v-model="startTime" type="datetime-local" />
            </div>
            <div class="field">
              <label>Fine</label>
              <input v-model="endTime" type="datetime-local" />
            </div>
          </div>
          <button class="btn btn-solid" :disabled="submitting" @click="submitBooking(space.id)">
            {{ submitting ? 'Prenotazione…' : 'Conferma prenotazione' }}
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

.page-header-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.page-header h1 {
  font-size: 1.9rem;
  margin-bottom: 0.5rem;
}

.page-header p {
  color: var(--muted);
  margin: 0;
}

.vibe {
  margin: 0.4rem 0 0;
  font-size: 0.85rem;
  color: var(--sage);
  max-width: 46ch;
}

.muted {
  color: var(--muted);
  font-size: 0.88rem;
}

.success-banner {
  border: 1px solid var(--sage);
  color: var(--sage);
  padding: 0.6rem 0.9rem;
  border-radius: 3px;
  font-size: 0.85rem;
  margin-bottom: 1.25rem;
}

.space-list {
  list-style: none;
  margin: 0;
  padding: 0;
  border-top: 1px solid var(--hairline);
}

.space-row {
  border-bottom: 1px solid var(--hairline);
}

.space-row-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.15rem 0.25rem;
  cursor: pointer;
}

.space-row-main h3 {
  font-size: 1.05rem;
  margin-bottom: 0.2rem;
}

.space-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.booking-panel {
  padding: 0 0.25rem 1.5rem;
}

.booking-fields {
  display: flex;
  gap: 1.25rem;
  flex-wrap: wrap;
}

.booking-fields .field {
  flex: 1;
  min-width: 200px;
}
</style>
