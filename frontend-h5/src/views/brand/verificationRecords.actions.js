import {
  calculateVerificationStats,
  mergeVerificationRecordsWithOrders,
} from './verificationRecords.logic.js'

export const enrichVerificationRecords = async (orderApi, verificationRecords) => {
  if (!verificationRecords || verificationRecords.length === 0) {
    return []
  }

  try {
    const ordersResp = await orderApi.getOrders()
    return mergeVerificationRecordsWithOrders(verificationRecords, ordersResp?.orders)
  } catch (_error) {
    return verificationRecords
  }
}

export const loadVerificationRecordsData = async (orderApi) => {
  const resp = await orderApi.getVerificationRecords()
  const baseRecords = resp?.records || []
  const records = await enrichVerificationRecords(orderApi, baseRecords)
  return {
    records,
    stats: calculateVerificationStats(records),
  }
}

export const unverifyAndReloadVerificationRecords = async (orderApi, orderId) => {
  await orderApi.unverifyOrder('', { orderId })
  return loadVerificationRecordsData(orderApi)
}
