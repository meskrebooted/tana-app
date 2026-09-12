<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const auth = useAuthStore()
const router = useRouter()

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(name.value, email.value, password.value)
    router.push({ name: 'spaces' })
  } catch (err) {
    error.value = err.response?.data?.error || 'Registrazione non riuscita.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-screen">
    <div class="auth-card">
      <h1>Crea il tuo account</h1>
      <p class="subtitle">Unisciti a Pergola per prenotare spazi di lavoro.</p>

      <div v-if="error" class="error-banner">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <div class="field">
          <label for="name">Nome</label>
          <input id="name" v-model="name" type="text" required autocomplete="name" />
        </div>
        <div class="field">
          <label for="email">Email</label>
          <input id="email" v-model="email" type="email" required autocomplete="email" />
        </div>
        <div class="field">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            minlength="6"
            autocomplete="new-password"
          />
        </div>
        <button class="btn btn-solid" type="submit" :disabled="loading" style="width: 100%">
          {{ loading ? 'Creazione in corso…' : 'Registrati' }}
        </button>
      </form>

      <p class="switch">
        Hai già un account?
        <router-link to="/login">Accedi</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-screen {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--ink);
  padding: 1.5rem;
}

.auth-card {
  width: 100%;
  max-width: 380px;
  background: var(--surface);
  border: 1px solid var(--hairline);
  border-radius: 4px;
  padding: 2.5rem 2.25rem;
}

.auth-card h1 {
  font-size: 1.5rem;
  margin-bottom: 0.4rem;
}

.subtitle {
  color: var(--muted);
  font-size: 0.9rem;
  margin: 0 0 1.75rem;
}

.switch {
  text-align: center;
  font-size: 0.85rem;
  color: var(--muted);
  margin-top: 1.5rem;
}

.switch a {
  color: var(--brass);
  text-decoration: none;
}
</style>
