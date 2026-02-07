// DMH 分销商端登录页面
let authToken = localStorage.getItem('distributor_token');
let userInfo = JSON.parse(localStorage.getItem('distributor_info') || 'null');

const API_BASE_URL = (() => {
    const configured = window.DMH_API_BASE_URL || localStorage.getItem('dmh_api_base');
    if (configured) return configured.replace(/\/$/, '');
    if (window.location.port === '3101') return 'http://localhost:8889/api/v1';
    return '/api/v1';
})();

const apiFetch = (path, options) => {
    const normalized = path.startsWith('/') ? path : '/' + path;
    return fetch(API_BASE_URL + normalized, options);
};

// 初始化
function init() {
    render();
    if (authToken && userInfo) {
        showMainPage();
        loadDistributorData();
    }
}

// 渲染应用
function render() {
    document.getElementById('app').innerHTML = `
        <div class="login-page" id="loginPage">
            <div class="login-card">
                <div class="logo">
                    <h1>DMH 分销中心</h1>
                    <p>推广活动，轻松赚取佣金</p>
                </div>
                <form id="loginForm">
                    <div class="form-group">
                        <label>用户名</label>
                        <input type="text" id="username" value="distributor001" required placeholder="请输入用户名">
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
                    <p>用户名: distributor001 | 密码: 123456</p>
                </div>
            </div>
        </div>

        <div class="main-page" id="mainPage">
            <div class="header">
                <h1>分销中心</h1>
                <button class="logout-btn" onclick="logout()">退出</button>
            </div>

            <div class="stats">
                <div class="stat-card purple">
                    <div class="number">¥<span id="balance">0.00</span></div>
                    <div class="label">可提现</div>
                </div>
                <div class="stat-card green">
                    <div class="number">¥<span id="totalEarnings">0.00</span></div>
                    <div class="label">累计收益</div>
                </div>
                <div class="stat-card orange">
                    <div class="number" id="totalOrders">0</div>
                    <div class="label">订单数</div>
                </div>
                <div class="stat-card red">
                    <div class="number" id="subordinatesCount">0</div>
                    <div class="label">下级数</div>
                </div>
            </div>

            <div class="section">
                <div class="section-header">
                    <div class="section-title">我的收益</div>
                </div>
                <div id="rewardsList">
                    <div class="empty-state">加载中...</div>
                </div>
            </div>
        </div>
    `;

    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}

// 处理登录
async function handleLogin(e) {
    e.preventDefault();

    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    const loginBtn = document.getElementById('loginBtn');
    const errorMsg = document.getElementById('errorMsg');

    if (!username || !password) {
        showError('请输入用户名和密码');
        return;
    }

    loginBtn.disabled = true;
    loginBtn.textContent = '登录中...';

    try {
        const response = await apiFetch('/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username: username, password: password })
        });

        loginBtn.disabled = false;
        loginBtn.textContent = '登录';

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.message || '登录失败');
        }

        const data = await response.json();
        
        // 检查是否是分销商角色
        if (!data.roles || !data.roles.includes('participant')) {
            showError('该账号不是分销商账号');
            return;
        }

        authToken = data.token;
        userInfo = {
            userId: data.userId,
            username: data.username,
            roles: data.roles
        };

        localStorage.setItem('distributor_token', authToken);
        localStorage.setItem('distributor_info', JSON.stringify(userInfo));

        showMainPage();
        await loadDistributorData();

    } catch (error) {
        showError(error.message);
    }
}

// 加载分销商数据
async function loadDistributorData() {
    try {
        const response = await apiFetch('/user/balance', {
            headers: {
                'Authorization': 'Bearer ' + authToken
            }
        });

        if (response.ok) {
            const data = await response.json();
            if (data.data) {
                document.getElementById('balance').textContent = data.data.balance.toFixed(2);
                document.getElementById('totalEarnings').textContent = data.data.totalReward.toFixed(2);
            }
        }

        // 加载奖励记录
        const rewardsResponse = await apiFetch('/rewards', {
            headers: {
                'Authorization': 'Bearer ' + authToken
            }
        });

        if (rewardsResponse.ok) {
            const rewardsData = await rewardsResponse.json();
            renderRewardsList(rewardsData.rewards || []);
        }

    } catch (error) {
        console.error('Failed to load distributor data:', error);
    }
}

// 渲染奖励列表
function renderRewardsList(rewards) {
    const container = document.getElementById('rewardsList');
    
    if (!rewards || rewards.length === 0) {
        container.innerHTML = '<div class="empty-state">暂无奖励记录</div>';
        return;
    }

    container.innerHTML = rewards.map(reward => 
        '<div class="reward-card">' +
            '<div class="reward-info">' +
                '<div class="reward-amount">+¥' + reward.amount.toFixed(2) + '</div>' +
                '<div class="reward-status ' + reward.status + '">' + getStatusText(reward.status) + '</div>' +
            '</div>' +
            '<div class="reward-date">' + new Date(reward.createdAt).toLocaleString('zh-CN') + '</div>' +
        '</div>'
    ).join('');
}

function getStatusText(status) {
    const statusMap = {
        'pending': '待结算',
        'settled': '已结算',
        'cancelled': '已取消'
    };
    return statusMap[status] || status;
}

// 显示主页面
function showMainPage() {
    document.getElementById('loginPage').classList.add('hidden');
    document.getElementById('mainPage').classList.add('active');
}

// 显示错误信息
function showError(message) {
    const errorMsg = document.getElementById('errorMsg');
    if (errorMsg) {
        errorMsg.textContent = message;
        errorMsg.style.display = 'block';
    }
}

// 退出登录
function logout() {
    localStorage.removeItem('distributor_token');
    localStorage.removeItem('distributor_info');
    authToken = null;
    userInfo = null;
    document.getElementById('loginPage').classList.remove('hidden');
    document.getElementById('mainPage').classList.remove('active');
    if (document.getElementById('errorMsg')) {
        document.getElementById('errorMsg').style.display = 'none';
    }
}

// 启动应用
document.addEventListener('DOMContentLoaded', init);
