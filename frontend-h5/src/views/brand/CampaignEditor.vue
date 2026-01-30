<template>
  <div class="campaign-editor">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <button @click="goBack" class="back-btn">
        <ArrowLeft :size="20" />
        <span>返回</span>
      </button>
      <h1 class="nav-title">{{ isEditMode ? '编辑活动' : '创建活动' }}</h1>
      <button @click="saveCampaign" :disabled="saving" class="save-btn">
        <Save :size="20" />
        <span>{{ saving ? '保存中...' : '保存' }}</span>
      </button>
    </div>

    <!-- 表单内容 -->
    <div class="form-container">
      <form @submit.prevent="saveCampaign" class="campaign-form">
        <!-- 基本信息 -->
        <div class="form-section">
          <h2 class="section-title">基本信息</h2>
          
          <div class="form-group">
            <label class="form-label">活动名称 *</label>
            <input
              v-model="form.name"
              type="text"
              class="form-input"
              placeholder="请输入活动名称（2-100字符）"
              maxlength="100"
              required
            >
          </div>

          <div class="form-group">
            <label class="form-label">活动描述</label>
            <textarea
              v-model="form.description"
              class="form-textarea"
              placeholder="请输入活动描述（最多500字符）"
              maxlength="500"
              rows="4"
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">开始时间 *</label>
              <input
                v-model="form.startTime"
                type="datetime-local"
                class="form-input"
                required
              >
            </div>
            <div class="form-group">
              <label class="form-label">结束时间 *</label>
              <input
                v-model="form.endTime"
                type="datetime-local"
                class="form-input"
                required
              >
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">奖励金额 *</label>
            <div class="input-with-unit">
              <input
                v-model.number="form.rewardRule"
                type="number"
                class="form-input"
                placeholder="请输入奖励金额"
                min="0.01"
                max="999.99"
                step="0.01"
                required
              >
              <span class="input-unit">元</span>
            </div>
          </div>

          <div v-if="isEditMode" class="form-group">
            <label class="form-label">活动状态</label>
            <select v-model="form.status" class="form-select">
              <option value="active">进行中</option>
              <option value="paused">已暂停</option>
              <option value="ended">已结束</option>
            </select>
          </div>
        </div>

        <!-- 支付配置 -->
        <div class="form-section">
          <h2 class="section-title">支付配置</h2>
          <p class="section-desc">设置活动支付方式和金额配置</p>

          <div class="form-group">
            <label class="form-label">是否需要支付</label>
            <select v-model="form.requirePayment" class="form-select">
              <option value="true">需要支付</option>
              <option value="false">免费活动</option>
            </select>
          </div>

          <div v-if="form.requirePayment === 'true' || form.requirePayment === true">
            <div class="form-group">
              <label class="form-label">支付类型</label>
              <select v-model="form.paymentType" class="form-select">
                <option value="deposit">订金支付</option>
                <option value="full">全款支付</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ form.paymentType === 'deposit' ? '订金金额' : '支付金额' }}
              </label>
              <div class="input-with-unit">
                <input
                  v-model.number="form.paymentAmount"
                  type="number"
                  class="form-input"
                  :placeholder="form.paymentType === 'deposit' ? '请输入订金金额' : '请输入支付金额'"
                  min="0.01"
                  max="99999.99"
                  step="0.01"
                >
                <span class="input-unit">元</span>
              </div>
            </div>

            <div class="qrcode-actions">
              <button
                type="button"
                @click="showQrcodePreview"
                :disabled="!isEditMode || (!form.requirePayment || form.requirePayment === 'false')"
                class="qrcode-btn"
              >
                <QrCode :size="18" />
                预览支付二维码
              </button>
            </div>
          </div>
        </div>

        <!-- 二维码预览模态框 -->
        <div v-if="showQrcode" class="qrcode-modal-overlay" @click="closeQrcode">
          <div class="qrcode-modal-content" @click.stop>
            <div class="qrcode-modal-header">
              <h3>支付二维码预览</h3>
              <button @click="closeQrcode" class="close-btn">
                <X :size="20" />
              </button>
            </div>
            <div class="qrcode-modal-body">
              <div v-if="qrcodeLoading" class="qrcode-loading">
                加载中...
              </div>
              <div v-else-if="qrcodeError" class="qrcode-error">
                {{ qrcodeError }}
              </div>
              <div v-else class="qrcode-display">
                <img :src="qrcodeData" alt="支付二维码" class="qrcode-image" />
                <p class="qrcode-expire">
                  有效期至：{{ qrcodeExpireAt }}
                </p>
                <button type="button" @click="refreshQrcode" class="refresh-btn">
                  <RefreshCw :size="16" />
                  刷新
                </button>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>

    <!-- 字段编辑模态框 -->
            </div>

            <!-- 二维码预览 -->
            <div v-if="showQrcode" class="qrcode-preview">
              <h3 class="qrcode-title">支付二维码</h3>
              <div class="qrcode-container">
                <div v-if="qrcodeLoading" class="qrcode-loading">
                  加载中...
                </div>
                <div v-else-if="qrcodeError" class="qrcode-error">
                  {{ qrcodeError }}
                </div>
                <div v-else class="qrcode-display">
                  <img :src="qrcodeData" alt="支付二维码" class="qrcode-image" />
                  <p class="qrcode-expire">
                    有效期至：{{ qrcodeExpireAt }}
                  </p>
                  <button type="button" @click="refreshQrcode" class="refresh-btn">
                    <RefreshCw :size="16" />
                    刷新
                  </button>
                </div>
              </div>
            </div>
            </div>
          </div>
        </div>

        <!-- 动态表单字段 -->
        <div class="form-section">
          <h2 class="section-title">动态表单字段</h2>
          <p class="section-desc">配置用户参与活动时需要填写的信息字段</p>
          
           <!-- 字段列表 -->
          <div class="field-list" id="field-list-container">
            <div v-for="(field, index) in form.formFields" :key="index" class="field-item">
              <div class="field-info">
                <div class="field-header">
                  <span class="field-name">{{ field.label }}</span>
                  <span class="field-type">{{ getFieldTypeText(field.type) }}</span>
                  <span v-if="field.required" class="field-required">必填</span>
                </div>
                <div class="field-details">
                  <span v-if="field.placeholder" class="field-placeholder">占位符: {{ field.placeholder }}</span>
                  <span v-if="field.options" class="field-options">选项: {{ field.options.join(', ') }}</span>
                  <span v-if="field.validationRules" class="field-validation">
                    <span v-if="field.validationRules.minLength" class="validation-item">最小{{ field.validationRules.minLength }}字符</span>
                    <span v-if="field.validationRules.maxLength" class="validation-item">最大{{ field.validationRules.maxLength }}字符</span>
                    <span v-if="field.validationRules.pattern" class="validation-item">正则验证</span>
                  </span>
                </div>
              </div>
              <div class="field-actions">
                <div class="drag-handle">
                  <GripVertical :size="16" class="drag-icon" />
                </div>
                <button type="button" @click="editField(index)" class="edit-btn">
                  <Edit2 :size="16" />
                </button>
                <button type="button" @click="removeField(index)" class="remove-btn">
                  <Trash2 :size="16" />
                </button>
              </div>
            </div>
            <div v-if="form.formFields.length === 0" class="empty-state">
              暂无表单字段，请添加字段
            </div>
          </div>

          <!-- 添加字段 -->
          <div class="add-field-section">
            <h3 class="subsection-title">添加字段</h3>
            <div class="field-type-buttons">
              <button
                type="button"
                v-for="fieldType in fieldTypes"
                :key="fieldType.type"
                @click="addField(fieldType.type)"
                class="field-type-btn"
              >
                <Plus :size="20" class="field-type-icon" />
                <span class="field-type-name">{{ fieldType.name }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- 表单预览 -->
        <div class="form-section">
          <h2 class="section-title">表单预览</h2>
          <p class="section-desc">用户看到的表单样式预览</p>
          
          <div class="form-preview">
            <div class="preview-header">
              <h3 class="preview-title">{{ form.name || '活动名称' }}</h3>
              <p class="preview-description">{{ form.description || '活动描述' }}</p>
            </div>
            
            <div class="preview-form">
               <div v-for="field in form.formFields" :key="field.name" class="preview-field">
                <label class="preview-label">
                  {{ field.label }}
                  <span v-if="field.required" class="required-mark">*</span>
                </label>

                <input
                  v-if="field.type === 'text'"
                  type="text"
                  :placeholder="field.placeholder"
                  class="preview-input"
                  disabled
                >

                <input
                  v-else-if="field.type === 'phone'"
                  type="tel"
                  :placeholder="field.placeholder || '请输入手机号'"
                  class="preview-input"
                  disabled
                >

                <input
                  v-else-if="field.type === 'email'"
                  type="email"
                  :placeholder="field.placeholder || '请输入邮箱'"
                  class="preview-input"
                  disabled
                >

                <textarea
                  v-else-if="field.type === 'textarea'"
                  :placeholder="field.placeholder || '请输入备注'"
                  class="preview-textarea"
                  rows="3"
                  disabled
                ></textarea>

                <input
                  v-else-if="field.type === 'address'"
                  type="text"
                  :placeholder="field.placeholder || '请输入地址'"
                  class="preview-input"
                  disabled
                >

                <select
                  v-else-if="field.type === 'select'"
                  class="preview-select"
                  disabled
                >
                  <option value="">请选择{{ field.label }}</option>
                  <option v-for="option in field.options" :key="option" :value="option">
                    {{ option }}
                  </option>
                </select>
              </div>
              
              <div v-if="form.formFields.length === 0" class="preview-empty">
                请先添加表单字段
              </div>
              
              <button type="button" class="preview-submit-btn" disabled>
                立即参与（奖励 ¥{{ form.rewardRule || 0 }}）
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>

    <!-- 字段编辑模态框 -->
    <div v-if="showFieldModal" class="modal-overlay" @click="closeFieldModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingFieldIndex >= 0 ? '编辑字段' : '添加字段' }}</h3>
          <button @click="closeFieldModal" class="close-btn">
            <X :size="20" />
          </button>
        </div>
        
          <div class="field-form">
            <label>字段类型</label>
            <select v-model="fieldForm.type" class="form-select" :disabled="editingFieldIndex >= 0">
              <option value="text">文本框</option>
              <option value="phone">手机号</option>
              <option value="email">邮箱</option>
              <option value="textarea">多行文本</option>
              <option value="address">地址</option>
              <option value="select">下拉选择</option>
            </select>
          </div>
          
          <div class="form-group">
            <label>字段标签 *</label>
            <input
              v-model="fieldForm.label"
              type="text"
              class="form-input"
              placeholder="如：姓名、联系方式"
              required
            >
          </div>
          
          <div class="form-group">
            <label>字段名称 *</label>
            <input
              v-model="fieldForm.name"
              type="text"
              class="form-input"
              placeholder="如：name、phone（英文，用于数据存储）"
              required
            >
          </div>
          
          <div v-if="fieldForm.type !== 'phone'" class="form-group">
            <label>占位符</label>
            <input
              v-model="fieldForm.placeholder"
              type="text"
              class="form-input"
              placeholder="输入框的提示文字"
            >
          </div>

          <!-- 验证规则配置 -->
          <div class="validation-rules-section">
            <h4 class="rules-title">验证规则</h4>

            <div class="form-group">
              <label>最小长度</label>
              <input
                v-model.number="fieldForm.validationRules.minLength"
                type="number"
                class="form-input"
                min="0"
                max="1000"
                placeholder="0表示不限制"
              >
            </div>

            <div class="form-group">
              <label>最大长度</label>
              <input
                v-model.number="fieldForm.validationRules.maxLength"
                type="number"
                class="form-input"
                min="1"
                max="1000"
                placeholder="0表示不限制"
              >
            </div>

            <div class="form-group">
              <label>正则表达式</label>
              <input
                v-model="fieldForm.validationRules.pattern"
                type="text"
                class="form-input"
                placeholder="可选，如：^[A-Za-z0-9]+$"
              >
            </div>

            <div class="form-group">
              <label>自定义验证消息</label>
              <input
                v-model="fieldForm.validationRules.customMessage"
                type="text"
                class="form-input"
                placeholder="可选，覆盖默认提示"
              >
            </div>
          </div>
          
          <div v-if="fieldForm.type === 'select'" class="form-group">
            <label>选项列表 *</label>
            <div class="options-list">
              <div v-for="(option, index) in fieldForm.options" :key="index" class="option-item">
                <input
                  v-model="fieldForm.options[index]"
                  type="text"
                  class="form-input"
                  placeholder="选项内容"
                >
                <button type="button" @click="removeOption(index)" class="remove-option-btn">删除</button>
              </div>
            </div>
            <button type="button" @click="addOption" class="add-option-btn">添加选项</button>
          </div>
          
          <div class="form-group">
            <label class="checkbox-label">
              <input
                v-model="fieldForm.required"
                type="checkbox"
                class="checkbox-input"
              >
              <span>必填字段</span>
            </label>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="closeFieldModal" class="cancel-btn">取消</button>
          <button @click="saveField" class="confirm-btn">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, Save, Plus, Edit2, Trash2, X, RefreshCw, QrCode, GripVertical } from 'lucide-vue-next'
