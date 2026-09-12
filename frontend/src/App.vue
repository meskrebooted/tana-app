<script setup>
import { computed, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from './store/auth'
import client from './api/client'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const showNav = computed(() => auth.isAuthenticated)
const streakDays = ref(0)

async function loadStreak() {
  if (!auth.isAuthenticated) return
  try {
    const { data } = await client.get('/bookings/streak')
    streakDays.value = data.streak_days
  } catch {
    streakDays.value = 0
  }
}

watch(() => auth.isAuthenticated, loadStreak, { immediate: true })
watch(() => route.name, loadStreak)

function handleLogout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="shell">
    <aside v-if="showNav" class="sidebar">
      <div class="brand">
        <span class="brand-mark">T</span>
        <div class="brand-text">
          <h1>Tana</h1>
          <span>co-working</span>
        </div>
      </div>

      <nav>
        <router-link to="/spazi" :class="{ active: route.name === 'spaces' }">Spazi</router-link>
        <router-link to="/prenotazioni" :class="{ active: route.name === 'bookings' }">
          Le mie prenotazioni
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div v-if="streakDays > 0" class="streak">
          🔥 {{ streakDays }} {{ streakDays === 1 ? 'giorno' : 'giorni' }} di fila in Tana
        </div>
        <div class="who">
          <strong>{{ auth.user?.name }}</strong>
          <span>{{ auth.user?.email }}</span>
        </div>
        <button class="btn btn-quiet" @click="handleLogout">Esci</button>
      </div>
    </aside>

    <main :class="{ 'main-full': !showNav }">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.shell {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: var(--surface);
  border-right: 1px solid var(--hairline);
  padding: 1.75rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.brand-mark {
  width: 34px;
  height: 34px;
  border: 1px solid var(--brass);
  color: var(--brass);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 1.1rem;
  border-radius: 3px;
}

.brand-text h1 {
  font-size: 1.15rem;
  color: var(--parchment);
}

.brand-text span {
  font-size: 0.75rem;
  color: var(--muted);
}

nav {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

nav a {
  text-decoration: none;
  color: var(--muted);
  padding: 0.6rem 0.75rem;
  border-left: 2px solid transparent;
  font-size: 0.92rem;
}

nav a:hover {
  color: var(--parchment);
}

nav a.active {
  color: var(--brass);
  border-left-color: var(--brass);
  background: var(--surface-raised);
}

.sidebar-footer {
  border-top: 1px solid var(--hairline);
  padding-top: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.streak {
  font-size: 0.82rem;
  color: var(--brass);
}

.who {
  display: flex;
  flex-direction: column;
  font-size: 0.85rem;
}

.who span {
  color: var(--muted);
  font-size: 0.78rem;
}

main {
  flex: 1;
  padding: 3rem 3.5rem;
  max-width: 1000px;
}

.main-full {
  max-width: none;
  padding: 0;
}

@media (max-width: 720px) {
  .shell {
    flex-direction: column;
  }
  .sidebar {
    width: auto;
    flex-direction: row;
    align-items: center;
    flex-wrap: wrap;
    gap: 1rem;
  }
  nav {
    flex-direction: row;
  }
  main {
    padding: 2rem 1.25rem;
  }
}
</style>
