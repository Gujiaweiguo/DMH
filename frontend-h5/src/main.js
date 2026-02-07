import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 引入Vant组件库
import { 
  Button, Field, Form, Cell, CellGroup, NavBar, Icon, Dialog, Toast, Popup, 
  ActionSheet, Switch, Tabs, Tab, Loading, Empty, Tag, Tabbar, TabbarItem,
  DatePicker, TimePicker, RadioGroup, Radio
} from 'vant'
import 'vant/lib/index.css'

// 引入触摸模拟器（用于桌面端调试）
// import '@vant/touch-emulator' // Temporarily disabled - package not installed

const app = createApp(App)

// 注册Vant组件
app.use(Button)
app.use(Field)
app.use(Form)
app.use(Cell)
app.use(CellGroup)
app.use(NavBar)
app.use(Icon)
app.use(Dialog)
app.use(Toast)
app.use(Popup)
app.use(ActionSheet)
app.use(Switch)
app.use(Tabs)
app.use(Tab)
app.use(Loading)
app.use(Empty)
app.use(Tag)
app.use(Tabbar)
app.use(TabbarItem)
app.use(DatePicker)
app.use(TimePicker)
app.use(RadioGroup)
app.use(Radio)

app.use(router)
app.mount('#app')
