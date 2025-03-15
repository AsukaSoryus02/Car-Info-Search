import axios from 'axios';

// 车辆API封装
const carApi = {
  // 获取所有车辆
  getAllCars() {
    return axios.get('/api/cars');
  },
  
  // 根据ID获取车辆
  getCarById(id) {
    return axios.get(`/api/cars/${id}`);
  },
  
  // 根据品牌获取车辆
  getCarsByBrand(brand) {
    return axios.get(`/api/cars/brand/${brand}`);
  },
  
  // 创建车辆
  createCar(car) {
    return axios.post('/api/cars', car);
  },
  
  // 更新车辆
  updateCar(id, car) {
    return axios.put(`/api/cars/${id}`, car);
  },
  
  // 删除车辆
  deleteCar(id) {
    return axios.delete(`/api/cars/${id}`);
  }
};

export default carApi;