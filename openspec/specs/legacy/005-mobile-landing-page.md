# Spec: 移动端落地页

**Module**: Mobile Landing Page  
**Priority**: P0  
**Status**: ✅ Approved  
**Related Proposal**: [../changes/001-dmh-mvp-core-features](../changes/001-dmh-mvp-core-features.md)

---

## 📋 模块概述

移动端落地页是用户参与营销活动的入口，支持多渠道访问、动态表单渲染、支付集成和专属海报生成。使用 Uni-app 框架开发，实现一套代码多端运行。

---

## 🎯 核心功能

### 1. 多渠道适配

#### 支持的渠道
- 微信浏览器（公众号/朋友圈）
- 微信小程序
- 抖音 App 内置浏览器
- 普通手机浏览器（Chrome/Safari）

#### 渠道识别
```javascript
// Uni-app 渠道识别
uni.getSystemInfo({
  success: (res) => {
    const platform = res.platform; // ios | android
    const env = uni.getEnv(); // weixin | alipay | douyin | h5
    
    // 根据环境调整功能
    if (env === 'weixin') {
      // 微信环境：支持分享、授权
    } else if (env === 'douyin') {
      // 抖音环境：调整分享接口
    }
  }
});
```

### 2. 来源追踪

#### URL 参数设计
```
https://h5.dmh.com/campaign/1?c_id=100&u_id=200

参数说明：
- campaign/1: 活动ID
- c_id: 渠道ID（channel_id）
- u_id: 推荐人ID（user_id）
```

#### 参数存储
```javascript
// 页面加载时解析URL参数
onLoad(options) {
  const campaignId = options.campaignId || this.$route.params.id;
  const channelId = options.c_id || 0;
  const referrerId = options.u_id || 0;
  
  // 存储到本地
  uni.setStorageSync('channel_id', channelId);
  uni.setStorageSync('referrer_id', referrerId);
  
  // 整个会话期间有效
  this.loadCampaign(campaignId);
}
```

### 3. 动态表单渲染

#### 表单配置加载
```javascript
async loadCampaignForm() {
  const res = await api.getCampaign(this.campaignId);
  this.campaign = res.data;
  this.formFields = res.data.formFields;
  
  // 初始化表单数据
  this.formData = {};
  this.formFields.forEach(field => {
    this.formData[field.name] = '';
  });
}
```

#### 动态表单渲染
```vue
<template>
  <view class="form-container">
    <view v-for="field in formFields" :key="field.name" class="form-item">
      <!-- 文本框 -->
      <input 
        v-if="field.type === 'text'"
        v-model="formData[field.name]"
        :placeholder="field.placeholder"
        :required="field.required"
      />
      
      <!-- 手机号 -->
      <input 
        v-else-if="field.type === 'phone'"
        v-model="formData[field.name]"
        type="number"
        maxlength="11"
        placeholder="请输入手机号"
        :required="field.required"
      />
      
      <!-- 下拉选择 -->
      <picker 
        v-else-if="field.type === 'select'"
        :value="formData[field.name]"
        :range="field.options"
        @change="onPickerChange(field.name, $event)"
      >
        <view class="picker">
          {{ formData[field.name] || field.placeholder }}
        </view>
      </picker>
    </view>
  </view>
</template>
```

#### 表单校验
```javascript
validateForm() {
  for (const field of this.formFields) {
    const value = this.formData[field.name];
    
    // 必填校验
    if (field.required && !value) {
      uni.showToast({
        title: `请填写${field.label}`,
        icon: 'none'
      });
      return false;
    }
    
    // 手机号校验
    if (field.type === 'phone' && !this.validatePhone(value)) {
      uni.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      });
      return false;
    }
  }
  
  return true;
}

validatePhone(phone) {
  return /^1[3-9]\d{9}$/.test(phone);
}
```

### 4. 微信授权

#### 静默授权（OpenID）
```javascript
// 微信公众号环境
async getWeChatAuth() {
  // 1. 检查是否已授权
  const code = this.$route.query.code;
  if (code) {
    // 2. 用code换取openid
    const res = await api.wechatAuth({ code });
    uni.setStorageSync('openid', res.data.openid);
    return res.data.openid;
  }
  
  // 3. 跳转微信授权页
  const appid = 'wx1234567890';
  const redirectUri = encodeURIComponent(window.location.href);
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appid}&redirect_uri=${redirectUri}&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect`;
  window.location.href = authUrl;
}
```

#### 手机号授权（小程序）
```vue
<button open-type="getPhoneNumber" @getphonenumber="getPhoneNumber">
  一键获取手机号
</button>

<script>
export default {
  methods: {
    async getPhoneNumber(e) {
      if (e.detail.errMsg === 'getPhoneNumber:ok') {
        const res = await api.decryptPhone({
          encryptedData: e.detail.encryptedData,
          iv: e.detail.iv
        });
        this.formData.phone = res.data.phone;
      }
    }
  }
}
</script>
```

### 5. 专属海报生成

#### 海报生成请求
```javascript
async generatePoster() {
  uni.showLoading({ title: '生成中...' });
  
  const res = await api.generatePoster({
    campaignId: this.campaignId,
    userId: this.userId,
    nickname: this.userInfo.nickname,
    avatar: this.userInfo.avatar
  });
  
  this.posterUrl = res.data.posterUrl;
  uni.hideLoading();
  
  // 显示海报预览
  this.showPoster = true;
}
```

#### 后端海报合成（Go）
```go
import "github.com/fogleman/gg"