import { campaignApi } from '../../services/brandApi.js'
import Sortable from 'sortablejs'

const router = useRouter()
const route = useRoute()

const form = reactive({
  brandId: 1, // TODO: 从用户信息获取
  name: '',
  description: '',
  formFields: [], // 动态表单字段配置
  rewardRule: 0,
  startTime: '',
  endTime: '',
  status: 'active',
  // 支付配置
  requirePayment: 'false',
  paymentType: 'deposit', // deposit（订金）或 full（全款）
  paymentAmount: 0
})

const saving = ref(false)
const loading = ref(false)
const showFieldModal = ref(false)
const editingFieldIndex = ref(-1)

// 二维码预览相关
const showQrcode = ref(false)
const qrcodeData = ref('')
const qrcodeLoading = ref(false)
const qrcodeError = ref('')
const qrcodeExpireAt = ref('')

const fieldForm = reactive({
  type: 'text',
  name: '',
  label: '',
  required: false,
  placeholder: '',
  options: [],
  // 验证规则
  validationRules: {
    minLength: 0,
    maxLength: 100,
    pattern: '',
    customMessage: ''
  }
})

const isEditMode = computed(() => !!route.params.id)

// 支持的字段类型（按照OpenSpec规范）
const fieldTypes = [
  { type: 'text', name: '文本框' },
  { type: 'phone', name: '手机号' },
  { type: 'email', name: '邮箱' },
  { type: 'textarea', name: '多行文本' },
  { type: 'address', name: '地址' },
  { type: 'select', name: '下拉选择' }
]

