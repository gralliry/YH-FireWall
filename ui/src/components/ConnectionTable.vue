<template>
    <el-table :data="connectionData" v-loading="loading" stripe highlight-current-row>
        <el-table-column label="Fd" prop="fd" align="center" width="80" sortable />
        <el-table-column label="Protocol" align="center" width="100">
            <template #default="{ row }">
                {{ protocol2str[row.protocol as 6 | 17] }}({{ family2str[row.family as 2 | 10] }})
            </template>
        </el-table-column>
        <el-table-column label="PID" prop="pid" align="center" width="100" sortable />
        <el-table-column label="Exe" prop="exe" align="center" sortable show-overflow-tooltip />
        <el-table-column label="Name" prop="name" align="center" width="100" />
        <el-table-column label="User" prop="username" align="center" width="100" />
        <el-table-column label="Interface" prop="interface" align="center" width="100" />
        <el-table-column label="LocalAddr" show-overflow-tooltip>
            <template #default="{ row }">
                <span v-if="row.family === 2">{{ row.localIP }}:{{ row.localPort }}</span>
                <span v-else-if="row.family === 10">[{{ row.localIP }}]:{{ row.localPort }}</span>
            </template>
        </el-table-column>
        <el-table-column width="50">
            <template #default="{ row }">
                <span v-if="row.direction === 0"> <el-icon>
                        <Back />
                    </el-icon> </span>
                <span v-if="row.direction === 1"><el-icon>
                        <Right />
                    </el-icon></span>
                <span v-if="row.direction === 2"> <el-icon>
                        <Back />
                        <Right />
                    </el-icon></span>
            </template>
        </el-table-column>
        <el-table-column label="RemoteAddr" show-overflow-tooltip>
            <template #default="{ row }">
                <span v-if="row.family === 2">{{ row.remoteIP }}:{{ row.remotePort }}</span>
                <span v-else-if="row.family === 10">[{{ row.remoteIP }}]:{{ row.remotePort }}</span>
            </template>
        </el-table-column>
        <el-table-column label="Status" prop="status" width="120" />
        <el-table-column label="Action" width="100" align="center">
            <template #default="{ $index, row }">
                <el-button v-if="row.status == 'ESTABLISHED'" size="small" type="danger"
                    @click="handleConnectionClose($index, row)">Close</el-button>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup lang="ts">
import { ref, onActivated } from 'vue';
import { axiosInstance } from '@/api/instance'
import { ElMessage } from 'element-plus'
import { Back, Right } from '@element-plus/icons-vue'

const protocol2str = {
    6: "TCP",
    17: "UDP",
} as const;

const family2str = {
    2: "v4",
    10: "v6",
} as const;

interface Connection {
    id: number

    fd: number
    family: 2 | 10
    protocol: 6 | 17
    pid: number
    exe: string
    name: string
    username: string
    interface: string

    localIP: string
    localPort: number

    direction: 0 | 1 | 2 | 3

    remoteIP: string
    remotePort: number

    status: string

    establishedTime: number
}

const connectionData = ref<Connection[]>([])

const loading = ref(false)

function handleGetConnections() {
    loading.value = true
    axiosInstance.get('/connection').then(res => {
        connectionData.value = res.data
        connectionData.value.sort((a, b) => b.establishedTime - a.establishedTime)
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