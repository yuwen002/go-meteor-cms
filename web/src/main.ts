import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'

// Element Plus 组件与样式由 unplugin-vue-components + ElementPlusResolver 按需注入，
// 菜单图标在 MenuItem.vue 里按映射表精确引入，这里不再整包 use(ElementPlus)。
const app = createApp(App)

app.use(createPinia())
app.use(router)
app.mount('#app')