const statusOptions = [
  { value: 'active', label: '进行中', desc: '活动正常进行，用户可以参与' },
  { value: 'paused', label: '已暂停', desc: '活动暂停，用户无法参与' },
  { value: 'ended', label: '已结束', desc: '活动结束，不再接受新的参与' }
]

const goBack = () => {
  router.back()
}

const getFieldTypeText = (type) => {
  const typeMap = {
    text: '文本框',
    phone: '手机号',
    email: '邮箱',
    textarea: '多行文本',
    address: '地址',
    select: '下拉选择'
  }
  return typeMap[type] || type
}

// 添加字段
const addField = (type) => {
  resetFieldForm()
  fieldForm.type = type

  // 根据字段类型设置默认值
  switch (type) {
    case 'text':
      fieldForm.label = '姓名'
      fieldForm.name = 'name'
      fieldForm.placeholder = '请输入姓名'
      break
    case 'phone':
      fieldForm.label = '手机号'
      fieldForm.name = 'phone'
      fieldForm.placeholder = '请输入手机号'
      fieldForm.required = true
      break
    case 'email':
      fieldForm.label = '邮箱'
      fieldForm.name = 'email'
      fieldForm.placeholder = '请输入邮箱地址'
      fieldForm.required = false
      break
    case 'textarea':
      fieldForm.label = '备注'
      fieldForm.name = 'remarks'
      fieldForm.placeholder = '请输入备注信息'
      fieldForm.required = false
      break
    case 'address':
      fieldForm.label = '收货地址'
      fieldForm.name = 'address'
      fieldForm.placeholder = '请输入详细地址'
      fieldForm.required = true
      break
    case 'select':
      fieldForm.label = '意向课程'
      fieldForm.name = 'course'
      fieldForm.options = ['前端开发', '后端开发']
      break
  }

  showFieldModal.value = true
}

