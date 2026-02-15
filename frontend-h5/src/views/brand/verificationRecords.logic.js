export const getVerificationStatusText = (status) => {
  const statusMap = {
    verified: '已核销',
    unverified: '未核销',
    cancelled: '已取消',
  }
  return statusMap[status] || status
}

export const formatDateTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export const filterVerificationRecords = (records, currentStatus, searchKeyword) => {
  let result = records

  if (currentStatus !== 'all') {
    result = result.filter((record) => record.verificationStatus === currentStatus)
  }

  if (searchKeyword) {
    const keyword = searchKeyword.toLowerCase()
    result = result.filter((record) =>
      record.orderId.toString().includes(keyword) ||
      record.userPhone?.includes(keyword),
    )
  }

  return result
}

export const mergeVerificationRecordsWithOrders = (verificationRecords, orders) => {
  const ordersMap = {}
  ;(orders || []).forEach((order) => {
    ordersMap[order.id] = order
  })

  return verificationRecords.map((record) => {
    const order = ordersMap[record.orderId]
    return {
      ...record,
      orderStatus: order?.status || '',
      orderAmount: order?.amount || 0,
      userPhone: order?.phone || '',
      verifiedByName: `用户${record.verifiedBy || '-'}`,
    }
  })
}

export const calculateVerificationStats = (records, todayDateString = new Date().toLocaleDateString('zh-CN')) => {
  const verifiedRecords = records.filter((record) => record.verificationStatus === 'verified')

  return {
    total: verifiedRecords.length,
    today: verifiedRecords.filter((record) => {
      const recordDate = new Date(record.verifiedAt).toLocaleDateString('zh-CN')
      return recordDate === todayDateString
    }).length,
    totalAmount: verifiedRecords.reduce((sum, record) => sum + (record.orderAmount || 0), 0),
  }
}
