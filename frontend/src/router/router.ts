// 使用HashHistory  变成/#/path
import { createRouter, createWebHashHistory  } from 'vue-router'

// 引入页面组件
import RuleTable from '@/components/RuleTable.vue'
import ProcessTable from '@/components/ConnectionTable.vue'
import InterfaceTable from '@/components/InterfaceTable.vue'
import ConfigArea from '@/components/ConfigArea.vue'

const routes = [
    { path: '/', redirect: '/rule' },
    { path: '/rule', name: 'Rule', component: RuleTable },
    { path: '/connection', name: 'Connection', component: ProcessTable },
    { path: '/interface', name: 'Interface', component: InterfaceTable },
    { path: '/config', name: 'Config', component: ConfigArea },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router
