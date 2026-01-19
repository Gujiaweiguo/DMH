<template>
  <div class="distributor-promotion">
    <van-nav-bar title="推广工具" left-arrow @click-left="$router.back()" />

    <div class="promotion-content">
      <!-- 活动选择 -->
      <van-cell-group inset>
        <van-field
          v-model="selectedCampaignId"
          name="campaign"
          label="选择活动"
          placeholder="请选择要推广的活动"
          is-link
          readonly
          @click="showCampaignPicker = true"
        >
          <template #input>
            {{ selectedCampaignName || '请选择活动' }}
          </template>
        </van-field>
      </van-cell-group>

      <!-- 推广链接信息 -->
      <div class="link-info" v-if="generatedLink">
        <van-cell-group inset>
          <van-cell title="推广链接" :value="generatedLink.link" />
          <van-cell title="推广码" :value="generatedLink.linkCode" />
        </van-cell-group>

        <!-- 二维码 -->
        <div class="qrcode-section">
          <h3>推广二维码</h3>
          <div class="qrcode-container">
            <img :src="qrcodeUrl" alt="推广二维码" v-if="qrcodeUrl" />
            <van-loading v-else-if="loadingQrcode" />
            <div class="qrcode-placeholder" v-else>
              <van-button size="small" type="primary" @click="loadQrcode">
                生成二维码
              </van-button>
            </div>
          </div>
          <p class="qrcode-tip">扫描二维码即可访问活动页面</p>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <van-button
            block
            type="primary"
            icon="link-o"
            @click="copyLink"
          >
            复制推广链接
          </van-button>
          <van-button
            block
            icon="down"
            @click="downloadQrcode"
            v-if="qrcodeUrl"
          >
            下载二维码
          </van-button>
        </div>
      </div>

      <!-- 我的推广链接 -->
      <div class="my-links">
        <h3>我的推广链接</h3>
        <van-empty v-if="links.length === 0" description="暂无推广链接" />
        <van-cell-group inset v-else>
          <van-cell
            v-for="link in links"
            :key="link.linkId"
            :title="link.campaignName || `活动 ${link.campaignId}`"
            :value="`点击 ${link.clickCount} 次`"
            is-link
            @click="viewLink(link)"
          />
        </van-cell-group>
      </div>

      <!-- 推广说明 -->
      <div class="promotion-guide">
        <h3>推广说明</h3>
        <div class="guide-content">
          <p>1. 选择要推广的活动，生成专属推广链接和二维码</p>
          <p>2. 将链接或二维码分享给朋友，引导他们下单</p>
          <p>3. 用户通过您的推广链接下单，您即可获得佣金奖励</p>
          <p>4. 最多支持三级分销，您的下级推广也能为您带来收益</p>
        </div>
      </div>
    </div>

    <!-- 活动选择器 -->
    <van-popup v-model:show="showCampaignPicker" position="bottom">
      <van-picker
        :columns="campaignOptions"
        @confirm="onCampaignConfirm"
        @cancel="showCampaignPicker = false"
      />
    </van-popup>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Toast } from 'vant'
import axios from '@/utils/axios'

