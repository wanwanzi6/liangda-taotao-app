// app.js
App({
  onLaunch() {
    // 1. 日志逻辑
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 2. 登录逻辑(等写完后端登录接口再联调)
    wx.login({
      success: res => {
        console.log("微信登录凭证 Code:", res.code)
      }
    })
  },

  // 3. 全局数据存储
  globalData: {
    userInfo: null,
    baseUrl: 'http://127.0.0.1:8080/api/v1' // 后端地址
  }
})
