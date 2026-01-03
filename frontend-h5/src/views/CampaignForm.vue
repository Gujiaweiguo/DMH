<template>
  <div class="container">
    <div class="card">
      <h2 class="title">报名信息</h2>
      
      <div class="field">
        <label class="label">手机号 *</label>
        <input 
          v-model="form.phone" 
          type="tel" 
          class="input" 
          placeholder="请输入手机号"
          maxlength="11"
        />
      </div>

      <div v-for="field in formFields" :key="field" class="field">
        <label class="label">{{ field }} *</label>
        <input 
          v-model="formData[field]" 
          type="text" 
          class="input" 
          :placeholder="'请输入' + field"
        />
      </div>
    </div>

    <div class="tips">提交即表示同意本活动的报名规则与隐私策略。</div>

    <button 
      class="btn-submit" 
      @click="handleSubmit"
      :disabled="submitting"
    >
      {{ submitting ? '提交中...' : '提交报名' }}
    </button>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const campaignId = route.params.id
const formFields = ref([])
const form = reactive({
  phone: ''
})
const formData = reactive({})
const submitting = ref(false)
const source = ref({ c_id: '', u_id: '' })

// 读取来源信息
const loadSource = () => {
  try {
    const saved = localStorage.getItem('dmh_source')
    if (saved) {
      source.value = JSON.parse(saved)
    }
  } catch (e) {
    console.error('读取来源信息失败', e)
  }
}

// 获取活动信息
const fetchCampaign = async () => {
  try {
    const response = await fetch(`/api/v1/campaigns/detail?id=${campaignId}`)
    if (response.ok) {
      const campaign = await response.json()
      formFields.value = campaign.formFields || []
      formFields.value.forEach(field => {
        formData[field] = ''
      })
    }
  } catch (error) {
    console.error('加载活动失败', error)
  }
}

// 表单校验
const validate = () => {
  if (!form.phone) {
    alert('请输入手机号')
    return false
  }
  if (!/^1[3-9]\d{9}$/.test(form.phone)) {
    alert('请输入正确的手机号')
    return false
  }
  for (const key of formFields.value) {
    if (!formData[key] || !formData[key].trim()) {
      alert(`请填写${key}`)
      return false
    }
  }
  return true
}

// 提交报名
const handleSubmit = async () => {
  if (submitting.value) return
  if (!validate()) return

  submitting.value = true

  try {
    const payload = {
      campaignId: Number(campaignId),
      phone: form.phone,
      formData: { ...formData },
      referrerId: source.value.u_id ? Number(source.value.u_id) : 0
    }

    const response = await fetch('/api/v1/orders', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })

    if (response.ok) {
      alert('报名成功')
      setTimeout(() => {
        router.push('/success')
      }, 500)
    } else {
      const error = await response.json()
      alert(error.message || '报名失败，请重试')
    }
  } catch (error) {
    console.error('提交失败', error)
    alert('网络错误，请重试')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadSource()
  fetchCampaign()
})
</script>

<style scoped>
.container {
  padding: 16px;
  padding-bottom: 100px;
}

.card {
  background-color: #fff;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
}

.title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
}

.field {
  margin-bottom: 16px;
}

.label {
  display: block;
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.input {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  padding: 12px 16px;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
}

.input:focus {
  border-color: #4f46e5;
}

.tips {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-bottom: 16px;
  line-height: 1.5;
}

.btn-submit {
  width: 100%;
  background-color: #4f46e5;
  color: #fff;
  border-radius: 999px;
  padding: 16px 0;
  text-align: center;
  font-size: 16px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-submit:active:not(:disabled) {
  background-color: #4338ca;
}

.btn-submit:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}
</style>