// 编辑字段
const editField = (index) => {
  const field = form.formFields[index]
  Object.assign(fieldForm, {
    type: field.type,
    name: field.name,
    label: field.label,
    required: field.required,
    placeholder: field.placeholder || '',
    options: field.options ? [...field.options] : [],
    validationRules: field.validationRules || {
      minLength: 0,
      maxLength: 100,
      pattern: '',
      customMessage: ''
    }
  })
  editingFieldIndex.value = index
  showFieldModal.value = true
}

// 删除字段
const removeField = (index) => {
  if (confirm('确定要删除这个字段吗？')) {
    form.formFields.splice(index, 1)
  }
}

// 重置字段表单
const resetFieldForm = () => {
  Object.assign(fieldForm, {
    type: 'text',
    name: '',
    label: '',
    required: false,
    placeholder: '',
    options: [],
    validationRules: {
      minLength: 0,
      maxLength: 100,
      pattern: '',
      customMessage: ''
    }
  })
  editingFieldIndex.value = -1
}

// 关闭字段模态框
const closeFieldModal = () => {
  showFieldModal.value = false
  resetFieldForm()
}

// 保存字段
const saveField = () => {
  // 验证字段
  if (!fieldForm.label.trim()) {
    alert('请输入字段标签')
    return
  }
  
  if (!fieldForm.name.trim()) {
    alert('请输入字段名称')
    return
  }
  
  if (fieldForm.type === 'select' && fieldForm.options.length === 0) {
    alert('下拉选择字段至少需要一个选项')
    return
  }
  
  // 检查字段名称是否重复（编辑时排除自己）
  const existingIndex = form.formFields.findIndex((field, index) => 
    field.name === fieldForm.name && index !== editingFieldIndex.value
  )
  if (existingIndex >= 0) {
    alert('字段名称已存在，请使用其他名称')
    return
  }
  
  // 构建字段对象
  const fieldData = {
    type: fieldForm.type,
    name: fieldForm.name,
    label: fieldForm.label,
    required: fieldForm.required,
    validationRules: fieldForm.validationRules
  }

  if (fieldForm.placeholder) {
    fieldData.placeholder = fieldForm.placeholder
  }

  if (fieldForm.type === 'select') {
    fieldData.options = fieldForm.options.filter(option => option.trim())
  }

  // 保存字段
  if (editingFieldIndex.value >= 0) {
    form.formFields[editingFieldIndex.value] = fieldData
  } else {
    form.formFields.push(fieldData)
  }
  
  closeFieldModal()
}

