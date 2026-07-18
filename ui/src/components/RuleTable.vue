<style scoped lang="less">
.el-container {
    width: 100%;
    height: 100%;
}

.header {
    width: 100%;
    height: 10%;
    display: flex;
    align-items: center;
    justify-content: space-between; // ⭐ 左右分布
    padding: 0 20px; // 可选，增加左右间距
    gap: 12px;

    .header-selector {
        display: flex;
        flex-wrap: wrap; // 如果太多可以换行
        gap: 8px;
        margin-left: 20px;
    }


    .header-buttons {
        display: flex;
        gap: 10px;
    }
}

.main {
    width: 100%;
    height: 90%;
}
</style>


<template>
    <el-container v-loading="loading">
        <!-- 顶部：添加规则 & 可见列控制 -->
        <el-header class="header">
            <div class="header-selector">
                <el-checkbox v-model="visibleCols.group">Group</el-checkbox>
                <el-checkbox v-model="visibleCols.comment">Comment</el-checkbox>
                <el-checkbox v-model="visibleCols.srcNets">SrcNet</el-checkbox>
                <el-checkbox v-model="visibleCols.srcPorts">SrcPort</el-checkbox>
                <el-checkbox v-model="visibleCols.dstNets">DstNet</el-checkbox>
                <el-checkbox v-model="visibleCols.dstPorts">DstPort</el-checkbox>
                <el-checkbox v-model="visibleCols.inDevs">InDev</el-checkbox>
                <el-checkbox v-model="visibleCols.outDevs">OutDev</el-checkbox>
                <el-checkbox v-model="visibleCols.protocols">Protocol</el-checkbox>
                <el-checkbox v-model="visibleCols.accept">Accept</el-checkbox>
                <el-checkbox v-model="visibleCols.priority">Priority</el-checkbox>
                <el-checkbox v-model="visibleCols.enable">Enable</el-checkbox>
            </div>
            <div class="header-buttons">
                <el-button :icon="Refresh" @click="handleRefresh" :disabled="loading" />
                <el-button :icon="Plus" @click="handleAdd" :disabled="loading" />
            </div>
        </el-header>

        <!-- 表格主体 -->
        <el-main class="main">
            <el-table :data="ruleData" stripe highlight-current-row
                :default-sort="{ prop: 'priority', order: 'ascending' }">
                <!-- Group -->
                <el-table-column v-if="visibleCols.group" prop="group" label="Group" sortable :filters="groupData"
                    :filter-method="filterGroup">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.group" size="small" />
                        <span v-else>{{ row.data.group }}</span>
                    </template>
                </el-table-column>

                <!-- Comment -->
                <el-table-column v-if="visibleCols.comment" prop="comment" label="Comment" sortable>
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.comment" size="small" />
                        <span v-else>{{ row.data.comment }}</span>
                    </template>
                </el-table-column>

                <!-- SrcNet -->
                <el-table-column v-if="visibleCols.srcNets" prop="srcNets" label="SrcNet">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.srcNets" type="textarea" autosize />
                        <span v-else>{{ row.data.srcNets }}</span>
                    </template>
                </el-table-column>

                <!-- SrcPort -->
                <el-table-column v-if="visibleCols.srcPorts" prop="srcPorts" label="SrcPort">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.srcPorts" type="textarea" autosize />
                        <span v-else>{{ row.data.srcPorts }}</span>
                    </template>
                </el-table-column>

                <!-- DstNet -->
                <el-table-column v-if="visibleCols.dstNets" prop="dstNets" label="DstNet">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.dstNets" type="textarea" autosize />
                        <span v-else>{{ row.data.dstNets }}</span>
                    </template>
                </el-table-column>

                <!-- DstPort -->
                <el-table-column v-if="visibleCols.dstPorts" prop="dstPorts" label="DstPort">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.dstPorts" type="textarea" autosize />
                        <span v-else>{{ row.data.dstPorts }}</span>
                    </template>
                </el-table-column>

                <!-- InDev -->
                <el-table-column v-if="visibleCols.inDevs" prop="inDevs" label="InDev">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.inDevs" type="textarea" autosize />
                        <span v-else>{{ row.data.inDevs }}</span>
                    </template>
                </el-table-column>

                <!-- OutDev -->
                <el-table-column v-if="visibleCols.outDevs" prop="outDevs" label="OutDev">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.outDevs" type="textarea" autosize />
                        <span v-else>{{ row.data.outDevs }}</span>
                    </template>
                </el-table-column>

                <!-- Protocol -->
                <el-table-column v-if="visibleCols.protocols" prop="protocols" label="Protocol" align="center">
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.protocols" type="textarea" autosize />
                        <span v-else>{{ row.data.protocols }}</span>
                    </template>
                </el-table-column>

                <!-- Accept -->
                <el-table-column v-if="visibleCols.accept" label="Accept" align="center"
                    :filters="[{ text: 'Accept', value: true }, { text: 'Drop', value: false }]"
                    :filter-method="filterAccept">
                    <template #default="{ row }">
                        <el-switch v-if="row.isEditing" v-model="row.cache.accept"
                            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
                        <el-switch v-else v-model="row.data.accept" disabled
                            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
                    </template>
                </el-table-column>

                <!-- Priority -->
                <el-table-column v-if="visibleCols.priority" prop="priority" label="Priority" align="center" sortable>
                    <template #default="{ row }">
                        <el-input v-if="row.isEditing" v-model="row.cache.priority" type="number" />
                        <span v-else>{{ row.data.priority }}</span>
                    </template>
                </el-table-column>

                <!-- Enable -->
                <el-table-column v-if="visibleCols.enable" label="Enable" align="center"
                    :filters="[{ text: 'Enable', value: true }, { text: 'Disable', value: false }]"
                    :filter-method="filterEnable">
                    <template #default="{ row }">
                        <el-switch v-if="row.isEditing" v-model="row.cache.enable" />
                        <el-switch v-else v-model="row.data.enable" disabled />
                    </template>
                </el-table-column>

                <!-- Action -->
                <el-table-column label="Action" align="center" width="200">
                    <template #default="{ row, $index }">
                        <el-button v-if="!row.isEditing" size="small" type="primary" @click="handleEdit(row)">
                            Edit
                        </el-button>
                        <el-popconfirm v-if="!row.isEditing" icon-color="#626AEF" title="Delete this?" placement="left"
                            @confirm="handleDelete($index, row)">
                            <template #reference>
                                <el-button size="small" type="danger">Delete</el-button>
                            </template>
                            <template #actions="{ confirm, cancel }">
                                <el-button size="small" @click="cancel">No</el-button>
                                <el-button size="small" @click="confirm" type="danger">Yes</el-button>
                            </template>
                        </el-popconfirm>

                        <el-button v-if="row.isEditing" type="success" size="small" @click="handleConfirm(row)">
                            Confirm
                        </el-button>
                        <el-button v-if="row.isEditing" size="small" type="info" @click="handleCancel($index, row)">
                            Cancel
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-main>
    </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onActivated } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { axiosInstance } from '@/api/instance'

