export const extractCodeFromText = (text) => {
  if (!text) {
    return ''
  }

  try {
    const url = new URL(text)
    const keys = ['code', 'orderCode', 'verification_code', 'verificationCode']
    for (const key of keys) {
      const value = url.searchParams.get(key)
      if (value) {
        return value
      }
    }
  } catch (_error) {
    return text
  }

  return text
}

export const getPaymentStatusText = (status) => {
  const statusMap = {
    paid: '已支付',
    unpaid: '未支付',
    refunded: '已退款',
  }
  return statusMap[status] || status
}

export const getVerifyStatusText = (status) => {
  const statusMap = {
    unverified: '未核销',
    verified: '已核销',
  }
  return statusMap[status] || status
}

export const getOrderStatusText = (status) => getVerifyStatusText(status)

export const isOperationSuccess = (response) => {
  return Boolean(response && response.code === 0 && response.data)
}

export const buildVerifiedOrderState = (scannedOrder, verifyNotes, nowIso) => {
  return {
    ...scannedOrder,
    verifyStatus: 'verified',
    verifiedBy: '品牌管理员',
    verifiedAt: nowIso,
    notes: verifyNotes,
  }
}

export const buildUnverifiedOrderState = (scannedOrder, cancelReason, nowIso) => {
  return {
    ...scannedOrder,
    verifyStatus: 'unverified',
    cancelledBy: '品牌管理员',
    cancelledAt: nowIso,
    cancelReason,
  }
}
