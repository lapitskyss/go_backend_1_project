import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import initAxios from './services/axios'

initAxios();

createApp(App).use(router).mount('#app');
