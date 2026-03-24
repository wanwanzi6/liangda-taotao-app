// utils/auth.js - token 存取封装

const TOKEN_KEY = 'liangda_token'
const USER_INFO_KEY = 'liangda_user'

// 保存 token
function setToken(token) {
  wx.setStorageSync(TOKEN_KEY, token)
}

// 获取 token
function getToken() {
  return wx.getStorageSync(TOKEN_KEY)
}

// 移除 token
function removeToken() {
  wx.removeStorageSync(TOKEN_KEY)
}

// 保存用户信息
function setUserInfo(userInfo) {
  wx.setStorageSync(USER_INFO_KEY, userInfo)
}

// 获取用户信息
function getUserInfo() {
  return wx.getStorageSync(USER_INFO_KEY)
}

// 移除用户信息
function removeUserInfo() {
  wx.removeStorageSync(USER_INFO_KEY)
}

// 检查是否已登录
function isLoggedIn() {
  return !!getToken()
}

// 清理登录态
function clearAuth() {
  removeToken()
  removeUserInfo()
}

module.exports = {
  setToken,
  getToken,
  removeToken,
  setUserInfo,
  getUserInfo,
  removeUserInfo,
  isLoggedIn,
  clearAuth
}
