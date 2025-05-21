// 万知相同社区平台 - 事件处理和初始化

// 初始化函数
async function initApp() {
    // 添加事件监听器
    setupEventListeners();
    
    // 检查用户是否已登录
    currentUser = await getCurrentUser();
    updateAuthUI();
    
    // 处理路由
    handleRoute();
    
    // 添加历史状态变化监听器，实现前端路由
    window.addEventListener('popstate', function() {
        handleRoute();
    });
}

// 设置事件监听器
function setupEventListeners() {
    // 登录表单提交
    document.getElementById('login-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const usernameInput = document.getElementById('login-username');
        const passwordInput = document.getElementById('login-password');
        const errorElement = document.getElementById('login-error');
        
        errorElement.textContent = '';
        
        const username = usernameInput.value.trim();
        const password = passwordInput.value;
        
        if (!username || !password) {
            errorElement.textContent = '请输入用户名和密码';
            return;
        }
        
        const success = await login(username, password);
        if (success) {
            loginModal.hide();
            showNotice('登录成功！');
            updateAuthUI();
            
            // 清空表单
            usernameInput.value = '';
            passwordInput.value = '';
        } else {
            errorElement.textContent = '用户名或密码错误';
        }
    });
    
    // 注册表单提交
    document.getElementById('register-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const usernameInput = document.getElementById('register-username');
        const nameInput = document.getElementById('register-name');
        const passwordInput = document.getElementById('register-password');
        const confirmInput = document.getElementById('register-password-confirm');
        const errorElement = document.getElementById('register-error');
        
        errorElement.textContent = '';
        
        const username = usernameInput.value.trim();
        const name = nameInput.value.trim();
        const password = passwordInput.value;
        const confirm = confirmInput.value;
        
        if (!username || !name || !password) {
            errorElement.textContent = '请填写所有必填字段';
            return;
        }
        
        if (password !== confirm) {
            errorElement.textContent = '两次输入的密码不一致';
            return;
        }
        
        const success = await register(username, name, password);
        if (success) {
            registerModal.hide();
            showNotice('注册成功！');
            updateAuthUI();
            
            // 清空表单
            usernameInput.value = '';
            nameInput.value = '';
            passwordInput.value = '';
            confirmInput.value = '';
        } else {
            errorElement.textContent = '注册失败，用户名可能已被使用';
        }
    });
    
    // 登录按钮点击
    document.getElementById('login-btn').addEventListener('click', function() {
        loginModal.show();
    });
    
    // 注册按钮点击
    document.getElementById('register-btn').addEventListener('click', function() {
        registerModal.show();
    });
    
    // 退出按钮点击
    document.getElementById('logout-btn').addEventListener('click', function() {
        logout();
    });
    
    // 页面导航链接点击
    document.querySelectorAll('a[href^="/"]').forEach(link => {
        link.addEventListener('click', function(e) {
            const href = this.getAttribute('href');
            
            // 跳过外部链接和带有特殊属性的链接
            if (href.startsWith('http') || this.hasAttribute('data-no-route')) {
                return;
            }
            
            e.preventDefault();
            
            // 更新历史状态并处理路由
            history.pushState(null, '', href);
            handleRoute();
        });
    });
}

// 动态添加事件监听器（针对动态创建的元素）
function addDynamicEventListeners() {
    // 社区加入按钮
    document.querySelectorAll('.join-community-btn').forEach(btn => {
        btn.addEventListener('click', async function() {
            if (!currentUser) {
                showNotice('请先登录后再加入社区', 'warning');
                loginModal.show();
                return;
            }
            
            const communityId = this.getAttribute('data-id');
            showNotice('加入社区成功！');
            this.textContent = '已加入';
            this.disabled = true;
        });
    });
    
    // 群组加入按钮
    document.querySelectorAll('.join-group-btn').forEach(btn => {
        btn.addEventListener('click', async function() {
            if (!currentUser) {
                showNotice('请先登录后再加入群组', 'warning');
                loginModal.show();
                return;
            }
            
            const groupId = this.getAttribute('data-id');
            const response = await apiRequest(`/groups/${groupId}/join`, 'POST');
            
            if (response && response.success) {
                showNotice('加入群组成功！');
                this.textContent = '已加入';
                this.disabled = true;
            } else {
                showNotice('加入群组失败，请稍后再试', 'danger');
            }
        });
    });
}

