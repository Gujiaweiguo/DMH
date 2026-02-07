// DMH H5 用户端
let authToken = localStorage.getItem('user_token');
let userInfo = JSON.parse(localStorage.getItem('user_info') || 'null');
let campaigns = [];
let myRecords = [];
let currentTab = 'home';

// 初始化
function init() {
    render();
    if (authToken && userInfo) {
        showMainPage();
    }
}

// 渲染应用
function render() {
    document.getElementById('app').innerHTML = `
        <!-- 登录页面 -->
        <div class="login-page" id="loginPage">
            <div class="login-card">
                <div class="logo">
                    <h1>DMH 活动中心</h1>
                    <p>发现精彩活动，赢取丰厚奖励</p>
                </div>
                <form id="loginForm">
                    <div class="form-group">
                        <label>用户名</label>
                        <input type="text" id="username" value="user001" required placeholder="请输入用户名">
                    </div>
                    <div class="form-group">
                        <label>密码</label>
                        <input type="password" id="password" value="123456" required placeholder="请输入密码">
                    </div>
                    <div id="errorMsg"></div>
                    <button type="submit" class="btn" id="loginBtn">登录</button>
                </form>
                <div class="test-info">
                    <p><strong>测试账号</strong></p>
                    <p>用户名: user001 | 密码: 123456</p>
                </div>
            </div>
        </div>

        <!-- 主页面 -->
        <div class="main-page" id="mainPage">
            <div class="header">
                <h1>DMH 活动中心</h1>
                <p>发现精彩活动，赢取丰厚奖励</p>
            </div>
            
            <div class="user-info">
                <div class="avatar">${userInfo?.username?.charAt(0)?.toUpperCase() || 'U'}</div>
                <div class="info">
                    <div class="name">${userInfo?.username || '用户'}</div>
                    <div class="role">普通用户</div>
                </div>
                <button class="logout-btn" onclick="logout()">退出</button>
            </div>

            <div id="pageContent"></div>

            <div class="tab-bar">
                <div class="tab-item ${currentTab === 'home' ? 'active' : ''}" onclick="switchTab('home')">
                    <div class="icon">🏠</div>首页
                </div>
                <div class="tab-item ${currentTab === 'campaigns' ? 'active' : ''}" onclick="switchTab('campaigns')">
                    <div class="icon">🎯</div>活动
                </div>
                <div class="tab-item ${currentTab === 'records' ? 'active' : ''}" onclick="switchTab('records')">
                    <div class="icon">📋</div>我的
                </div>
            </div>
        </div>

        <!-- 报名模态框 -->
        <div class="modal" id="joinModal">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>活动报名</h3>
                    <button class="modal-close" onclick="closeModal('joinModal')">&times;</button>
                </div>
                <div class="modal-body" id="joinFormContent"></div>
            </div>
        </div>
    `;
    bindEvents();
    if (authToken && userInfo) {
        renderPageContent();
    }
}


// 绑定事件
function bindEvents() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}

// 登录处理
async function handleLogin(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const loginBtn = document.getElementById('loginBtn');
    const errorMsg = document.getElementById('errorMsg');
    
    loginBtn.disabled = true;
    loginBtn.textContent = '登录中...';
    errorMsg.innerHTML = '';
    
    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        if (response.ok && data.token) {
            authToken = data.token;
            userInfo = { username, userId: data.userId };
            localStorage.setItem('user_token', authToken);
            localStorage.setItem('user_info', JSON.stringify(userInfo));
            showMainPage();
        } else {
            throw new Error(data.message || '登录失败');
        }
    } catch (error) {
        errorMsg.innerHTML = `<div class="error-msg">登录失败: ${error.message}</div>`;
    } finally {
        loginBtn.disabled = false;
        loginBtn.textContent = '登录';
    }
}

// 退出登录
function logout() {
    localStorage.removeItem('user_token');
    localStorage.removeItem('user_info');
    authToken = null;
    userInfo = null;
    location.reload();
}

// 显示主页面
function showMainPage() {
    document.getElementById('loginPage').classList.add('hidden');
    document.getElementById('mainPage').classList.add('active');
    loadCampaigns();
    loadMyRecords();
    renderPageContent();
}

// 切换标签
function switchTab(tab) {
    currentTab = tab;
    render();
    if (authToken) {
        showMainPage();
    }
}


// 渲染页面内容
function renderPageContent() {
    const content = document.getElementById('pageContent');
    if (!content) return;
    
    switch (currentTab) {
        case 'home':
            content.innerHTML = renderHomePage();
            break;
        case 'campaigns':
            content.innerHTML = renderCampaignsPage();
            break;
        case 'records':
            content.innerHTML = renderRecordsPage();
            break;
    }
}

// 渲染首页
function renderHomePage() {
    const activeCampaigns = campaigns.filter(c => c.status === 'ACTIVE' || c.status === 'active').slice(0, 3);
    return `
        <div class="section">
            <div class="section-title">🔥 热门活动</div>
            ${activeCampaigns.length === 0 ? '<div class="empty-state">暂无活动</div>' : 
              activeCampaigns.map(c => renderCampaignCard(c)).join('')}
        </div>
        <div class="section">
            <div class="section-title">📢 最新公告</div>
            <div style="background: white; padding: 15px; border-radius: 12px;">
                <p style="color: #666; font-size: 14px; line-height: 1.6;">欢迎使用DMH活动中心！参与活动即可获得丰厚奖励，快来看看有哪些精彩活动吧~</p>
            </div>
        </div>
    `;
}

