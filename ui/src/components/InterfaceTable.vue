<template>
    <el-table :data="interfaceData" v-loading="loading" stripe highlight-current-row>
        <el-table-column label="Index" prop="index" align="center" />
        <el-table-column label="Name" prop="name" align="center" />
        <el-table-column label="MTU" prop="mtu" align="center" />
        <el-table-column label="Mac" prop="mac" align="center" />
        <el-table-column label="Flags">
            <template #default="{ row }">
                <div v-for="flag in row.flags" :key="flag">{{ flag }}</div>
            </template>
        </el-table-column>
        <el-table-column label="Networks">
            <template #default="{ row }">
                <div v-for="addr in row.addrs" :key="addr">{{ addr }}</div>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup lang="ts">
import { ref, onActivated } from 'vue';
import { ElMessage } from 'element-plus'
import { axiosInstance } from '@/api/instance'

interface Interface {
    index: number
    name: string
    mac: string
    mtu: number
    flags: string[]
    addrs: string[]
}

// 接口数据
const interfaceData = ref<Interface[]>([])

const loading = ref(false)

// 获取接口列表
function handleGetInterfaces() {
    loading.value = true
    axiosInstance.get('/interface').then(res => {
        interfaceData.value = res.data
        ElMessage.success('Interface list refreshed')
    }).catch(err => {
        interfaceData.value = []
        ElMessage.error(err.response?.data || err.message || 'Failed to get interface list')
    }).finally(() => {
        loading.value = false
    })
}

// 监听页面加载
onActivated(() => {
    handleGetInterfaces()
})

</script>