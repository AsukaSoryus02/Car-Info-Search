<template>
  <div class="car-list">
    <h2>车辆信息列表</h2>
    
    <!-- 搜索和筛选 -->
    <div class="search-bar">
      <a-row :gutter="16">
        <a-col :span="8">
          <a-input-search
            v-model:value="searchText"
            placeholder="搜索车型"
            @search="handleSearch"
            style="width: 100%"
          />
        </a-col>
        <a-col :span="8">
          <a-select
            v-model:value="filterBrand"
            placeholder="按品牌筛选"
            style="width: 100%"
            allowClear
            @change="handleFilterChange"
          >
            <a-select-option v-for="brand in brands" :key="brand" :value="brand">
              {{ brand }}
            </a-select-option>
          </a-select>
        </a-col>
        <a-col :span="8">
          <a-select
            v-model:value="filterFuelType"
            placeholder="按燃油类型筛选"
            style="width: 100%"
            allowClear
            @change="handleFilterChange"
          >
            <a-select-option value="汽油">汽油</a-select-option>
            <a-select-option value="柴油">柴油</a-select-option>
            <a-select-option value="电动">电动</a-select-option>
            <a-select-option value="混合动力">混合动力</a-select-option>
            <a-select-option value="插电混动">插电混动</a-select-option>
          </a-select>
        </a-col>
      </a-row>
    </div>

    <!-- 数据表格 -->
    <a-table
      :columns="columns"
      :data-source="filteredCars"
      :loading="loading"
      :pagination="{ pageSize: 10 }"
      rowKey="id"
      style="margin-top: 20px"
    >
      <!-- 自定义列渲染 -->
      <template #bodyCell="{ column, record }">
        <!-- 使用场景列 -->
        <template v-if="column.dataIndex === 'usageScenario'">
          <span>
            <a-tag v-for="tag in record.usageScenario" :key="tag" color="blue">{{ tag }}</a-tag>
            <span v-if="!record.usageScenario || record.usageScenario.length === 0">-</span>
          </span>
        </template>
        
        <!-- 操作列 -->
        <template v-if="column.dataIndex === 'action'">
          <a-button type="link" size="small" @click="viewCarDetail(record)">查看</a-button>
          <a-button type="link" size="small" danger @click="deleteCar(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <!-- 详情抽屉 -->
    <a-drawer
      title="车辆详细信息"
      placement="right"
      :visible="drawerVisible"
      :width="500"
      @close="drawerVisible = false"
    >
      <div v-if="selectedCar">
        <a-descriptions bordered :column="1">
          <a-descriptions-item label="编号">{{ selectedCar.id }}</a-descriptions-item>
          <a-descriptions-item label="品牌">{{ selectedCar.brand }}</a-descriptions-item>
          <a-descriptions-item label="车型">{{ selectedCar.model }}</a-descriptions-item>
          <a-descriptions-item label="油耗(L/100km)">
            {{ selectedCar.fuelConsumption || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="燃油类型">
            {{ selectedCar.fuelType || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="行驶里程(km)">
            {{ selectedCar.mileage || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="年均行驶里程(km)">
            {{ selectedCar.annualMileage || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="存放环境">
            {{ selectedCar.storageEnvironment || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="使用场景">
            <a-tag v-for="tag in selectedCar.usageScenario" :key="tag" color="blue">{{ tag }}</a-tag>
            <span v-if="!selectedCar.usageScenario || selectedCar.usageScenario.length === 0">-</span>
          </a-descriptions-item>
          <a-descriptions-item label="备注">
            {{ selectedCar.remarks || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="添加时间">
            {{ formatDate(selectedCar.createdAt) }}
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { message, Modal } from 'ant-design-vue';
import carApi from '../api/carApi';
import carModelsData from '../data/carModels.json';

// 表格列定义
const columns = [
  {
    title: '编号',
    dataIndex: 'id',
    key: 'id',
    width: 80,
  },
  {
    title: '品牌',
    dataIndex: 'brand',
    key: 'brand',
    width: 100,
  },
  {
    title: '车型',
    dataIndex: 'model',
    key: 'model',
    width: 120,
  },
  {
    title: '燃油类型',
    dataIndex: 'fuelType',
    key: 'fuelType',
    width: 100,
  },
  {
    title: '油耗(L/100km)',
    dataIndex: 'fuelConsumption',
    key: 'fuelConsumption',
    width: 120,
  },
  {
    title: '使用场景',
    dataIndex: 'usageScenario',
    key: 'usageScenario',
    width: 200,
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    width: 120,
  },
];

// 状态变量
const cars = ref([]);
const loading = ref(false);
const searchText = ref('');
const filterBrand = ref(undefined);
const filterFuelType = ref(undefined);
const drawerVisible = ref(false);
const selectedCar = ref(null);
const brands = ref([]);

// 初始化数据
onMounted(() => {
  // 加载品牌数据
  brands.value = carModelsData.map(item => item.brand);
  
  // 加载车辆列表
  fetchCarList();
});

// 获取车辆列表数据
const fetchCarList = () => {
  loading.value = true;
  carApi.getAllCars()
    .then(response => {
      cars.value = response.data;
      loading.value = false;
    })
    .catch(error => {
      console.error('获取车辆列表失败:', error);
      message.error('获取车辆列表失败，请稍后重试');
      loading.value = false;
    });
};

// 筛选后的车辆列表
const filteredCars = computed(() => {
  return cars.value.filter(car => {
    // 搜索文本筛选
    const textMatch = !searchText.value || 
      car.brand.includes(searchText.value) || 
      car.model.includes(searchText.value);
    
    // 品牌筛选
    const brandMatch = !filterBrand.value || car.brand === filterBrand.value;
    
    // 燃油类型筛选
    const fuelTypeMatch = !filterFuelType.value || car.fuelType === filterFuelType.value;
    
    return textMatch && brandMatch && fuelTypeMatch;
  });
});

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在计算属性中实现
};

// 筛选变更处理
const handleFilterChange = () => {
  // 筛选逻辑已在计算属性中实现
};

// 查看车辆详情
const viewCarDetail = (record) => {
  selectedCar.value = record;
  drawerVisible.value = true;
};

// 删除车辆
const deleteCar = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除编号为 ${record.id} 的车辆信息吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: () => {
      carApi.deleteCar(record.id)
        .then(() => {
          message.success('删除成功');
          fetchCarList(); // 重新加载列表
        })
        .catch(error => {
          console.error('删除失败:', error);
          message.error('删除失败，请稍后重试');
        });
    },
  });
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
};
</script>

<style scoped>
.car-list {
  padding: 0 20px;
}

.search-bar {
  margin-bottom: 20px;
  margin-top: 20px;
}
</style>