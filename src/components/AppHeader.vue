<template>
  <a-layout-header>
    <div class="logo">车辆信息分析反馈系统</div>
    <a-menu
      v-model:selectedKeys="selectedKeys"
      theme="dark"
      mode="horizontal"
      :style="{ lineHeight: '64px' }"
    >
      <a-menu-item key="home">
        <router-link to="/">首页</router-link>
      </a-menu-item>
      <a-menu-item key="add">
        <router-link to="/add">添加车辆</router-link>
      </a-menu-item>
      <a-menu-item key="list">
        <router-link to="/list">车辆列表</router-link>
      </a-menu-item>
    </a-menu>
  </a-layout-header>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const selectedKeys = ref([]);

// 根据当前路由更新选中的菜单项
watch(
  () => route.path,
  (path) => {
    if (path === '/') selectedKeys.value = ['home'];
    else if (path === '/add') selectedKeys.value = ['add'];
    else if (path === '/list') selectedKeys.value = ['list'];
  },
  { immediate: true }
);
</script>

<style scoped>
.logo {
  float: left;
  color: white;
  font-size: 18px;
  font-weight: bold;
  margin-right: 30px;
}
</style>