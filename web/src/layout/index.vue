<template>
  <el-container class="app-layout">
    <el-aside width="220px" class="app-aside">
      <div class="logo">Meteor CMS</div>
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header class="app-header">
        <Header />
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
</script>

<style scoped>
.app-layout {
  height: 100%;
}
.app-aside {
  background: #001529;
  color: #fff;
}
.logo {
  height: 56px;
  line-height: 56px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background: #000c17;
}
.app-header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  border-bottom: 1px solid #eee;
  background: #fff;
}
.app-main {
  background: #f0f2f5;
  padding: 16px;
}
</style>
