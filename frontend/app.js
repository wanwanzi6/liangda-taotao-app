// app.js
const auth = require('./utils/auth.js')

App({
  onLaunch() {
    // 1. 日志逻辑
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 2. 检查登录态，如果已有 token 先尝试自动登录
    this.checkLoginAndRefresh()
  },

  // 检查登录态并刷新
  checkLoginAndRefresh() {
    if (auth.isLoggedIn()) {
      // 已登录，刷新用户信息
      this.getUserProfile()
    }
  },

  // 微信登录
  login(callback) {
    wx.login({
      success: res => {
        // 发送 code 到后端换取 token
        wx.request({
          url: `${this.globalData.baseUrl}/login`,
          method: 'POST',
          data: { code: res.code },
          success: loginRes => {
            if (loginRes.data.code === 200 && loginRes.data.data.token) {
              // 保存 token 和用户信息
              auth.setToken(loginRes.data.data.token)
              auth.setUserInfo(loginRes.data.data.user)

              this.globalData.userInfo = loginRes.data.data.user

              // 回调
              if (callback) callback(true, loginRes.data.data)
            } else {
              console.error('登录失败:', loginRes.data)
              if (callback) callback(false, loginRes.data)
            }
          },
          fail: err => {
            console.error('登录请求失败:', err)
            if (callback) callback(false, err)
          }
        })
      },
      fail: err => {
        console.error('wx.login 失败:', err)
        if (callback) callback(false, err)
      }
    })
  },

  // 获取用户信息（需要先通过 wx.login 获取 token）
  getUserProfile() {
    wx.request({
      url: `${this.globalData.baseUrl}/user`,
      method: 'GET',
      header: {
        'Authorization': 'Bearer ' + auth.getToken()
      },
      success: res => {
        if (res.data.code === 200) {
          auth.setUserInfo(res.data.data)
          this.globalData.userInfo = res.data.data
        } else if (res.data.error === 'token 无效或已过期') {
          // token 过期，重新登录
          auth.clearAuth()
          this.globalData.userInfo = null
        }
      },
      fail: () => {
        // 网络错误，静默处理
      }
    })
  },

  // 封装带 token 的请求方法
  requestWithAuth(options) {
    const token = auth.getToken()
    const header = options.header || {}

    if (token) {
      header['Authorization'] = 'Bearer ' + token
    }

    return wx.request({
      ...options,
      header,
      fail: err => {
        if (options.fail) options.fail(err)
      },
      success: res => {
        // token 过期处理
        if (res.statusCode === 401) {
          auth.clearAuth()
          this.globalData.userInfo = null
          wx.showToast({ title: '请重新登录', icon: 'none' })
        } else if (options.success) {
          options.success(res)
        }
      }
    })
  },

  // 退出登录
  logout() {
    auth.clearAuth()
    this.globalData.userInfo = null
  },

  // 全局数据存储
  globalData: {
    userInfo: null,
    baseUrl: 'http://127.0.0.1:8080/api/v1'
  }
})
