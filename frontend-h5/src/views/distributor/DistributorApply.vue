<template>
  <div class="distributor-apply">
    <van-nav-bar title="申请成为分销商" left-arrow @click-left="$router.back()" />

    <div class="apply-form">
      <van-cell-group inset>
        <van-cell title="分销商说明" :value="distributorDescription" />
      </van-cell-group>

      <van-form @submit="onSubmit">
        <van-cell-group inset>
          <van-field
            v-model="selectedBrandId"
            name="brandId"
            label="选择品牌"
            placeholder="请选择要推广的品牌"
            is-link
            readonly
            @click="showBrandPicker = true"
            :rules="[{ required: true, message: '请选择品牌' }]"
          >
            <template #input>
              {{ selectedBrandName || '请选择品牌' }}
            </template>
          </van-field>

          <van-field
            v-model="form.reason"
            name="reason"
            label="申请理由"
            type="textarea"
            rows="4"
            placeholder="请简述您申请成为分销商的理由，以及您的推广优势..."
            maxlength="500"
            show-word-limit
            :rules="[{ required: true, message: '请填写申请理由' }]"
          />

          <div class="terms">
            <van-checkbox v-model="agreeTerms">
              我已阅读并同意
              <span class="link" @click.prevent="showTerms">《分销商协议》</span>
            </van-checkbox>
          </div>
        </van-cell-group>

        <div class="submit-section">
          <van-button
            round
            block
            type="primary"
            native-type="submit"
            :loading="submitting"
            :disabled="!agreeTerms"
          >
            提交申请
          </van-button>
        </div>
      </van-form>
    </div>

    <!-- 品牌选择器 -->
    <van-popup v-model:show="showBrandPicker" position="bottom">
      <van-picker
        :columns="brandOptions"
        @confirm="onBrandConfirm"
        @cancel="showBrandPicker = false"
      />
    </van-popup>

    <!-- 协议弹窗 -->
    <van-popup v-model:show="showTermsPopup" round position="bottom" :style="{ height: '80%' }">
      <div class="terms-popup">
        <van-nav-bar
          title="分销商协议"
          left-text="关闭"
          @click-left="showTermsPopup = false"
        />
        <div class="terms-content">
          <h3>一、分销商资格</h3>
          <p>1.1 分销商需经品牌管理员审批后方可获得分销资格。</p>
          <p>1.2 分销商需遵守平台规则和品牌相关规定。</p>

          <h3>二、推广规则</h3>
          <p>2.1 分销商可通过专属推广链接和二维码进行推广。</p>
          <p>2.2 推广过程中不得进行虚假宣传或误导消费者。</p>

          <h3>三、奖励说明</h3>
          <p>3.1 分销商可获得推广订单的佣金奖励。</p>
          <p>3.2 奖励比例根据分销商级别和品牌设置而定。</p>
          <p>3.3 最多支持三级分销，符合相关法规要求。</p>

          <h3>四、违规处理</h3>
          <p>4.1 违反协议规定的分销商将被暂停或取消分销资格。</p>
          <p>4.2 涉嫌违法行为的将移交司法机关处理。</p>
        </div>
        <div class="terms-footer">
          <van-button type="primary" block @click="showTermsPopup = false">
            我已阅读并同意
          </van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast, Dialog } from 'vant'
import axios from '@/utils/axios'
import {
  DISTRIBUTOR_DESCRIPTION,
  getDefaultForm,
  buildBrandOptions,
  canSubmit,
  buildApplyPayload,
  TERMS_CONTENT
} from './distributorApply.logic.js'

export default {
  name: 'DistributorApply',
  setup() {
    const router = useRouter()
    const form = ref(getDefaultForm())
    const selectedBrandId = ref(null)
    const selectedBrandName = ref('')
    const showBrandPicker = ref(false)
    const showTermsPopup = ref(false)
    const agreeTerms = ref(false)
    const submitting = ref(false)
    const brands = ref([])

    const distributorDescription = DISTRIBUTOR_DESCRIPTION

    const brandOptions = computed(() => buildBrandOptions(brands.value))

    // 加载品牌列表
    const loadBrands = async () => {
      try {
		const data = await axios.get('/brands')
        if (data.code === 200) {
          brands.value = data.data.brands || []
        }
      } catch (error) {
        console.error('获取品牌列表失败:', error)
      }
    }

    // 选择品牌
    const onBrandConfirm = ({ selectedOptions }) => {
      const option = selectedOptions[0]
      selectedBrandId.value = option.value
      selectedBrandName.value = option.text
      form.value.brandId = option.value
      showBrandPicker.value = false
    }

    // 显示协议
    const showTerms = () => {
      showTermsPopup.value = true
    }

    // 提交申请
    const onSubmit = async () => {
      if (!canSubmit(form.value, agreeTerms.value)) {
        Toast('请先同意分销商协议')
        return
      }

      submitting.value = true
      try {
		const data = await axios.post('/distributor/apply', buildApplyPayload(form.value))

        if (data.code === 200) {
          Dialog.alert({
            title: '申请提交成功',
            message: '您的分销商申请已提交，请等待品牌管理员审核',
            confirmButtonText: '我知道了'
          }).then(() => {
            router.back()
          })
        } else {
          Toast(data.message || '提交失败')
        }
      } catch (error) {
        const message = error.response?.data?.message || '提交失败，请稍后重试'
        Toast(message)
      } finally {
        submitting.value = false
      }
    }

    onMounted(() => {
      loadBrands()
    })

    return {
      form,
      selectedBrandId,
      selectedBrandName,
      showBrandPicker,
      showTermsPopup,
      agreeTerms,
      submitting,
      brandOptions,
      distributorDescription,
      onBrandConfirm,
      showTerms,
      onSubmit
    }
  }
}
</script>

<style scoped>
.distributor-apply {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.apply-form {
  padding: 16px;
}

.van-cell-group {
  margin-bottom: 16px;
}

.terms {
  padding: 16px;
  font-size: 14px;
  color: #666;
}

.terms .link {
  color: #667eea;
}

.submit-section {
  margin-top: 24px;
  padding: 0 16px;
}

.terms-popup {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.terms-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  font-size: 14px;
  line-height: 1.8;
  color: #333;
}

.terms-content h3 {
  font-size: 16px;
  margin: 16px 0 8px;
  color: #333;
}

.terms-content p {
  margin: 8px 0;
  color: #666;
}

.terms-footer {
  padding: 16px;
  border-top: 1px solid #eee;
}
</style>
