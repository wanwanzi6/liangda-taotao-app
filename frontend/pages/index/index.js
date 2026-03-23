// index.js
const app = getApp();

Page({
  data: {
    categories: [], // 存放分类列表
    products: [],  // 存放商品列表
    activeCategoryId: 0, // 0 代表“全部”
    loading: false
  },

  // 1. 页面加载生命周期
  onLoad: function() {
    this.fetchCategories();
    this.fetchProducts();
  },

  // 2. 获取后端分类数据
  fetchCategories() {
    wx.request({
      url: `${app.globalData.baseUrl}/categories`,
      method: 'GET', 
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({ categories: res.data.data });
        }
      }
    });
  },

  // 3. 点击分类的函数
  switchCategory(e) {
    const id = e.currentTarget.dataset.id || 0;
    this.setData({
      activeCategoryId: id, 
      products: []
    });
    this.fetchProducts();
  },

  // 3. 获取后端商品列表, 带上 CategoryId
  fetchProducts() {
    this.setData({ loading: true });
    wx.request({
      url: `${app.globalData.baseUrl}/products`,
      method: 'GET',
      data: {
        category_id: this.data.activeCategoryId,
        page: 1,
        pagesize: 20
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({ products: res.data.data });
        }
      },
      complete: () => {
        this.setData({ loading: false });
      }
    });
  },

  // 4. 跳转详情页
  goToDetail(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({
      url: `/pages/detail/detail?id=${id}`
    });
  },

  // 5. 下拉刷新
  onPullDownRefresh() {
    this.fetchProducts();
    wx.stopPullDownRefresh();
  }
})