// 渲染活动列表页
function renderCampaignsPage() {
    return `
        <div class="section">
            <div class="section-title">🎯 全部活动</div>
            ${campaigns.length === 0 ? '<div class="empty-state">暂无活动</div>' : 
              campaigns.map(c => renderCampaignCard(c)).join('')}
        </div>
    `;
}

// 渲染我的记录页
function renderRecordsPage() {
    return `
        <div class="section">
            <div class="section-title">📋 我的参与记录</div>
            <div class="my-records">
                ${myRecords.length === 0 ? '<div class="empty-state">暂无参与记录</div>' : 
                  myRecords.map(r => `
                    <div class="record-item">
                        <div class="record-info">
                            <h4>${r.campaignName || '活动'}</h4>
                            <p>报名时间: ${r.createdAt || '-'}</p>
                        </div>
                        <span class="record-status ${r.status === 'approved' ? 'success' : 'pending'}">
                            ${r.status === 'approved' ? '已通过' : '待审核'}
                        </span>
                    </div>
                  `).join('')}
            </div>
        </div>
    `;
}

// 渲染活动卡片
function renderCampaignCard(c) {
    const isActive = c.status === 'ACTIVE' || c.status === 'active';
    return `
        <div class="campaign-card">
            <div class="campaign-img">🎉</div>
            <div class="campaign-content">
                <h3>${c.name}</h3>
                <p>${c.description || '精彩活动等你来参与！'}</p>
                <div class="campaign-meta">
                    <span class="status ${isActive ? 'active' : ''}">${isActive ? '进行中' : '已结束'}</span>
                    <span>👥 ${c.orderCount || 0}人参与</span>
                    <span>📅 ${(c.endTime || '').substring(0, 10)}</span>
                </div>
                <button class="join-btn" ${!isActive ? 'disabled' : ''} onclick="openJoinModal(${c.id})">
                    ${isActive ? '立即参与' : '活动已结束'}
                </button>
            </div>
        </div>
    `;
}


// 加载活动列表
async function loadCampaigns() {
    try {
        const response = await fetch('/api/v1/campaigns?page=1&pageSize=100', {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (response.ok) {
            const data = await response.json();
            campaigns = data.campaigns || data.list || [];
            renderPageContent();
        }
    } catch (error) {
        console.error('加载活动失败:', error);
    }
}

// 加载我的参与记录
async function loadMyRecords() {
    try {
        const response = await fetch('/api/v1/orders/my?page=1&pageSize=100', {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (response.ok) {
            const data = await response.json();
            myRecords = data.orders || data.list || [];
            renderPageContent();
        }
    } catch (error) {
        console.error('加载记录失败:', error);
        myRecords = [];
    }
}

// 打开报名模态框
function openJoinModal(campaignId) {
    const campaign = campaigns.find(c => c.id === campaignId);
    if (!campaign) return;
    
    document.getElementById('joinFormContent').innerHTML = `
        <form id="joinForm" onsubmit="submitJoin(event, ${campaignId})">
            <h4 style="margin-bottom: 15px; color: #333;">${campaign.name}</h4>
            <p style="color: #666; font-size: 14px; margin-bottom: 20px;">${campaign.description || '填写以下信息完成报名'}</p>
            
            <div class="form-group">
                <label>姓名 *</label>
                <input type="text" id="joinName" required placeholder="请输入您的姓名">
            </div>
            <div class="form-group">
                <label>手机号 *</label>
                <input type="tel" id="joinPhone" required placeholder="请输入手机号" pattern="[0-9]{11}">
            </div>
            <div class="form-group">
                <label>备注</label>
                <input type="text" id="joinRemark" placeholder="选填">
            </div>
            
            <button type="submit" class="btn" id="submitBtn">提交报名</button>
        </form>
    `;
    openModal('joinModal');
}

// 提交报名
async function submitJoin(e, campaignId) {
    e.preventDefault();
    const name = document.getElementById('joinName').value;
    const phone = document.getElementById('joinPhone').value;
    const remark = document.getElementById('joinRemark').value;
    const submitBtn = document.getElementById('submitBtn');
    
    submitBtn.disabled = true;
    submitBtn.textContent = '提交中...';
    
    try {
        const response = await fetch('/api/v1/orders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                campaignId,
                formData: { name, phone, remark }
            })
        });
        
        if (response.ok) {
            alert('报名成功！');
            closeModal('joinModal');
            loadCampaigns();
            loadMyRecords();
        } else {
            const data = await response.json();
            throw new Error(data.message || '报名失败');
        }
    } catch (error) {
        alert('报名失败: ' + error.message);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = '提交报名';
    }
}

// 模态框操作
function openModal(id) { document.getElementById(id).classList.add('active'); }
function closeModal(id) { document.getElementById(id).classList.remove('active'); }

// 初始化
init();