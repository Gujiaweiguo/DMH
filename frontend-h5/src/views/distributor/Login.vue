<template>
  <div class="distributor-login">
    <div class="login-container">
      <div class="login-header">
        <div class="logo">🎁</div>
        <h1 class="title">分销员登录</h1>
        <p class="subtitle">DMH数字营销中台</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            required
          >
        </div>

        <div class="form-group">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="请输入密码"
            required
          >
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <button
          type="submit"
          class="login-btn"
          :disabled="loading"
        >
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, computed } from "vue";
import { useRouter } from "vue-router";
import { authApi } from "../../services/brandApi.js";
import {
  getDefaultForm,
  hasDistributorRole,
  saveDistributorAuth,
  isAlreadyLoggedIn,
  getLoginErrorMessage,
  getButtonText
} from "./distributorLogin.logic.js";

const router = useRouter();

const form = reactive(getDefaultForm());

const loading = ref(false);
const errorMessage = ref("");

onMounted(() => {
  if (isAlreadyLoggedIn()) {
    router.replace("/distributor");
  }
});

// biome-ignore lint/correctness/noUnusedVariables: used in template
const handleLogin = async () => {
  loading.value = true;
  errorMessage.value = "";

  try {
    const data = await authApi.login(form.username, form.password);
    if (!data) {
      throw new Error("登录响应为空");
    }

    if (!hasDistributorRole(data)) {
      throw new Error("您没有分销员权限");
    }

    saveDistributorAuth(data);
    router.push("/distributor");
  } catch (error) {
    errorMessage.value = getLoginErrorMessage(error);
  } finally {
    loading.value = false;
  }
};

// biome-ignore lint/correctness/noUnusedVariables: used in template
const loginButtonText = computed(() => getButtonText(loading.value));
</script>

<style scoped>
.distributor-login {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  font-size: 48px;
  margin-bottom: 16px;
}

.title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.login-form {
  margin-bottom: 10px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 16px;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  font-size: 16px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  margin-bottom: 20px;
  text-align: center;
}

.login-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.3s;
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