// 显示创建社区表单
async function showCreateCommunityForm() {
    if (!currentUser) {
        showNotice('请先登录后再创建社区', 'warning');
        loginModal.show();
        return;
    }
    
    // 添加模态框到文档
    const modalHTML = createCommunityFormHtml();
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = modalHTML;
    document.body.appendChild(tempDiv.firstElementChild);
    
    // 初始化模态框
    const modal = new bootstrap.Modal(document.getElementById('create-community-modal'));
    modal.show();
    
    // 添加表单提交事件
    document.getElementById('create-community-form').addEventListener('submit', handleCreateCommunity);
    
    // 模态框隐藏后清理
    document.getElementById('create-community-modal').addEventListener('hidden.bs.modal', function() {
        this.remove();
    });
}

// 显示创建群组表单
async function showCreateGroupForm() {
    if (!currentUser) {
        showNotice('请先登录后再创建群组', 'warning');
        loginModal.show();
        return;
    }
    
    // 获取社区列表
    const communitiesData = await apiRequest('/communities');
    const communities = communitiesData && communitiesData.communities ? communitiesData.communities : [];
    
    if (communities.length === 0) {
        showNotice('请先创建或加入一个社区后再创建群组', 'warning');
        return;
    }
    
    // 添加模态框到文档
    const modalHTML = createGroupFormHtml(communities);
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = modalHTML;
    document.body.appendChild(tempDiv.firstElementChild);
    
    // 初始化模态框
    const modal = new bootstrap.Modal(document.getElementById('create-group-modal'));
    modal.show();
    
    // 添加表单提交事件
    document.getElementById('create-group-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const communitySelect = document.getElementById('group-community');
        const nameInput = document.getElementById('group-name');
        const descInput = document.getElementById('group-desc');
        const tagsInput = document.getElementById('group-tags');
        const errorElement = document.getElementById('group-create-error');
        
        errorElement.textContent = '';
        
        const communityId = communitySelect.value;
        const name = nameInput.value.trim();
        const description = descInput.value.trim();
        let tags = tagsInput.value.trim() ? tagsInput.value.split(',').map(tag => tag.trim()) : [];
        
        // 验证
        if (!communityId) {
            errorElement.textContent = '请选择所属社区';
            return;
        }
        
        if (!name) {
            errorElement.textContent = '请输入群组名称';
            return;
        }
        
        if (tags.length > 5) {
            errorElement.textContent = '标签最多只能添加5个';
            return;
        }
        
        // 发送请求创建群组
        const response = await apiRequest('/groups', 'POST', {
            name,
            description,
            community_id: communityId,
            owner_id: currentUser.id,
            owner_name: currentUser.name,
            tags
        });
        
        if (response && response.id) {
            modal.hide();
            showNotice('群组创建成功！');
            
            // 重新加载群组列表
            loadGroupsPage();
        } else {
            errorElement.textContent = '创建群组失败，请稍后再试';
        }
    });
    
    // 模态框隐藏后清理
    document.getElementById('create-group-modal').addEventListener('hidden.bs.modal', function() {
        this.remove();
    });
}

// 创建占位logo图片
function createLogoPlaceholder() {
    // 创建logo占位图片
    const logoPlaceholder = document.createElement('img');
    logoPlaceholder.src = "data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' width='100' height='100' viewBox='0 0 100 100'%3e%3crect width='100' height='100' fill='%23FF6F00'/%3e%3ctext x='50' y='60' font-size='32' text-anchor='middle' fill='white' font-family='Arial, sans-serif'%3e万知%3c/text%3e%3c/svg%3e";
    logoPlaceholder.alt = "万知相同";
    logoPlaceholder.id = "logo-placeholder";
    logoPlaceholder.style.display = "none";
    document.body.appendChild(logoPlaceholder);
}

// 页面加载时初始化应用
document.addEventListener('DOMContentLoaded', function() {
    // 创建logo占位图
    createLogoPlaceholder();
    
    // 为空logo设置占位图
    document.querySelectorAll('.logo-image').forEach(img => {
        img.onerror = function() {
            this.src = document.getElementById('logo-placeholder').src;
        };
    });
    
    // 初始化应用
    initApp();
});
