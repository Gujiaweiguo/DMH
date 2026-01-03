import React, { useState, useEffect } from 'react';
import {
  Card,
  Tabs,
  Table,
  Button,
  Form,
  Input,
  Switch,
  InputNumber,
  Modal,
  message,
  Tag,
  Space,
  Popconfirm,
  DatePicker,
  Select,
  Descriptions,
  Progress,
  Alert,
} from 'antd';
import {
  SecurityScanOutlined,
  LockOutlined,
  AuditOutlined,
  UserOutlined,
  WarningOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';

const { TabPane } = Tabs;
const { RangePicker } = DatePicker;
const { Option } = Select;

// 类型定义
interface PasswordPolicy {
  id: number;
  minLength: number;
  requireUppercase: boolean;
  requireLowercase: boolean;
  requireNumbers: boolean;
  requireSpecialChars: boolean;
  maxAge: number;
  historyCount: number;
  maxLoginAttempts: number;
  lockoutDuration: number;
  sessionTimeout: number;
  maxConcurrentSessions: number;
  createdAt: string;
  updatedAt: string;
}

interface AuditLog {
  id: number;
  userId?: number;
  username: string;
  action: string;
  resource: string;
  resourceId: string;
  details: string;
  clientIp: string;
  userAgent: string;
  status: string;
  errorMsg: string;
  createdAt: string;
}

interface LoginAttempt {
  id: number;
  userId?: number;
  username: string;
  clientIp: string;
  userAgent: string;
  success: boolean;
  failReason: string;
  createdAt: string;
}

interface UserSession {
  id: string;
  userId: number;
  clientIp: string;
  userAgent: string;
  loginAt: string;
  lastActiveAt: string;
  expiresAt: string;
  status: string;
  createdAt: string;
}

interface SecurityEvent {
  id: number;
  eventType: string;
  severity: string;
  userId?: number;
  username: string;
  clientIp: string;
  userAgent: string;
  description: string;
  details: string;
  handled: boolean;
  handledBy?: number;
  handledAt?: string;
  createdAt: string;
}

interface PasswordStrength {
  score: number;
  level: string;
  message: string;
}

const SecurityManagementView: React.FC = () => {
  const [activeTab, setActiveTab] = useState('policy');
  const [loading, setLoading] = useState(false);
  
  // 密码策略相关状态
  const [passwordPolicy, setPasswordPolicy] = useState<PasswordPolicy | null>(null);
  const [policyForm] = Form.useForm();
  const [policyModalVisible, setPolicyModalVisible] = useState(false);
  
  // 审计日志相关状态
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [auditTotal, setAuditTotal] = useState(0);
  const [auditPage, setAuditPage] = useState(1);
  const [auditFilters, setAuditFilters] = useState<any>({});
  
  // 登录尝试相关状态
  const [loginAttempts, setLoginAttempts] = useState<LoginAttempt[]>([]);
  const [attemptTotal, setAttemptTotal] = useState(0);
  const [attemptPage, setAttemptPage] = useState(1);
  const [attemptFilters, setAttemptFilters] = useState<any>({});
  
  // 用户会话相关状态
  const [userSessions, setUserSessions] = useState<UserSession[]>([]);
  const [sessionTotal, setSessionTotal] = useState(0);
  const [sessionPage, setSessionPage] = useState(1);
  
  // 安全事件相关状态
  const [securityEvents, setSecurityEvents] = useState<SecurityEvent[]>([]);
  const [eventTotal, setEventTotal] = useState(0);
  const [eventPage, setEventPage] = useState(1);
  const [eventFilters, setEventFilters] = useState<any>({});
  
  // 密码强度检查
  const [passwordStrength, setPasswordStrength] = useState<PasswordStrength | null>(null);
  const [strengthModalVisible, setStrengthModalVisible] = useState(false);

  // 获取密码策略
  const fetchPasswordPolicy = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/security/password-policy', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setPasswordPolicy(data);
        policyForm.setFieldsValue(data);
      } else {
        message.error(data.message || '获取密码策略失败');
      }
    } catch (error) {
      message.error('获取密码策略失败');
    } finally {
      setLoading(false);
    }
  };

  // 更新密码策略
  const updatePasswordPolicy = async (values: any) => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/security/password-policy', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(values),
      });
      const data = await response.json();
      if (response.ok) {
        message.success('密码策略更新成功');
        setPasswordPolicy(data);
        setPolicyModalVisible(false);
        fetchPasswordPolicy();
      } else {
        message.error(data.message || '更新密码策略失败');
      }
    } catch (error) {
      message.error('更新密码策略失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取审计日志
  const fetchAuditLogs = async (page = 1, filters = {}) => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        pageSize: '20',
        ...filters,
      });
      const response = await fetch(`/api/v1/security/audit-logs?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setAuditLogs(data.logs || []);
        setAuditTotal(data.total || 0);
        setAuditPage(page);
      } else {
        message.error(data.message || '获取审计日志失败');
      }
    } catch (error) {
      message.error('获取审计日志失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取登录尝试记录
  const fetchLoginAttempts = async (page = 1, filters = {}) => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        pageSize: '20',
        ...filters,
      });
      const response = await fetch(`/api/v1/security/login-attempts?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setLoginAttempts(data.attempts || []);
        setAttemptTotal(data.total || 0);
        setAttemptPage(page);
      } else {
        message.error(data.message || '获取登录尝试记录失败');
      }
    } catch (error) {
      message.error('获取登录尝试记录失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取用户会话
  const fetchUserSessions = async (page = 1) => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        pageSize: '20',
      });
      const response = await fetch(`/api/v1/security/sessions?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setUserSessions(data.sessions || []);
        setSessionTotal(data.total || 0);
        setSessionPage(page);
      } else {
        message.error(data.message || '获取用户会话失败');
      }
    } catch (error) {
      message.error('获取用户会话失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取安全事件
  const fetchSecurityEvents = async (page = 1, filters = {}) => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        pageSize: '20',
        ...filters,
      });
      const response = await fetch(`/api/v1/security/events?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setSecurityEvents(data.events || []);
        setEventTotal(data.total || 0);
        setEventPage(page);
      } else {
        message.error(data.message || '获取安全事件失败');
      }
    } catch (error) {
      message.error('获取安全事件失败');
    } finally {
      setLoading(false);
    }
  };

  // 撤销会话
  const revokeSession = async (sessionId: string) => {
    try {
      const response = await fetch(`/api/v1/security/sessions/${sessionId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        message.success('会话已撤销');
        fetchUserSessions(sessionPage);
      } else {
        message.error(data.message || '撤销会话失败');
      }
    } catch (error) {
      message.error('撤销会话失败');
    }
  };

  // 强制用户下线
  const forceLogoutUser = async (userId: number, reason?: string) => {
    try {
      const response = await fetch(`/api/v1/security/force-logout/${userId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ reason: reason || '管理员操作' }),
      });
      const data = await response.json();
      if (response.ok) {
        message.success('用户已被强制下线');
        fetchUserSessions(sessionPage);
      } else {
        message.error(data.message || '强制下线失败');
      }
    } catch (error) {
      message.error('强制下线失败');
    }
  };

  // 处理安全事件
  const handleSecurityEvent = async (eventId: number, note?: string) => {
    try {
      const response = await fetch(`/api/v1/security/events/${eventId}/handle`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ note: note || '' }),
      });
      const data = await response.json();
      if (response.ok) {
        message.success('安全事件已处理');
        fetchSecurityEvents(eventPage, eventFilters);
      } else {
        message.error(data.message || '处理安全事件失败');
      }
    } catch (error) {
      message.error('处理安全事件失败');
    }
  };

  // 检查密码强度
  const checkPasswordStrength = async (password: string) => {
    try {
      const response = await fetch('/api/v1/security/check-password-strength', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ newPassword: password }),
      });
      const data = await response.json();
      if (response.ok) {
        setPasswordStrength(data);
      } else {
        message.error(data.message || '检查密码强度失败');
      }
    } catch (error) {
      message.error('检查密码强度失败');
    }
  };

  useEffect(() => {
    if (activeTab === 'policy') {
      fetchPasswordPolicy();
    } else if (activeTab === 'audit') {
      fetchAuditLogs(1, auditFilters);
    } else if (activeTab === 'attempts') {
      fetchLoginAttempts(1, attemptFilters);
    } else if (activeTab === 'sessions') {
      fetchUserSessions(1);
    } else if (activeTab === 'events') {
      fetchSecurityEvents(1, eventFilters);
    }
  }, [activeTab]);

  // 审计日志表格列
  const auditColumns: ColumnsType<AuditLog> = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 120,
    },
    {
      title: '资源',
      dataIndex: 'resource',
      key: 'resource',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status) => (
        <Tag color={status === 'success' ? 'green' : 'red'}>
          {status === 'success' ? '成功' : '失败'}
        </Tag>
      ),
    },
    {
      title: 'IP地址',
      dataIndex: 'clientIp',
      key: 'clientIp',
      width: 120,
    },
    {
      title: '详情',
      dataIndex: 'details',
      key: 'details',
      ellipsis: true,
    },
  ];

  // 登录尝试表格列
  const attemptColumns: ColumnsType<LoginAttempt> = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '结果',
      dataIndex: 'success',
      key: 'success',
      width: 80,
      render: (success) => (
        <Tag color={success ? 'green' : 'red'} icon={success ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
          {success ? '成功' : '失败'}
        </Tag>
      ),
    },
    {
      title: 'IP地址',
      dataIndex: 'clientIp',
      key: 'clientIp',
      width: 120,
    },
    {
      title: '失败原因',
      dataIndex: 'failReason',
      key: 'failReason',
      ellipsis: true,
      render: (text) => text || '-',
    },
    {
      title: '用户代理',
      dataIndex: 'userAgent',
      key: 'userAgent',
      ellipsis: true,
    },
  ];

  // 用户会话表格列
  const sessionColumns: ColumnsType<UserSession> = [
    {
      title: '用户ID',
      dataIndex: 'userId',
      key: 'userId',
      width: 80,
    },
    {
      title: '登录时间',
      dataIndex: 'loginAt',
      key: 'loginAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '最后活跃',
      dataIndex: 'lastActiveAt',
      key: 'lastActiveAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '过期时间',
      dataIndex: 'expiresAt',
      key: 'expiresAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '活跃' : '已失效'}
        </Tag>
      ),
    },
    {
      title: 'IP地址',
      dataIndex: 'clientIp',
      key: 'clientIp',
      width: 120,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space>
          <Popconfirm
            title="确定要撤销此会话吗？"
            onConfirm={() => revokeSession(record.id)}
          >
            <Button size="small" danger>撤销</Button>
          </Popconfirm>
          <Popconfirm
            title="确定要强制此用户下线吗？"
            onConfirm={() => forceLogoutUser(record.userId)}
          >
            <Button size="small" danger>强制下线</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 安全事件表格列
  const eventColumns: ColumnsType<SecurityEvent> = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '事件类型',
      dataIndex: 'eventType',
      key: 'eventType',
      width: 120,
    },
    {
      title: '严重程度',
      dataIndex: 'severity',
      key: 'severity',
      width: 100,
      render: (severity) => {
        const colors = {
          low: 'blue',
          medium: 'orange',
          high: 'red',
          critical: 'purple',
        };
        return <Tag color={colors[severity as keyof typeof colors]}>{severity.toUpperCase()}</Tag>;
      },
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: 'IP地址',
      dataIndex: 'clientIp',
      key: 'clientIp',
      width: 120,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'handled',
      key: 'handled',
      width: 80,
      render: (handled) => (
        <Tag color={handled ? 'green' : 'orange'}>
          {handled ? '已处理' : '待处理'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_, record) => (
        !record.handled && (
          <Popconfirm
            title="确定要标记此事件为已处理吗？"
            onConfirm={() => handleSecurityEvent(record.id)}
          >
            <Button size="small" type="primary">处理</Button>
          </Popconfirm>
        )
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title={<><SecurityScanOutlined /> 安全管理</>}>
        <Tabs activeKey={activeTab} onChange={setActiveTab}>
          <TabPane tab={<><LockOutlined /> 密码策略</>} key="policy">
            <Card
              title="密码安全策略"
              extra={
                <Space>
                  <Button
                    type="primary"
                    onClick={() => setPolicyModalVisible(true)}
                    disabled={!passwordPolicy}
                  >
                    修改策略
                  </Button>
                  <Button
                    icon={<ReloadOutlined />}
                    onClick={() => fetchPasswordPolicy()}
                    loading={loading}
                  >
                    刷新
                  </Button>
                </Space>
              }
            >
              {passwordPolicy && (
                <Descriptions column={2} bordered>
                  <Descriptions.Item label="最小长度">{passwordPolicy.minLength} 位</Descriptions.Item>
                  <Descriptions.Item label="需要大写字母">
                    {passwordPolicy.requireUppercase ? '是' : '否'}
                  </Descriptions.Item>
                  <Descriptions.Item label="需要小写字母">
                    {passwordPolicy.requireLowercase ? '是' : '否'}
                  </Descriptions.Item>
                  <Descriptions.Item label="需要数字">
                    {passwordPolicy.requireNumbers ? '是' : '否'}
                  </Descriptions.Item>
                  <Descriptions.Item label="需要特殊字符">
                    {passwordPolicy.requireSpecialChars ? '是' : '否'}
                  </Descriptions.Item>
                  <Descriptions.Item label="密码有效期">{passwordPolicy.maxAge} 天</Descriptions.Item>
                  <Descriptions.Item label="历史密码记录">{passwordPolicy.historyCount} 个</Descriptions.Item>
                  <Descriptions.Item label="最大登录尝试">{passwordPolicy.maxLoginAttempts} 次</Descriptions.Item>
                  <Descriptions.Item label="锁定时长">{passwordPolicy.lockoutDuration} 分钟</Descriptions.Item>
                  <Descriptions.Item label="会话超时">{passwordPolicy.sessionTimeout} 分钟</Descriptions.Item>
                  <Descriptions.Item label="最大并发会话">{passwordPolicy.maxConcurrentSessions} 个</Descriptions.Item>
                  <Descriptions.Item label="最后更新">
                    {dayjs(passwordPolicy.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
                  </Descriptions.Item>
                </Descriptions>
              )}
            </Card>

            <Card title="密码强度检查" style={{ marginTop: 16 }}>
              <Space direction="vertical" style={{ width: '100%' }}>
                <Input.Password
                  placeholder="输入密码以检查强度"
                  onChange={(e) => {
                    if (e.target.value) {
                      checkPasswordStrength(e.target.value);
                    } else {
                      setPasswordStrength(null);
                    }
                  }}
                />
                {passwordStrength && (
                  <div>
                    <Progress
                      percent={passwordStrength.score}
                      status={passwordStrength.score >= 80 ? 'success' : passwordStrength.score >= 60 ? 'normal' : 'exception'}
                      format={() => passwordStrength.level}
                    />
                    <Alert
                      message={passwordStrength.message}
                      type={passwordStrength.score >= 80 ? 'success' : passwordStrength.score >= 60 ? 'info' : 'warning'}
                      style={{ marginTop: 8 }}
                    />
                  </div>
                )}
              </Space>
            </Card>
          </TabPane>

          <TabPane tab={<><AuditOutlined /> 审计日志</>} key="audit">
            <Card
              title="操作审计日志"
              extra={
                <Button
                  icon={<ReloadOutlined />}
                  onClick={() => fetchAuditLogs(auditPage, auditFilters)}
                  loading={loading}
                >
                  刷新
                </Button>
              }
            >
              <Table
                columns={auditColumns}
                dataSource={auditLogs}
                rowKey="id"
                loading={loading}
                pagination={{
                  current: auditPage,
                  total: auditTotal,
                  pageSize: 20,
                  onChange: (page) => fetchAuditLogs(page, auditFilters),
                  showSizeChanger: false,
                  showQuickJumper: true,
                  showTotal: (total) => `共 ${total} 条记录`,
                }}
                scroll={{ x: 1200 }}
              />
            </Card>
          </TabPane>

          <TabPane tab="登录尝试" key="attempts">
            <Card
              title="登录尝试记录"
              extra={
                <Button
                  icon={<ReloadOutlined />}
                  onClick={() => fetchLoginAttempts(attemptPage, attemptFilters)}
                  loading={loading}
                >
                  刷新
                </Button>
              }
            >
              <Table
                columns={attemptColumns}
                dataSource={loginAttempts}
                rowKey="id"
                loading={loading}
                pagination={{
                  current: attemptPage,
                  total: attemptTotal,
                  pageSize: 20,
                  onChange: (page) => fetchLoginAttempts(page, attemptFilters),
                  showSizeChanger: false,
                  showQuickJumper: true,
                  showTotal: (total) => `共 ${total} 条记录`,
                }}
                scroll={{ x: 1200 }}
              />
            </Card>
          </TabPane>

          <TabPane tab={<><UserOutlined /> 用户会话</>} key="sessions">
            <Card
              title="活跃用户会话"
              extra={
                <Button
                  icon={<ReloadOutlined />}
                  onClick={() => fetchUserSessions(sessionPage)}
                  loading={loading}
                >
                  刷新
                </Button>
              }
            >
              <Table
                columns={sessionColumns}
                dataSource={userSessions}
                rowKey="id"
                loading={loading}
                pagination={{
                  current: sessionPage,
                  total: sessionTotal,
                  pageSize: 20,
                  onChange: (page) => fetchUserSessions(page),
                  showSizeChanger: false,
                  showQuickJumper: true,
                  showTotal: (total) => `共 ${total} 条记录`,
                }}
                scroll={{ x: 1200 }}
              />
            </Card>
          </TabPane>

          <TabPane tab={<><WarningOutlined /> 安全事件</>} key="events">
            <Card
              title="安全事件监控"
              extra={
                <Button
                  icon={<ReloadOutlined />}
                  onClick={() => fetchSecurityEvents(eventPage, eventFilters)}
                  loading={loading}
                >
                  刷新
                </Button>
              }
            >
              <Table
                columns={eventColumns}
                dataSource={securityEvents}
                rowKey="id"
                loading={loading}
                pagination={{
                  current: eventPage,
                  total: eventTotal,
                  pageSize: 20,
                  onChange: (page) => fetchSecurityEvents(page, eventFilters),
                  showSizeChanger: false,
                  showQuickJumper: true,
                  showTotal: (total) => `共 ${total} 条记录`,
                }}
                scroll={{ x: 1200 }}
              />
            </Card>
          </TabPane>
        </Tabs>
      </Card>

      {/* 密码策略编辑模态框 */}
      <Modal
        title="修改密码策略"
        open={policyModalVisible}
        onCancel={() => setPolicyModalVisible(false)}
        onOk={() => policyForm.submit()}
        confirmLoading={loading}
        width={600}
      >
        <Form
          form={policyForm}
          layout="vertical"
          onFinish={updatePasswordPolicy}
        >
          <Form.Item
            name="minLength"
            label="最小长度"
            rules={[{ required: true, message: '请输入最小长度' }]}
          >
            <InputNumber min={6} max={50} addonAfter="位" />
          </Form.Item>
          
          <Form.Item name="requireUppercase" label="需要大写字母" valuePropName="checked">
            <Switch />
          </Form.Item>
          
          <Form.Item name="requireLowercase" label="需要小写字母" valuePropName="checked">
            <Switch />
          </Form.Item>
          
          <Form.Item name="requireNumbers" label="需要数字" valuePropName="checked">
            <Switch />
          </Form.Item>
          
          <Form.Item name="requireSpecialChars" label="需要特殊字符" valuePropName="checked">
            <Switch />
          </Form.Item>
          
          <Form.Item
            name="maxAge"
            label="密码有效期"
            rules={[{ required: true, message: '请输入密码有效期' }]}
          >
            <InputNumber min={0} max={365} addonAfter="天" />
          </Form.Item>
          
          <Form.Item
            name="historyCount"
            label="历史密码记录数量"
            rules={[{ required: true, message: '请输入历史密码记录数量' }]}
          >
            <InputNumber min={0} max={20} addonAfter="个" />
          </Form.Item>
          
          <Form.Item
            name="maxLoginAttempts"
            label="最大登录尝试次数"
            rules={[{ required: true, message: '请输入最大登录尝试次数' }]}
          >
            <InputNumber min={3} max={20} addonAfter="次" />
          </Form.Item>
          
          <Form.Item
            name="lockoutDuration"
            label="锁定时长"
            rules={[{ required: true, message: '请输入锁定时长' }]}
          >
            <InputNumber min={5} max={1440} addonAfter="分钟" />
          </Form.Item>
          
          <Form.Item
            name="sessionTimeout"
            label="会话超时时间"
            rules={[{ required: true, message: '请输入会话超时时间' }]}
          >
            <InputNumber min={30} max={1440} addonAfter="分钟" />
          </Form.Item>
          
          <Form.Item
            name="maxConcurrentSessions"
            label="最大并发会话数"
            rules={[{ required: true, message: '请输入最大并发会话数' }]}
          >
            <InputNumber min={1} max={10} addonAfter="个" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SecurityManagementView;