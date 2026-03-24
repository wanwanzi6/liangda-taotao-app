// pages/profile/profile.js
const app = getApp();
const auth = require('../../utils/auth.js');

Page({
  data: {
    userInfo: null,
    products: []
  },

  onShow() {
    // 每次显示页面时刷新数据
    this.checkLogin();
  },

  checkLogin() {
    if (!auth.isLoggedIn()) {
      // 未登录，引导用户登录
      wx.showModal({
        title: '提示',
        content: '请先登录',
        confirmText: '去登录',
        success: (res) => {
          if (res.confirm) {
            app.login((success) => {
              if (success) {
                this.loadData();
              }
            });
          }
        }
      });
    } else {
      this.loadData();
    }
  },

  loadData() {
    this.setData({
      userInfo: app.globalData.userInfo || auth.getUserInfo()
    });
    this.fetchMyProducts();
  },

  fetchMyProducts() {
    app.requestWithAuth({
      url: `${app.globalData.baseUrl}/user/products`,
      method: 'GET',
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({ products: res.data.data || [] });
        }
      }
    });
  },

  goToEdit() {
    wx.navigateTo({
      url: '/pages/profile/edit/edit'
    });
  },

  goToPost() {
    wx.switchTab({
      url: '/pages/post/post'
    });
  },

  refreshProduct(e) {
    const id = e.currentTarget.dataset.id;
    wx.showLoading({ title: '擦亮中...' });

    app.requestWithAuth({
      url: `${app.globalData.baseUrl}/products/${id}/refresh`,
      method: 'POST',
      success: (res) => {
        wx.hideLoading();
        if (res.data.code === 200) {
          wx.showToast({ title: '擦亮成功', icon: 'success' });
          this.fetchMyProducts(); // 刷新列表
        } else {
          wx.showToast({ title: res.data.msg || '擦亮失败', icon: 'none' });
        }
      },
      fail: () => {
        wx.hideLoading();
        wx.showToast({ title: '网络请求失败', icon: 'none' });
      }
    });
  },

  deleteProduct(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '确认删除',
      content: '确定要下架这个商品吗？',
      success: (res) => {
        if (res.confirm) {
          wx.showLoading({ title: '删除中...' });
          app.requestWithAuth({
            url: `${app.globalData.baseUrl}/products/${id}`,
            method: 'DELETE',
            success: (res) => {
              wx.hideLoading();
              if (res.data.code === 200) {
                wx.showToast({ title: '已下架', icon: 'success' });
                this.fetchMyProducts(); // 刷新列表
              } else {
                wx.showToast({ title: res.data.msg || '删除失败', icon: 'none' });
              }
            },
            fail: () => {
              wx.hideLoading();
              wx.showToast({ title: '网络请求失败', icon: 'none' });
            }
          });
        }
      }
    });
  }
});
