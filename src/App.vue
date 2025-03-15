<template>
  <a-layout class="layout">
    <app-header />
    <a-layout-content style="padding: 0 50px">
      <div :style="{ background: '#fff', padding: '24px', minHeight: '280px' }">
        <router-view></router-view>
      </div>
    </a-layout-content>
    <a-layout-footer style="text-align: 'center'">
      车辆信息记录系统 ©2023
    </a-layout-footer>
  </a-layout>
</template>

<script setup>
import AppHeader from './components/AppHeader.vue';
</script>

<style>
.layout {
  min-height: 100vh;
}

.ant-layout-content {
  padding: 0 50px;
}

.content-container {
  background: #fff;
  padding: 24px;
  min-height: 280px;
  margin-top: 16px;
}
</style>