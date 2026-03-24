// pages/post/post.js
const app = getApp();
const auth = require('../../utils/auth.js')

Page({
  data: {
    categories: [],
    selectedCategory: null,
    title: '',
    description: '',
    price: '',
    tempImageUrl: '',  // 本地预览路径
    remoteImageUrl: '', // 服务器返回的相对路径
    isUploading: false,
    isLoggedIn: false,
  },

  onLoad() {
    this.checkLogin();
  },

  // 检查登录状态
  checkLogin() {
    if (!auth.isLoggedIn()) {
      // 未登录，引导用户登录
      wx.showModal({
        title: '提示',
        content: '请先登录后再发布商品',
        confirmText: '去登录',
        success: (res) => {
          if (res.confirm) {
            app.login((success) => {
              if (success) {
                this.setData({ isLoggedIn: true });
                this.fetchCategories();
              }
            });
          } else {
            wx.switchTab({ url: '/pages/index/index' });
          }
        }
      });
    } else {
      this.setData({ isLoggedIn: true });
      this.fetchCategories();
    }
  },

  // 1. 选择图片
  chooseImage() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'], // 这里之前漏了逗号
      success: (res) => {
        const path = res.tempFiles[0].tempFilePath;
        this.setData({
          tempImageUrl: path,
          isUploading: true
        });
        this.uploadFile(path);
      }
    })
  },

  // 2. 上传文件到后端
  uploadFile(path) {
    wx.uploadFile({
      url: `${app.globalData.baseUrl}/upload`,
      filePath: path,
      name: 'image',
      success: (res) => {
        // wx.uploadFile 返回的是字符串，必须手动解析
        const resData = JSON.parse(res.data);
        // 这里根据你后端返回的 JSON 结构调整，假设返回的是 {code: 200, url: "/uploads/..."}
        if (resData.url) {
          this.setData({
            remoteImageUrl: resData.url,
            isUploading: false
          });
          wx.showToast({ title: '图片上传成功', icon: 'success' });
        } else {
          this.setData({ isUploading: false });
          wx.showToast({ title: '解析路径失败', icon: 'none' });
        }
      },
      fail: (err) => {
        console.error("上传错误", err);
        this.setData({ isUploading: false });
        wx.showToast({ title: '网络请求失败', icon: 'none' });
      }
    });
  },

  // 3. 获取分类
  fetchCategories() {
    wx.request({
      url: `${app.globalData.baseUrl}/categories`,
      success: (res) => {
        // 注意：适配你后端的返回格式，如果是 gin 直接返回数组，就用 res.data
        // 如果包裹了 code，就按你写的 res.data.data
        const list = res.data.data || res.data; 
        this.setData({ categories: list });
      }
    });
  },

  // 4. 表单输入处理
  onInputTitle(e) { this.setData({ title: e.detail.value }); },
  onInputDesc(e) { this.setData({ description: e.detail.value }); },
  onInputPrice(e) { this.setData({ price: e.detail.value }); },
  onCategoryChange(e) {
    const index = e.detail.value;
    this.setData({ selectedCategory: this.data.categories[index] });
  },

  // 5. 提交发布
  submitPost() {
    const { title, description, price, selectedCategory, remoteImageUrl, isUploading, isLoggedIn } = this.data;

    if (!isLoggedIn) {
      wx.showToast({ title: '请先登录', icon: 'none' });
      return;
    }

    if (isUploading) {
      wx.showToast({ title: '图片还在上传中...', icon: 'none' });
      return;
    }

    if (!title || !price || !selectedCategory) {
      wx.showToast({ title: '请填写完整信息', icon: 'none' });
      return;
    }

    // 使用带 token 的请求
    app.requestWithAuth({
      url: `${app.globalData.baseUrl}/products`,
      method: 'POST',
      data: {
        Title: title,
        Description: description,
        Price: parseFloat(price),
        CategoryID: selectedCategory.ID,
        ImageURL: remoteImageUrl,
        Status: 1
      },
      success: (res) => {
        // 兼容不同的后端返回 code 逻辑
        if (res.statusCode === 200 || (res.data && res.data.code === 200)) {
          wx.showToast({ title: '发布成功', icon: 'success' });
          setTimeout(() => {
            wx.switchTab({ url: '/pages/index/index' });
          }, 1500);
        } else {
          wx.showToast({ title: '发布失败: ' + (res.data.error || ''), icon: 'none' });
        }
      }
    });
  }
})