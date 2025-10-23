import { createApp } from 'vue'
import { createPinia } from 'pinia'

import './style.css'
import App from './App.vue'
import router from '@/global/router.js';


const pinia = createPinia()

createApp(App).use(pinia).use(router).mount('#app')
