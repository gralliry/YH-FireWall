<style scoped lang="less">
.wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;

    display: flex;
    flex-direction: column;

    margin: 0;
    padding: 0;
    box-sizing: border-box;
    overflow: hidden;
}

.header {
    height: 40px;
    flex-shrink: 0;
    text-align: center;

    .header-content {
        display: flex;
        align-items: center;
        justify-content: center;
        /* 先让整体居中 */
        position: relative;
        /* 为右侧绝对定位做基础 */
        height: 100%;
    }

    /* 标题居中 */
    .header-title {
        margin: 0;
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
    }

    /* 右侧暗黑模式切换按钮 */
    .header-right {
        position: absolute;
        right: 10px;

        .el-button {
            border: none;
        }
    }
}

.main {
    width: 100%;
    display: flex;
    overflow: hidden;

    .main-nav {
        width: 150px;
        overflow-y: auto;
    }

    .main-content {
        width: 100%;
        overflow: auto;
    }
}

.footer {
    height: 40px;
    flex-shrink: 0;
    text-align: center;
}
</style>

<template>
    <el-container class="wrapper">
        <el-header class="header">
            <div class="header-content">
                <h3 class="header-title">YH Firewall Web Manager</h3>
                <div class="header-right">
                    <el-button :icon="isDark ? Moon : Sunny" @click="toggleDark()" />
                </div>
            </div>
        </el-header>

        <el-main class="main">
            <!-- 左侧导航菜单 -->
            <el-menu class="main-nav" router default-active="/rule" :unique-opened="true">
                <el-menu-item index="/rule">Rule</el-menu-item>
                <el-menu-item index="/connection">Connection</el-menu-item>
                <el-menu-item index="/interface">Interface</el-menu-item>
                <el-menu-item index="/config">Config</el-menu-item>
            </el-menu>

            <!-- 右侧内容区：由路由控制 -->
            <div class="main-content">
                <router-view v-slot="{ Component }">
                    <keep-alive>
                        <component :is="Component" />
                    </keep-alive>
                </router-view>
            </div>
        </el-main>
    </el-container>
</template>

<script setup lang="ts">
// 引入VueUse暗黑模式核心
import { useDark, useToggle } from '@vueuse/core'
import { Sunny, Moon } from '@element-plus/icons-vue'

// 初始化暗黑模式检测器
const isDark = useDark()
// 生成切换函数（类型安全）
const toggleDark = useToggle(isDark)
</script>
