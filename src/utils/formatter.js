// 格式化工具函数

/**
 * 格式化日期时间
 * @param {string|Date} dateTime - 日期时间字符串或Date对象
 * @param {string} format - 格式化模式，默认为 'YYYY-MM-DD HH:mm'
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDateTime(dateTime, format = 'YYYY-MM-DD HH:mm') {
  if (!dateTime) return '';
  
  const date = typeof dateTime === 'string' ? new Date(dateTime) : dateTime;
  
  if (isNaN(date.getTime())) return '';
  
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
}

/**
 * 格式化数字，保留指定小数位
 * @param {number} num - 要格式化的数字
 * @param {number} digits - 小数位数，默认为2
 * @returns {string} 格式化后的数字字符串
 */
export function formatNumber(num, digits = 2) {
  if (num === undefined || num === null) return '';
  return Number(num).toFixed(digits);
}