// 添加选项
const addOption = () => {
  fieldForm.options.push('')
}

// 删除选项
const removeOption = (index) => {
  fieldForm.options.splice(index, 1)
}

const loadCampaign = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const campaign = await campaignApi.getCampaign(route.params.id)

    // 填充表单数据
    Object.assign(form, {
      brandId: campaign.brandId,
      name: campaign.name,
      description: campaign.description,
      formFields: [...campaign.formFields],
      rewardRule: campaign.rewardRule,
      startTime: formatDateTimeLocal(campaign.startTime),
      endTime: formatDateTimeLocal(campaign.endTime),
      status: campaign.status,
      // 支付配置
      requirePayment: campaign.paymentConfig?.requirePayment ? 'true' : 'false',
      paymentType: campaign.paymentConfig?.paymentType || 'deposit',
      paymentAmount: campaign.paymentConfig?.paymentAmount || 0
    })
  } catch (error) {
    console.error('加载活动失败:', error)
    alert('加载活动信息失败')
    goBack()
  } finally {
    loading.value = false
    // 初始化sortable
    initSortable()
  }
}

const saveCampaign = async () => {
  // 表单验证
  if (!form.name.trim()) {
    alert('请输入活动名称')
    return
  }

  if (form.name.length < 2 || form.name.length > 100) {
    alert('活动名称长度必须在2-100个字符之间')
    return
  }

  if (form.rewardRule < 0.01 || form.rewardRule > 999.99) {
    alert('奖励金额必须在0.01-999.99元之间')
    return
  }

  if (!form.startTime || !form.endTime) {
    alert('请设置活动时间')
    return
  }

  if (new Date(form.startTime) >= new Date(form.endTime)) {
    alert('结束时间必须晚于开始时间')
    return
  }

  if (form.formFields.length === 0) {
    alert('请至少添加一个表单字段')
    return
  }

  // 支付配置验证
  if (form.requirePayment === 'true' || form.requirePayment === true) {
    if (!form.paymentType || form.paymentType === '') {
      alert('请选择支付类型')
      return
    }
    if (!form.paymentAmount || form.paymentAmount < 0.01) {
      alert('请输入有效的支付金额')
      return
    }
  }

  saving.value = true
  try {
    // 构建paymentConfig对象
    const paymentConfig = (form.requirePayment === 'true' || form.requirePayment === true)
      ? {
          requirePayment: true,
          paymentType: form.paymentType,
          paymentAmount: form.paymentAmount
        }
      : {
          requirePayment: false
        }

    const campaignData = {
      brandId: form.brandId,
      name: form.name,
      description: form.description,
      formFields: form.formFields,
      rewardRule: form.rewardRule,
      startTime: formatDateTime(form.startTime),
      endTime: formatDateTime(form.endTime),
      status: form.status,
      paymentConfig: paymentConfig
    }

    if (isEditMode.value) {
      await campaignApi.updateCampaign(route.params.id, campaignData)
      alert('更新成功')
    } else {
      await campaignApi.createCampaign(campaignData)
      alert('创建成功')
    }

    router.push('/brand/campaigns')
  } catch (error) {
    console.error('保存活动失败:', error)
    const errorMessage = error.message || '保存失败'
    alert(`保存失败: ${errorMessage}`)
  } finally {
    saving.value = false
  }
}

