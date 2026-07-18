<template>
    <el-table :data="connectionData" v-loading="loading" stripe highlight-current-row>
        <el-table-column label="Protocol" align="center" width="100">
            <template #default="{ row }">
                {{ proto2str[row.protocol] }}
            </template>
        </el-table-column>
        <el-table-column label="PID" align="center" width="100" sortable>
            <template #default="{ row }">
                {{ row.process?.pid }}
            </template>
        </el-table-column>
        <el-table-column label="Exe" align="center" sortable show-overflow-tooltip>
            <template #default="{ row }">
                {{ row.process?.exe }}
            </template>
        </el-table-column>
        <el-table-column label="Name" align="center" width="120">
            <template #default="{ row }">
                {{ row.process?.name }}
            </template>
        </el-table-column>
        <el-table-column label="User" align="center" width="100">
            <template #default="{ row }">
                {{ row.process?.username }}
            </template>
        </el-table-column>
        <el-table-column label="LocalAddr" show-overflow-tooltip>
            <template #default="{ row }">
                {{ row.localIP }}
            </template>
        </el-table-column>
        <el-table-column label="" width="50">
            <template #default="{ row }">
                <span v-if="row.direction === 0"><el-icon><Back /></el-icon></span>
                <span v-if="row.direction === 1"><el-icon><Right /></el-icon></span>
                <span v-if="row.direction === 2"><el-icon><Back /><Right /></el-icon></span>
            </template>
        </el-table-column>
        <el-table-column label="RemoteAddr" show-overflow-tooltip>
            <template #default="{ row }">
                {{ row.remoteIP }}
            </template>
        </el-table-column>
        <el-table-column label="Action" width="100" align="center">
            <template #default="{ $index, row }">
                <el-popconfirm title="Close this connection?" @confirm="handleConnectionClose($index, row)">
                    <template #reference>
                        <el-button size="small" type="danger">Close</el-button>
                    </template>
                </el-popconfirm>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup lang="ts">
import { ref, onActivated } from 'vue';
import { axiosInstance } from '@/api/instance'
import { ElMessage } from 'element-plus'
import { Back, Right } from '@element-plus/icons-vue'

const proto2str: Record<number, string> = {
    6: "TCP",
    17: "UDP",
}

interface Connection {
    id: string
    protocol: number
    localIP: string
    remoteIP: string
    direction: number
    establishTime: number
    process: {
        pid: number
        exe: string
        name: string
        cmdline: string
        username: string
    } | null
}

const connectionData = ref<Connection[]>([])

const loading = ref(false)

function handleGetConnections() {
    loading.value = true
    axiosInstance.get('/connection').then(res => {
        connectionData.value = res.data
        connectionData.value.sort((a, b) => b.establishTime - a.establishTime)
        ElMessage.success('Connection list refreshed')
    }).catch(err => {
        connectionData.value = []
        ElMessage.error(err.response?.data || err.message || 'Failed to get connection list')
    }).finally(() => {
        loading.value = false
    })
}

function handleConnectionClose(index: number, row: Connection) {
    axiosInstance.delete(`/connection/${row.id}`).then(() => {
        connectionData.value.splice(index, 1)
        ElMessage.success('Connection closed')
    }).catch(err => {
        ElMessage.error(err.response?.data || err.message || 'Failed to close connection')
    })
}

onActivated(() => {
    handleGetConnections()
})
</script>
