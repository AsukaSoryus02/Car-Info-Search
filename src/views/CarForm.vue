<template>
  <div class="car-form">
    <h2>添加车辆信息</h2>
    <a-form
      :model="formState"
      name="carForm"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 14 }"
      autocomplete="off"
      @finish="onFinish"
      @finishFailed="onFinishFailed"
    >
      <!-- 车型选择（必填） -->
      <a-form-item
        label="品牌"
        name="brand"
        :rules="[{ required: true, message: '请选择车辆品牌' }]"
      >
        <a-select
          v-model:value="formState.brand"
          placeholder="请选择品牌"
          @change="handleBrandChange"
        >
          <a-select-option v-for="brand in brands" :key="brand" :value="brand">
            {{ brand }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item
        label="车型"
        name="model"
        :rules="[{ required: true, message: '请选择车型' }]"
      >
        <a-select
          v-model:value="formState.model"
          placeholder="请选择车型"
          :disabled="!formState.brand"
        >
          <a-select-option v-for="model in models" :key="model" :value="model">
            {{ model }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <!-- 油耗（非必填） -->
      <a-form-item label="油耗(L/100km)" name="fuelConsumption">
        <a-input-number
          v-model:value="formState.fuelConsumption"
          :min="0"
          :max="30"
          :step="0.1"
          style="width: 100%"
        />
      </a-form-item>

      <!-- 燃油类型（非必填） -->
      <a-form-item label="燃油类型" name="fuelType">
        <a-select v-model:value="formState.fuelType" placeholder="请选择燃油类型">
          <a-select-option value="92">92</a-select-option>
          <a-select-option value="95">95</a-select-option>
        </a-select>
      </a-form-item>

      <!-- 行驶里程（非必填） -->
      <a-form-item label="行驶里程(Wkm)" name="mileage">
        <a-input-number
          v-model:value="formState.mileage"
          :min="0"
          :max="100"
          :step="0.01"
          style="width: 100%"
        />
      </a-form-item>

      <!-- 年均行驶里程（非必填） -->
      <a-form-item label="年均行驶里程(km)" name="annualMileage">
        <a-input-number
          v-model:value="formState.annualMileage"
          :min="0"
          :max="100000"
          :step="1000"
          style="width: 100%"
        />
      </a-form-item>

      <!-- 存放环境（非必填） -->
      <a-form-item label="存放环境" name="storageEnvironment">
        <a-select v-model:value="formState.storageEnvironment" placeholder="请选择存放环境">
          <a-select-option value="地下停车场">地下停车场</a-select-option>
          <a-select-option value="露天停车场">露天停车场</a-select-option>
          <a-select-option value="路边停车">路边停车</a-select-option>
        </a-select>
      </a-form-item>

      <!-- 使用场景（非必填） -->
      <a-form-item label="使用场景" name="usageScenario">
        <a-select
          v-model:value="formState.usageScenario"
          mode="multiple"
          placeholder="请选择使用场景（可多选）"
        >
          <a-select-option value="通勤">通勤</a-select-option>
          <a-select-option value="商务">商务</a-select-option>
          <a-select-option value="家用">家用</a-select-option>
          <a-select-option value="越野">越野</a-select-option>
          <a-select-option value="长途旅行">长途旅行</a-select-option>
          <a-select-option value="市区代步">市区代步</a-select-option>
        </a-select>
      </a-form-item>

      <!-- 其他信息（非必填） -->
      <a-form-item label="其他信息" name="remarks">
        <a-textarea
          v-model:value="formState.remarks"
          :rows="4"
          placeholder="请输入其他补充信息"
        />
      </a-form-item>

      <a-form-item :wrapper-col="{ offset: 6, span: 14 }">
        <a-button type="primary" html-type="submit">提交</a-button>
        <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import carApi from '../api/carApi';
import carModelsData from '../data/carModels.json';

// 表单数据
const formState = reactive({
  brand: '',
  model: '',
  fuelConsumption: undefined,
  fuelType: undefined,
  mileage: undefined,
  annualMileage: undefined,
  storageEnvironment: undefined,
  usageScenario: [],
  remarks: ''
});

// 品牌和车型数据
const brands = ref([]);
const models = ref([]);
const carModelsMap = ref({});

// 初始化品牌和车型数据
onMounted(() => {
  // 提取所有品牌
  brands.value = carModelsData.map(item => item.brand);
  
  // 创建品牌到车型的映射
  carModelsData.forEach(item => {
    carModelsMap.value[item.brand] = item.models;
  });
});

// 品牌变更时更新车型列表
const handleBrandChange = (value) => {
  formState.model = '';
  models.value = carModelsMap.value[value] || [];
};

// 表单提交成功
const onFinish = (values) => {
  // 发送数据到后端
  carApi.createCar(values)
    .then(response => {
      message.success(`车辆信息添加成功，编号: ${response.data.id}`);
      resetForm();
    })
    .catch(error => {
      console.error('提交失败:', error);
      message.error('提交失败，请稍后重试');
    });
};

// 表单提交失败
const onFinishFailed = (errorInfo) => {
  console.log('Failed:', errorInfo);
  message.error('表单验证失败，请检查输入');
};

// 重置表单
const resetForm = () => {
  Object.keys(formState).forEach(key => {
    if (key === 'usageScenario') {
      formState[key] = [];
    } else {
      formState[key] = '';
    }
  });
};
</script>

<style scoped>
.car-form {
  max-width: 800px;
  margin: 0 auto;
}
</style>