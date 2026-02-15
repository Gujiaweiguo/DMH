<template>
  <div class="container">
    <div class="card">
      <div class="icon">✅</div>
      <h2 class="title">报名成功</h2>
      <p class="desc">感谢你的参与，我们已为你记录本次报名信息。</p>
    </div>
    <button class="btn-back" @click="goBack">返回首页</button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { goHome } from './success.logic.js'

const router = useRouter()

const goBack = () => goHome(router)
</script>

<style scoped>
.container {
  padding: 40px 20px;
  text-align: center;
}

.card {
  background-color: #fff;
  border-radius: 16px;
  padding: 40px 24px;
  margin-bottom: 24px;
}

.icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.title {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.desc {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}

.btn-back {
  width: 100%;
  max-width: 320px;
  background-color: #4f46e5;
  color: #fff;
  border-radius: 999px;
  padding: 14px 0;
  text-align: center;
  font-size: 16px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-back:active {
  background-color: #4338ca;
}
</style>