// 格式化日期时间为 datetime-local 输入格式
const formatDateTimeLocal = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

// 格式化日期时间为后端格式
const formatDateTime = (dateTimeLocal) => {
  const date = new Date(dateTimeLocal)
  return date.toISOString().slice(0, 19).replace('T', ' ')
}

// 显示二维码预览
const showQrcodePreview = async () => {
  if (!isEditMode.value) {
    alert('请先保存活动后再生成支付二维码')
    return
  }

  if (!form.requirePayment || form.requirePayment === 'false') {
    alert('免费活动不需要支付二维码')
    return
  }

  showQrcode.value = true
  await generateQrcode()
}

// 生成支付二维码
const generateQrcode = async () => {
  qrcodeLoading.value = true
  qrcodeError.value = ''
  try {
    const resp = await campaignApi.getPaymentQrcode(route.params.id)
    qrcodeData.value = resp.qrcodeURL
    qrcodeExpireAt.value = resp.expireAt
  } catch (error) {
    console.error('获取支付二维码失败:', error)
    qrcodeError.value = error.message || '生成二维码失败'
  } finally {
    qrcodeLoading.value = false
  }
}

// 刷新二维码
const refreshQrcode = async () => {
  await generateQrcode()
}

// 关闭二维码预览
const closeQrcode = () => {
  showQrcode.value = false
}

// 初始化sortable
const initSortable = () => {
  const fieldList = document.getElementById('field-list-container')
  if (!fieldList) return

  new Sortable(fieldList, {
    animation: 150,
    handle: '.drag-handle',
    ghostClass: 'sortable-ghost',
    dragClass: 'sortable-drag',
    onEnd: (evt) => {
      const oldIndex = evt.oldIndex
      const newIndex = evt.newIndex

      if (oldIndex !== newIndex) {
        const movedField = form.formFields.splice(oldIndex, 1)[0]
        form.formFields.splice(newIndex, 0, movedField)
      }
    }
  })
}

onMounted(() => {
  loadCampaign()
})
</script>

<style scoped>
/* 验证规则样式 */
.validation-rules-section {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.rules-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #e5e7eb;
}

/* 新增样式 */
.input-with-unit {
  display: flex;
  align-items: center;
}

.input-unit {
  margin-left: 8px;
  color: #6b7280;
  font-size: 14px;
}