export default {
  name: 'DistributorPromotion',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const brandId = ref(parseInt(route.query.brandId) || 0)
    const selectedCampaignId = ref(null)
    const selectedCampaignName = ref('')
    const showCampaignPicker = ref(false)
    const generatedLink = ref(null)
    const qrcodeUrl = ref('')
    const loadingQrcode = ref(false)
    const links = ref([])
    const campaigns = ref([])

    const campaignOptions = computed(() => {
      return campaigns.value.map(c => ({
        text: c.name,
        value: c.id
      }))
    })

    // 加载活动列表
    const loadCampaigns = async () => {
      try {
        const { data } = await axios.get('/api/v1/campaigns', {
          params: { status: 'active', pageSize: 100 }
        })
        if (data.code === 200) {
          campaigns.value = data.data.campaigns || []
        }
      } catch (error) {
        console.error('获取活动列表失败:', error)
      }
    }

    // 加载我的推广链接
    const loadMyLinks = async () => {
      try {
        const { data } = await axios.get('/api/v1/distributor/links')
        if (data.code === 200) {
          links.value = data.data || []
          // 获取活动名称
          for (const link of links.value) {
            const campaign = campaigns.value.find(c => c.id === link.campaignId)
            if (campaign) {
              link.campaignName = campaign.name
            }
          }
        }
      } catch (error) {
        console.error('获取推广链接失败:', error)
      }
    }

    // 选择活动
    const onCampaignConfirm = async ({ selectedOptions }) => {
      const option = selectedOptions[0]
      selectedCampaignId.value = option.value
      selectedCampaignName.value = option.text
      showCampaignPicker.value = false

      // 生成推广链接
      await generateLink(option.value)
    }

    // 生成推广链接
    const generateLink = async (campaignId) => {
      try {
        const { data } = await axios.post('/api/v1/distributor/link/generate', {
          campaignId: campaignId
        })
        if (data.code === 200) {
          generatedLink.value = data.data
          qrcodeUrl.value = ''
          // 重新加载链接列表
          loadMyLinks()
        }
      } catch (error) {
        Toast(error.response?.data?.message || '生成链接失败')
      }
    }

    // 加载二维码
    const loadQrcode = () => {
      loadingQrcode.value = true
      // 使用二维码API
      qrcodeUrl.value = generatedLink.value.qrcodeUrl || `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(generatedLink.value.link)}`
      loadingQrcode.value = false
    }

    // 复制链接
    const copyLink = () => {
      if (generatedLink.value) {
        // 使用 Clipboard API
        if (navigator.clipboard) {
          navigator.clipboard.writeText(generatedLink.value.link)
          Toast('链接已复制')
        } else {
          // 降级方案
          const input = document.createElement('input')
          input.value = generatedLink.value.link
          document.body.appendChild(input)
          input.select()
          document.execCommand('copy')
          document.body.removeChild(input)
          Toast('链接已复制')
        }
      }
    }

    // 下载二维码
    const downloadQrcode = () => {
      if (qrcodeUrl.value) {
        const link = document.createElement('a')
        link.href = qrcodeUrl.value
        link.download = `推广二维码_${generatedLink.value.linkCode}.png`
        link.click()
        Toast('正在下载二维码...')
      }
    }

    // 查看链接详情
    const viewLink = (link) => {
      selectedCampaignId.value = link.campaignId
      const campaign = campaigns.value.find(c => c.id === link.campaignId)
      selectedCampaignName.value = campaign ? campaign.name : `活动 ${link.campaignId}`
      generatedLink.value = link
      qrcodeUrl.value = ''
    }

    onMounted(() => {
      loadCampaigns()
      loadMyLinks()
    })

    return {
      brandId,
      selectedCampaignId,
      selectedCampaignName,
      showCampaignPicker,
      generatedLink,
      qrcodeUrl,
      loadingQrcode,
      links,
      campaignOptions,
      onCampaignConfirm,
      loadQrcode,
      copyLink,
      downloadQrcode,
      viewLink
    }
  }
}
</script>

<style scoped>
.distributor-promotion {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.promotion-content {
  padding: 16px;
}

.van-cell-group {
  margin-bottom: 16px;
}

.link-info {
  margin-top: 16px;
}

.qrcode-section {
  background: white;
  margin: 16px;
  padding: 16px;
  border-radius: 8px;
  text-align: center;
}

.qrcode-section h3 {
  margin: 0 0 16px;
  font-size: 16px;
  color: #333;
}

.qrcode-container {
  width: 200px;
  height: 200px;
  margin: 0 auto 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 8px;
}

.qrcode-container img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.qrcode-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qrcode-tip {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.action-buttons {
  padding: 0 16px;
  margin-bottom: 16px;
}

.action-buttons .van-button {
  margin-bottom: 12px;
}

.my-links {
  margin-top: 24px;
}

.my-links h3 {
  padding: 0 16px;
  font-size: 16px;
  color: #333;
  margin-bottom: 12px;
}

.promotion-guide {
  margin-top: 24px;
  background: white;
  padding: 16px;
  border-radius: 8px;
}

.promotion-guide h3 {
  margin: 0 0 12px;
  font-size: 16px;
  color: #333;
}

.guide-content p {
  margin: 8px 0;
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}
</style>
