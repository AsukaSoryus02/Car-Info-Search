import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import CarForm from '../views/CarForm.vue';
import CarList from '../views/CarList.vue';

const routes = [
  { 
    path: '/', 
    name: 'home',
    component: Home,
    meta: { title: '首页' }
  },
  { 
    path: '/add', 
    name: 'add',
    component: CarForm,
    meta: { title: '添加车辆' }
  },
  { 
    path: '/list', 
    name: 'list',
    component: CarList,
    meta: { title: '车辆列表' }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 设置页面标题
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 车辆信息记录系统` : '车辆信息记录系统';
  next();
});

export default router;