// ---------- 状态 ----------
interface Rule {
    id: string
    group: string
    comment: string
    srcNets: string
    srcPorts: string
    dstNets: string
    dstPorts: string
    inDevs: string
    outDevs: string
    protocols: string
    accept: boolean
    priority: number
    enable: boolean
}

interface ERule {
    data: Rule
    cache: Rule
    isEditing: boolean
    isNew: boolean
}

const ruleData = ref<ERule[]>([])

// ---------- 可见列 ----------
const visibleCols = ref({
    group: true,
    comment: true,
    srcNets: true,
    srcPorts: true,
    dstNets: true,
    dstPorts: true,
    inDevs: true,
    outDevs: true,
    protocols: true,
    accept: true,
    priority: true,
    enable: true,
})

// ---------- 工具函数 ----------
function filterGroup(value: string, row: ERule) { row.data.group === value }
function filterAccept(value: boolean, row: ERule) { row.data.accept === value }
function filterEnable(value: boolean, row: ERule) { row.data.enable === value }

const groupData = computed(() => {
    const groups = Array.from(new Set(ruleData.value.map(r => r.data.group).filter(Boolean)))
    return groups.map(g => ({ text: g, value: g }))
})

// ---------- 事件 ----------
function handleEdit(row: ERule) {
    Object.assign(row.cache, row.data)
    row.isEditing = true
}

function handleCancel(index: number, row: ERule) {
    if (row.isNew) {
        ruleData.value.splice(index, 1)
    } else {
        row.isEditing = false
    }
}

function handleConfirm(row: ERule) {
    if (row.isNew) {
        axiosInstance.post('/rule', row.cache).then(res => {
            row.cache.id = res.data
            Object.assign(row.data, row.cache)
            // 结束编辑
            row.isEditing = false
            row.isNew = false
            ElMessage.success('Added successfully')
        }).catch(err => {
            ElMessage.error(err.response?.data || err.message || 'Failed to add rule')
        })
    } else {
        // 修改的值才加入列表
        const data = Object.fromEntries(
            Object.entries(row.cache).filter(([key, value]) => {
                const k = key as keyof Rule
                return row.data[k] !== value
            })
        ) as Partial<Rule>
        axiosInstance.put(`/rule/${row.data.id}`, data).then(() => {
            Object.assign(row.data, row.cache)
            row.isEditing = false
            ElMessage.success('Saved successfully')
        }).catch(err => {
            ElMessage.error(err.response?.data || err.message || 'Failed to save rule')
        })
    }
}

function handleDelete(index: number, row: ERule) {
    axiosInstance.delete(`/rule/${row.data.id}`).then(() => {
        ruleData.value.splice(index, 1)
        ElMessage.success('Deleted successfully')
    }).catch(err => {
        ElMessage.error(err.response?.data || err.message || 'Failed to delete rule')
    })
}

function handleAdd() {
    const newRule: Rule = {
        id: '',
        group: '',
        comment: '',
        srcNets: '',
        srcPorts: '',
        dstNets: '',
        dstPorts: '',
        inDevs: '',
        outDevs: '',
        protocols: '',
        accept: true,
        priority: 0,
        enable: true,
    }
    const newERule: ERule = {
        data: { ...newRule },
        cache: { ...newRule },
        isEditing: true,
        isNew: true,
    }
    ruleData.value.unshift(newERule)
}

const loading = ref(false)

function handleRefresh() {
    loading.value = true
    axiosInstance.get('/rule').then(res => {
        ruleData.value = res.data.map((r: Rule) => ({
            data: { ...r },
            cache: { ...r },
            isEditing: false,
            isNew: false,
        }))
        ElMessage.success(`Rule list refreshed`)
    }).catch(err => {
        ruleData.value = []
        ElMessage.error(err.response?.data || err.message || 'Failed to get rule list')
    }).finally(() => {
        loading.value = false
    })
}

onActivated(() => {
    handleRefresh()
})

</script>