.preview-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: #fafafa;
  color: #374151;
  font-size: 14px;
  resize: none;
}

/* 二维码预览样式 */
.qrcode-actions {
  margin-top: 20px;
  text-align: right;
}

.qrcode-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #667eea;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.qrcode-btn:hover:not(:disabled) {
  background: #5a67d8;
  transform: translateY(-1px);
}

.qrcode-btn:active:not(:disabled) {
  transform: translateY(0);
}

.qrcode-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.qrcode-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.qrcode-modal-content {
  background: white;
  border-radius: 12px;
  padding: 24px;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.qrcode-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.qrcode-modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.qrcode-modal-body {
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.qrcode-loading,
.qrcode-error {
  padding: 60px 20px;
  font-size: 14px;
  color: #6b7280;
  text-align: center;
  background: #f8fafc;
  border-radius: 8px;
  width: 100%;
}

.qrcode-error {
  color: #dc2626;
  background: #fef2f2;
}

.qrcode-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 30px 40px;
  background: #f8fafc;
  border-radius: 8px;
  width: 100%;
}

.qrcode-image {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  border: 2px solid #e2e8f0;
  background: white;
  padding: 8px;
}

.qrcode-expire {
  font-size: 12px;
  color: #64748b;
  margin: 0;
  text-align: center;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: white;
  color: #667eea;
  border: 1px solid #667eea;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover {
  background: #667eea;
  color: white;
  transform: translateY(-1px);
}

.refresh-btn:active {
  transform: translateY(0);
}

<style scoped>
.campaign-editor {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-nav {
  background: white;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 10;
}

.back-btn {
  background: none;
  border: none;
  color: #667eea;
  font-size: 16px;
  cursor: pointer;
 8px;


.nav-title {
  font-size: 18px;
  font-weighd;
  margin: 0;
 3;
}

.save-btn {
  background: #a;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
 nter;
;
}

.save-btn:disabled {
  0.6;

}

.{
16px;
}

.form-section {
  background: white;
  border-radius: 8px;
 
g: 24px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
 ight: 600;

  margin: 0 0 2
  padding-bottom: 12px;
  border-botto
}


  font-size: px;
  color: #6b7280;
 0px 0;
}

.form-group {
  margin-bottom: 2
}

.form-label {
 
14px;
  font-weight;
  color: #37411;
  margin-bottom: 8px;
}

.form-input,
.form-textarea,
.form-select {
 
;
  border: 1px solid
  border-radius: 6px;
  font-size: 14px;
 rit;

  box-sizing: boer-box;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px 0.1);
}

. {
;
  min-height: 80px;
}

.{
;
  grid-template-cofr;
  gap: 16px;
}

.input-with-u
  position: relative;
}

.input-unit {
  position: absolute;
  right: 12px;
 

  color: #6b280;
  font-size: 14px;
}

.subsection-title {
 px;
 600;
  color: #374151;
  margin: 0 0 16px 0;
}

/* 字段列表 */
.field-list {
  background: #f9fa
  border: 1px 
  border-radius: 6px;
 

  min-heigpx;
}

.field-item {
  cursor: grab;
  transition: all 0.2s;
}

.field-item:active {
  cursor: grabbing;
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.drag-handle {
  cursor: grab;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.drag-icon {
  color: #9ca3af;
  transition: color 0.2s;
}

.field-item:hover .drag-icon {
  color: #667eea;
}

.sortable-ghost {
  opacity: 0.5;
  background: #f8fafc;
  border: 2px dashed #667eea;
  border-radius: 8px;
}

.sortable-drag {
  opacity: 0.8;
}

.field-item:last-child {
 
}

.field-info {
  flex: 1;
}

.field-header {
  display: flex;
  center;

  margin-bottom;
}

.field-name {
  font-size: 14px;
  500;
;
}

.

  color: #6b7280
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 1px;
}

uired {
  background: #f
  color: #dc2626;
  padding: 2px
  border-rad4px;
 ;
}

.field-details {
  font-size: 12px;
  color: #6b7280;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.field-placeholder,
.field-options {
  margin-right: 12px;
}

.field-validation {
  display: flex;
  gap: 8px;
  align-items: center;
  background: #fef3c7;
  padding: 4px 8px;
  border-radius: 4px;
}

.validation-item {
  font-size: 11px;
  color: #059669;
  font-weight: 500;
}

older,
.field-options {
  margin-right: 12px;
}

.ctions {
y: flex;
  gap: 8px;
}

.edit-btn,
.remove-btn {
  padding: 4px 8px;
  border: none;
  border-radius: 4;
  font-size: 12px;
 nter;
500;
}

.edit-btn {
 e;

}

.edit-btn:hover
  background: #bfdbfe;
}

{
  background:2;
  color: #dc2626;
}

 {
  background: #ecaca;
}

.empty-state {
 r;
px;
  color: #9ca3
  font-size: 1
}

/* 字段类型按钮 */
.field-type-button
 

  gap: 12px;
}

