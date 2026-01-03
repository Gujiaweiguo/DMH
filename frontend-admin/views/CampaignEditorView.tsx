<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { Save, ArrowLeft, Plus, Trash2 } from 'lucide-vue-next';
import { campaignApi, type CreateCampaignRequest, type UpdateCampaignRequest } from '../services/campaignApi';
import type { Campaign } from '../types';

const props = defineProps<{
  campaignId?: number;
}>();

const emit = defineEmits(['back', 'saved']);

const form = ref<CreateCampaignRequest & { status?: 'active' | 'paused' | 'ended' }>({
  brandId: 1, // 默认品牌ID，TODO: 从品牌列表中选择
  name: '',
  description: '',
  formFields: [],
  rewardRule: 0,
  startTime: '',
  endTime: '',
  status: 'active'
});

const newField = ref('');
const loading = ref(false);
const saving = ref(false);

const isEditMode = computed(() => !!props.campaignId);

const loadCampaign = async () => {
  if (!props.campaignId) return;
  
  loading.value = true;
  try {
    const campaign = await campaignApi.getCampaign(props.campaignId);
    form.value = {
      brandId: campaign.brandId,
      name: campaign.name,
      description: campaign.description,
      formFields: [...campaign.formFields],
      rewardRule: campaign.rewardRule,
      startTime: campaign.startTime,
      endTime: campaign.endTime,
      status: campaign.status
    };
  } catch (error) {
    console.error('Failed to load campaign:', error);
    alert('加载活动失败');
  } finally {
    loading.value = false;
  }
};

const addFormField = () => {
  if (!newField.value.trim()) return;
  form.value.formFields.push(newField.value.trim());
  newField.value = '';
};

const removeFormField = (index: number) => {
  form.value.formFields.splice(index, 1);
};

const validateForm = (): boolean => {
  if (!form.value.name.trim()) {
    alert('请输入活动名称');
    return false;
  }
  if (form.value.name.length < 2 || form.value.name.length > 100) {
    alert('活动名称长度必须在2-100个字符之间');
    return false;
  }
  if (form.value.rewardRule < 0) {
    alert('奖励金额不能为负数');
    return false;
  }
  if (!form.value.startTime || !form.value.endTime) {
    alert('请选择活动时间范围');
    return false;
  }
  if (new Date(form.value.startTime) > new Date(form.value.endTime)) {
    alert('开始时间不能晚于结束时间');
    return false;
  }
  return true;
};

const handleSave = async () => {
  if (!validateForm()) return;

  saving.value = true;
  try {
    if (isEditMode.value && props.campaignId) {
      await campaignApi.updateCampaign(props.campaignId, form.value as UpdateCampaignRequest);
      alert('更新成功');
    } else {
      await campaignApi.createCampaign(form.value);
      alert('创建成功');
    }
    emit('saved');
  } catch (error) {
    console.error('Failed to save campaign:', error);
    const errorMessage = error instanceof Error ? error.message : '保存失败';
    alert(`保存失败: ${errorMessage}`);
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  loadCampaign();
});
</script>

<template>
  <div class="campaign-editor-view">
    <div class="header">
      <button class="btn-back" @click="emit('back')">
        <ArrowLeft :size="20" />
        <span>返回</span>
      </button>
      <h1 class="title">{{ isEditMode ? '编辑活动' : '创建活动' }}</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="form-container">
      <div class="form-section">
        <h2 class="section-title">基本信息</h2>
        
        <div class="form-group">
          <label>活动名称 *</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="请输入活动名称（2-100字符）"
            maxlength="100"
          />
        </div>

        <div class="form-group">
          <label>活动描述</label>
          <textarea
            v-model="form.description"
            placeholder="请输入活动描述"
            rows="4"
          ></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>开始时间 *</label>
            <input v-model="form.startTime" type="datetime-local" />
          </div>
          <div class="form-group">
            <label>结束时间 *</label>
            <input v-model="form.endTime" type="datetime-local" />
          </div>
        </div>

        <div class="form-group">
          <label>奖励金额 *</label>
          <input
            v-model.number="form.rewardRule"
            type="number"
            step="0.01"
            min="0"
            placeholder="请输入奖励金额"
          />
        </div>

        <div v-if="isEditMode" class="form-group">
          <label>活动状态</label>
          <select v-model="form.status">
            <option value="active">进行中</option>
            <option value="paused">已暂停</option>
            <option value="ended">已结束</option>
          </select>
        </div>
      </div>

      <div class="form-section">
        <h2 class="section-title">动态表单字段</h2>
        
        <div class="field-list">
          <div v-for="(field, index) in form.formFields" :key="index" class="field-item">
            <span class="field-name">{{ field }}</span>
            <button class="btn-remove" @click="removeFormField(index)">
              <Trash2 :size="16" />
            </button>
          </div>
          <div v-if="form.formFields.length === 0" class="empty-state">
            暂无表单字段，请添加
          </div>
        </div>

        <div class="add-field">
          <input
            v-model="newField"
            type="text"
            placeholder="输入字段名称（如：姓名、手机号）"
            @keyup.enter="addFormField"
          />
          <button class="btn-add" @click="addFormField">
            <Plus :size="18" />
            <span>添加</span>
          </button>
        </div>
      </div>

      <div class="form-actions">
        <button class="btn-cancel" @click="emit('back')">取消</button>
        <button class="btn-save" @click="handleSave" :disabled="saving">
          <Save :size="20" />
          <span>{{ saving ? '保存中...' : '保存' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.campaign-editor-view {
  padding: 24px;
  max-width: 1000px;
  margin: 0 auto;
}

.header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.btn-back:hover {
  background: #f9fafb;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
}

.loading {
  text-align: center;
  padding: 60px;
  color: #6b7280;
}

.form-container {
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  padding: 24px;
}

.form-section {
  margin-bottom: 32px;
}

.form-section:last-of-type {
  margin-bottom: 0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #f3f4f6;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.field-list {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 16px;
  min-height: 100px;
}

.field-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  margin-bottom: 8px;
}

.field-item:last-child {
  margin-bottom: 0;
}

.field-name {
  font-size: 14px;
  color: #374151;
}

.btn-remove {
  padding: 6px;
  border: none;
  background: transparent;
  color: #ef4444;
  cursor: pointer;
  border-radius: 4px;
}

.btn-remove:hover {
  background: #fee2e2;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: #9ca3af;
  font-size: 14px;
}

.add-field {
  display: flex;
  gap: 12px;
}

.add-field input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-add:hover {
  background: #2563eb;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel {
  padding: 10px 24px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-cancel:hover {
  background: #f9fafb;
}

.btn-save {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-save:hover:not(:disabled) {
  background: #2563eb;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
