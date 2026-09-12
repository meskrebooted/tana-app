import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'

import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import SpacesView from '../views/SpacesView.vue'
import BookingsView from '../views/BookingsView.vue'

const routes = [
  { path: '/', redirect: '/spazi' },
  { path: '/login', name: 'login', component: LoginView, meta: { guestOnly: true } },
  { path: '/registrati', name: 'register', component: RegisterView, meta: { guestOnly: true } },
  { path: '/spazi', name: 'spaces', component: SpacesView, meta: { requiresAuth: true } },
  { path: '/prenotazioni', name: 'bookings', component: BookingsView, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' }
  }
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'spaces' }
  }
  return true
})

export default router
