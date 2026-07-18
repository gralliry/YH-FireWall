<style scoped>
.container {
    display: flex;
    flex-direction: column;
    height: 80vh;
}

.header {
    height: 25px;
    margin-left: 25px;
    display: flex;
}

.main {
    flex: 1;
    overflow: hidden;
    /* 不让外层滚动 */
    display: flex;
}

.textarea-with-lines {
    display: flex;
    width: 100%;
    height: 100%;
    border-radius: 6px;
    overflow: hidden;
}

/* 行号区域 */
.line-numbers {
    text-align: right;
    user-select: none;
    padding: 6px 8px;
    border-right: 1px solid;
    overflow: hidden;
    font-size: 14px;
    line-height: 20px;
}

/* 输入区域 */
.code-input {
    flex: 1;
}

/* Element Plus 内部 textarea */
.code-input :deep(.el-textarea__inner) {
    width: 100%;
    height: 100%;
    border: none;
    outline: none;
    resize: none;
    overflow-y: auto;
    box-sizing: border-box;
    background-color: transparent;
    font-size: 14px;
    line-height: 20px;
}
</style>


<template>
    <el-container class="container" v-loading="loading">
        <el-header class="header">
            <el-button :icon="Refresh" @click="handleRefresh" :disabled="loading" />
            <el-button :icon="Select" @click="handleSave" :disabled="loading" />
        </el-header>
        <el-main class="main">
            <div class="textarea-with-lines">
                <div class="line-numbers">
                    <div v-for="i in lineCount" :key="i">{{ i }}</div>
                </div>
                <el-input v-model="value" type="textarea" class="code-input" />
            </div>
        </el-main>
    </el-container>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, onActivated } from 'vue'

import { ElMessage } from 'element-plus'

import { Select, Refresh } from '@element-plus/icons-vue'

import { axiosInstance } from '@/api/instance'

const value = ref<string>('')

const lineCount = ref<number>(1)

const loading = ref(false)

// 初始化行数，同步一次
watch(value, () => { lineCount.value = value.value.split('\n').length || 1 }, { immediate: true })

// 处理刷新配置文件
function handleRefresh() {
    loading.value = true
    axiosInstance.get('/config').then(res => {
        value.value = res.data
        ElMessage.success(`Config file refreshed`)
    }).catch(err => {
        value.value = ''
        ElMessage.error(err.response?.data || err.message || "Failed to load config file")
    }).finally(() => {
        loading.value = false
    })
}

// 组件时刷新一次配置文件
onActivated(() => {
    handleRefresh()
})

// 处理保存配置文件
function handleSave() {
    loading.value = true
    axiosInstance.post('/config', value.value).then(() => {
        ElMessage.success('Config file saved')
    }).catch(err => {
        ElMessage.error(err.response?.data || err.message || "Failed to save config file")
    }).finally(() => {
        loading.value = false
    })
}

// ctrl+s 监听
function handleKeydown(e: KeyboardEvent) {
    // 检查是否按下 Ctrl+S 或 Cmd+S（Mac）
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
        e.preventDefault() // 阻止浏览器默认保存页面
        handleSave()
    }
}

onMounted(() => {
    window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeydown)
})

</script>