func GeneratePoster(req *PosterRequest) (string, error) {
    // 1. 加载背景图
    bg, err := gg.LoadImage("template/poster-bg.jpg")
    dc := gg.NewContextForImage(bg)
    
    // 2. 绘制活动标题
    dc.SetRGB(0, 0, 0)
    dc.LoadFontFace("fonts/msyh.ttf", 36)
    dc.DrawStringAnchored(req.CampaignName, 375, 200, 0.5, 0.5)
    
    // 3. 生成推荐二维码
    qrCode := generateQRCode(fmt.Sprintf(
        "https://h5.dmh.com/campaign/%d?u_id=%d",
        req.CampaignId, req.UserId
    ))
    qrImg, _ := gg.LoadImage(qrCode)
    dc.DrawImage(qrImg, 275, 800)
    
    // 4. 绘制用户头像
    avatar, _ := loadImageFromURL(req.Avatar)
    dc.DrawCircle(375, 650, 50)
    dc.Clip()
    dc.DrawImage(avatar, 325, 600)
    
    // 5. 保存文件
    filename := fmt.Sprintf("posters/%s.jpg", uuid.New())
    dc.SavePNG(filename)
    
    // 6. 上传到OSS
    ossUrl := uploadToOSS(filename)
    
    return ossUrl, nil
}
```

### 6. 支付集成

#### 微信支付流程
```javascript
async handleSubmit() {
  // 1. 表单校验
  if (!this.validateForm()) return;
  
  // 2. 创建订单
  const orderRes = await api.createOrder({
    campaignId: this.campaignId,
    phone: this.formData.phone,
    formData: this.formData,
    referrerId: uni.getStorageSync('referrer_id')
  });
  
  const orderId = orderRes.data.id;
  
  // 3. 发起支付
  const payRes = await api.createPayment({
    orderId: orderId,
    payType: 'wechat',
    clientType: this.getClientType()
  });
  
  // 4. 调起支付
  if (this.isWeChatMiniProgram()) {
    // 小程序支付
    wx.requestPayment({
      ...payRes.data.payParams,
      success: () => {
        this.onPaySuccess(orderId);
      },
      fail: (err) => {
        this.onPayFail(err);
      }
    });
  } else {
    // H5支付
    window.location.href = payRes.data.mwebUrl;
  }
}
```

#### 支付结果处理
```javascript
onPaySuccess(orderId) {
  uni.showToast({
    title: '支付成功',
    icon: 'success'
  });
  
  // 跳转到成功页
  uni.navigateTo({
    url: `/pages/success/index?orderId=${orderId}`
  });
}

// 支付成功页轮询订单状态
async checkOrderStatus() {
  const timer = setInterval(async () => {
    const res = await api.getOrder(this.orderId);
    if (res.data.payStatus === 'paid') {
      clearInterval(timer);
      this.orderPaid = true;
      // 显示报名成功码
      this.showSuccessInfo();
    }
  }, 2000);
}
```

---

## 🎨 页面设计

### 页面列表
```
/pages/campaign/index     - 活动详情页
/pages/campaign/form      - 报名表单页
/pages/success/index      - 支付成功页
/pages/poster/index       - 海报页
/pages/my/orders          - 我的订单
/pages/my/rewards         - 我的奖励
```

### 活动详情页
```vue
<template>
  <view class="campaign-page">
    <!-- 活动主图 -->
    <image :src="campaign.mainImage" mode="widthFix" />
    
    <!-- 活动信息 -->
    <view class="campaign-info">
      <text class="title">{{ campaign.name }}</text>
      <text class="desc">{{ campaign.description }}</text>
      <view class="time">
        <text>活动时间：{{ formatTime(campaign.startTime) }} - {{ formatTime(campaign.endTime) }}</text>
      </view>
    </view>
    
    <!-- 报名按钮 -->
    <view class="action-bar">
      <button class="btn-primary" @click="goToForm">立即报名</button>
      <button class="btn-secondary" @click="sharePoster">分享赚奖励</button>
    </view>
  </view>
</template>
```

---

## 🔌 API 接口调用

### API 封装
```javascript
// api/campaign.js
export default {
  // 获取活动详情
  getCampaign(id) {
    return request.get(`/api/v1/campaigns/${id}`);
  },
  
  // 创建订单
  createOrder(data) {
    return request.post('/api/v1/orders', data);
  },
  
  // 发起支付
  createPayment(data) {
    return request.post('/api/v1/orders/payment', data);
  },
  
  // 生成海报
  generatePoster(data) {
    return request.post('/api/v1/qrcode/generate', data);
  },
  
  // 微信授权
  wechatAuth(data) {
    return request.post('/api/v1/auth/wechat', data);
  }
};
```

---

## ✅ 验收标准

### 功能验收
- [ ] 多渠道正常访问
- [ ] 来源参数正确追踪
- [ ] 动态表单正确渲染
- [ ] 表单校验生效
- [ ] 微信授权正常
- [ ] 海报生成正常
- [ ] 支付流程完整

### 用户体验验收
- [ ] 页面加载 < 3 秒
- [ ] 支付体验流畅
- [ ] 错误提示友好
- [ ] 响应式适配

---

## 📝 开发清单

### 前端开发
- [ ] 初始化 Uni-app 项目
- [ ] 创建活动详情页
- [ ] 创建报名表单页
- [ ] 实现动态表单渲染
- [ ] 实现表单校验
- [ ] 集成微信授权
- [ ] 集成微信支付
- [ ] 创建支付成功页
- [ ] 创建海报页
- [ ] 创建我的订单页
- [ ] 创建我的奖励页
- [ ] 页面联调测试

### 后端开发
- [ ] 实现海报生成接口
- [ ] 实现微信授权接口
- [ ] 实现支付接口

---

## 🔗 相关文档
- [Proposal: DMH MVP 核心功能](../changes/001-dmh-mvp-core-features.md)
- [Spec: 订单与支付系统](./002-order-payment-system.md)
