// pages/detail/detail.js
const app = getApp();

Page({
  data: {
    product: null,
    loading: true
  },

  onLoad(options) {
    // options.id 就是从首页传过来的商品 ID
    this.fetchDetail(options.id);
  },

  fetchDetail(id) {
    wx.request({
      url: `${app.globalData.baseUrl}/products/${id}`,
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({ 
            product: res.data.data,
            loading: false 
          });
        }
      }
    });
  }
})