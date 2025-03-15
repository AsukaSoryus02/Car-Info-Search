<template>
  <div class="home">
    <h1>欢迎使用车辆反馈系统</h1>
    <p>本系统用于记录和管理车辆信息，包括车型、油耗、燃油类型等多种数据并给出指导建议。</p>
    <div class="action-buttons">
      <a-button type="primary" size="large" @click="$router.push('/add')">
        <template #icon><plus-outlined /></template>
        添加车辆信息
      </a-button>
      <a-button size="large" @click="$router.push('/list')" style="margin-left: 16px">
        <template #icon><unordered-list-outlined /></template>
        查看车辆列表
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { PlusOutlined, UnorderedListOutlined } from '@ant-design/icons-vue';
</script>

<style scoped>
.home {
  text-align: center;
  padding: 20px;
}

h1 {
  margin-bottom: 20px;
}

p {
  font-size: 16px;
  margin-bottom: 40px;
}

.action-buttons {
  margin-top: 30px;
}
</style>