.field-type-btn {
  display: fex;
  column;
r;
  padding: 16px;
  background: whe;
  border: 2px solid #e5e7
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.field-type-btn:hover {
 3b82f6;
8faff;
}

.field-type-icon {
 x;
 8px;
}

.field-type-name {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

//
{
  background: #;
  border: 1px solid #e5e7eb;
  border-radius
 
}

.preview-header {
  text-align: cter;
 4px;
;
  border-botto
}

.preview-title {
 20px;
 bold;
  color: #1a1aa;
  margin: 0 0 8px 0;
}

.

  color: #6b7280;
  margin: 0;
}

.-form {

  margin: 0o;
}

.preview-field {
 : 16px;
}

.preview-label {
  display: block;
 
00;
  color: #374151;
  margin-bottom: 6px;
}

.required-mark {
  color: #dc2626
  margin-left: 2px;
}

.preview-input,
.preview-select {
  width: 100%;
 x 12px;
d5db;
  border-radi;
  font-size: 14px;
  background: hite;
 der-box;


.preview-empty {
  text-align: c
  padding: 32px;
  color: #9ca3
  font-size: 14;
}

.preview-submit-btn {
  width: 100%;
  padding: 12px;
  background: #3b82f6;
 e;
ne;
  border-radius: 6px;
  font-size: 16px;
 600;
0px;
  opacity: 0.7;
}

/
 {
  position: fixe;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0.5);
  display: flex;
  align-items: cen
  justify-content: center;
 000;
20px;
}

.modal-content {
 

  width: 100%;
  max-width: 400px;
  max-height: 80vh;
 to;
}

.modal-header {
  padding: 20px 0;
  display: flex;
 
ter;
}

.modal-header h3 {
  margin: 0;
 
0;
  color: #1a1a;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
 6b7280;
x;
}

.field-form {
  padding: 20px;
}

.options-list {
  margin-bottom
}

 {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

{
  flex: 1;
}

.remove-opti
  background: #f2e2;
  color: #dc2626;
  border: none;
  padding: 8px 12p
  border-radius: 6px;
 x;
r;
  white-space: nowrap;
}

.
f6;
  color: #374151;
  border: 1px soli1d5db;
 
;
  font-size: 14x;
  cursor: pointer;
}

.
e5e7eb;
}

.checkbox-label {
  flex;
ter;
  gap: 8px;
  cursor: pointer;
}

t {
  width: auto;
  margin: 0;
}

.modal-actions {
  padding: 20px;
  display: flex;
  gap: 12px;
  border-top: 1px solid7eb;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
 
  border>
</stylefr;
  }
}lumns: 1plate-cotem grid-
   uttons {eld-type-b
  
  .fi 1fr;
  }te-columns:rid-templa
    g{.form-row {
  px) h: 768(max-widta 
@medi

}#2563eb;ound: ackgrer {
  bbtn:hovnfirm-

.co
}der: none;
  borr: white;olo;
  cnd: #3b82f6  backgrou{
m-btn confir
}

.b;f9fafd: #kgrouner {
  bacn:hovcel-bt
.can74151;
}
  color: #37eb;
id #e5esolrder: 1px 
  bowhite;ground:  backel-btn {
 
.cancr;
}
ointer: p
  cursoght: 500;nt-wei4px;
  fosize: 1 font-x;
 -radius: 6p