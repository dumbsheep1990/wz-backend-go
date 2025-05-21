// 万知相同社区平台前端JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // 全局变量
    const API_BASE_URL = '/api/v1';
    let currentUser = null;
    let currentPage = window.location.pathname;
    const appElement = document.getElementById('app');
    
    // Bootstrap 模态框实例
    const loginModal = new bootstrap.Modal(document.getElementById('login-modal'));
    const registerModal = new bootstrap.Modal(document.getElementById('register-modal'));
    
    // 工具函数
    function showNotice(message, type = 'success') {
        const alertElement = document.createElement('div');
        alertElement.className = `alert alert-${type} alert-notice`;
        alertElement.innerHTML = message;
        document.body.appendChild(alertElement);
        
        setTimeout(() => {
            alertElement.classList.add('alert-show');
        }, 10);
        
        setTimeout(() => {
            alertElement.classList.remove('alert-show');
            setTimeout(() => {
                document.body.removeChild(alertElement);
            }, 300);
        }, 3000);
    }
    
    function getToken() {
        return localStorage.getItem('token');
    }
    
    function setToken(token) {
        localStorage.setItem('token', token);
    }
    
    function clearToken() {
        localStorage.removeItem('token');
    }
    
    function createApiHeaders() {
        const headers = {
            'Content-Type': 'application/json'
        };
        
        const token = getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        return headers;
    }
    
    async function apiRequest(endpoint, method = 'GET', data = null) {
        const url = `${API_BASE_URL}${endpoint}`;
        const options = {
            method,
            headers: createApiHeaders()
        };
        
        if (data && (method === 'POST' || method === 'PUT')) {
            options.body = JSON.stringify(data);
        }
        
        try {
            const response = await fetch(url, options);
            
            // 如果响应是401未授权，清除令牌并刷新页面
            if (response.status === 401) {
                clearToken();
                showNotice('登录已过期，请重新登录', 'warning');
                updateAuthUI();
                return null;
            }
            
            return await response.json();
        } catch (error) {
            console.error('API请求错误:', error);
            return null;
        }
    }
    
    // 认证相关功能
    async function login(username, password) {
        const data = await apiRequest('/auth/login', 'POST', { username, password });
        if (data && data.token) {
            setToken(data.token);
            currentUser = {
                id: data.user_id,
                name: data.user_name
            };
            return true;
        }
        return false;
    }
    
    async function register(username, name, password) {
        const data = await apiRequest('/auth/register', 'POST', { 
            username, 
            name, 
            password 
        });
        
        if (data && data.token) {
            setToken(data.token);
            currentUser = {
                id: data.user_id,
                name: data.user_name
            };
            return true;
        }
        return false;
    }
    
    async function getCurrentUser() {
        if (!getToken()) return null;
        
        const data = await apiRequest('/auth/me');
        if (data && data.user_id) {
            return {
                id: data.user_id,
                name: data.user_name
            };
        }
        
        return null;
    }
    
    function logout() {
        clearToken();
        currentUser = null;
        showNotice('您已成功退出登录');
        updateAuthUI();
    }
    
    // UI更新函数
    function updateAuthUI() {
        const userInfoElement = document.getElementById('user-info');
        const loginBtnsElement = document.getElementById('login-btns');
        
        if (currentUser) {
            userInfoElement.classList.remove('d-none');
            loginBtnsElement.classList.add('d-none');
            document.getElementById('user-name').textContent = currentUser.name;
        } else {
            userInfoElement.classList.add('d-none');
            loginBtnsElement.classList.remove('d-none');
        }
    }
    
    // 路由处理函数
    function handleRoute() {
        currentPage = window.location.pathname;
        
        // 为了简单起见，我们只处理几个主要路由
        if (currentPage === '/' || currentPage === '/index.html') {
            loadHomePage();
        } else if (currentPage.startsWith('/communities')) {
            loadCommunitiesPage();
        } else if (currentPage.startsWith('/groups')) {
            loadGroupsPage();
        } else if (currentPage.startsWith('/posts/')) {
            if (currentPage.includes('/hot')) {
                loadHotPostsPage();
            } else if (currentPage.includes('/new')) {
                loadNewPostsPage();
            } else {
                loadPostDetail(currentPage.split('/').pop());
            }
        } else {
            loadNotFoundPage();
        }
    }
    
    // 页面加载函数
    async function loadHomePage() {
        appElement.innerHTML = '<div class="text-center py-5"><div class="spinner-border" role="status"><span class="visually-hidden">加载中...</span></div></div>';
        
        // 获取推荐社区
        const communitiesData = await apiRequest('/communities?page_size=6');
        // 获取热门群组
        const groupsData = await apiRequest('/groups?page_size=5');
        // 获取最新帖子
        const postsData = await apiRequest('/posts?page_size=10');
        
        let html = `
            <div class="row mb-4">
                <div class="col-md-8">
                    <div class="card border-0 shadow-sm">
                        <div class="card-body">
                            <h4 class="mb-3">欢迎来到万知相同社区</h4>
                            <p>这里是万知相同社区平台，您可以在这里加入感兴趣的社区和群组，与志同道合的朋友交流讨论。</p>
                            <div class="mt-3">
                                <a href="/communities" class="btn btn-primary">浏览社区</a>
                                <a href="/groups" class="btn btn-outline-primary ms-2">查看群组</a>
                            </div>
                        </div>
                    </div>
                    
                    <h5 class="mt-4 mb-3">最新帖子</h5>
                    <div class="posts-list">
        `;
        
        if (postsData && postsData.posts && postsData.posts.length > 0) {
            postsData.posts.forEach(post => {
                html += createPostItemHtml(post);
            });
        } else {
            html += '<div class="text-center py-4 text-muted">暂无帖子</div>';
        }
        
        html += `
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="card border-0 shadow-sm mb-4">
                        <div class="card-header bg-white">
                            <h5 class="mb-0">热门群组</h5>
                        </div>
                        <div class="card-body p-0">
                            <ul class="list-group list-group-flush">
        `;
        
        if (groupsData && groupsData.groups && groupsData.groups.length > 0) {
            groupsData.groups.forEach(group => {
                html += `
                    <li class="list-group-item">
                        <a href="/groups/${group.id}" class="text-decoration-none">${group.name}</a>
                        <span class="badge bg-light text-dark float-end">${group.member_count}人</span>
                    </li>
                `;
            });
        } else {
            html += '<li class="list-group-item text-center text-muted">暂无群组</li>';
        }
        
        html += `
                            </ul>
                        </div>
                        <div class="card-footer bg-white">
                            <a href="/groups" class="text-decoration-none">查看更多 &raquo;</a>
                        </div>
                    </div>
                    
                    <div class="card border-0 shadow-sm">
                        <div class="card-header bg-white">
                            <h5 class="mb-0">社区推荐</h5>
                        </div>
                        <div class="card-body">
                            <div class="row g-2">
        `;
        
        if (communitiesData && communitiesData.communities && communitiesData.communities.length > 0) {
            communitiesData.communities.forEach(community => {
                html += `
                    <div class="col-6">
                        <div class="card community-card">
                            <div class="community-body">
                                <h6 class="community-title">${community.name}</h6>
                                <div class="community-meta">${community.member_count}成员 · ${community.post_count}帖子</div>
                            </div>
                        </div>
                    </div>
                `;
            });
        } else {
            html += '<div class="text-center text-muted">暂无推荐社区</div>';
        }
        
        html += `
                            </div>
                        </div>
                        <div class="card-footer bg-white">
                            <a href="/communities" class="text-decoration-none">查看更多 &raquo;</a>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        appElement.innerHTML = html;
    }
    
    async function loadCommunitiesPage() {
        appElement.innerHTML = '<div class="text-center py-5"><div class="spinner-border" role="status"><span class="visually-hidden">加载中...</span></div></div>';
        
        const communitiesData = await apiRequest('/communities');
        
        let html = `
            <div class="row mb-4">
                <div class="col">
                    <h3>社区列表</h3>
                    <p class="text-muted">探索并加入感兴趣的社区</p>
                </div>
                ${currentUser ? '<div class="col-auto"><button id="create-community-btn" class="btn btn-primary">创建社区</button></div>' : ''}
            </div>
            
            <div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
        `;
        
        if (communitiesData && communitiesData.communities && communitiesData.communities.length > 0) {
            communitiesData.communities.forEach(community => {
                html += createCommunityCardHtml(community);
            });
        } else {
            html += '<div class="col-12 text-center py-5 text-muted">暂无社区，快来创建第一个社区吧！</div>';
        }
        
        html += `
            </div>
        `;
        
        appElement.innerHTML = html;
        
        // 如果用户已登录，添加创建社区按钮事件监听
        if (currentUser) {
            const createBtn = document.getElementById('create-community-btn');
            if (createBtn) {
                createBtn.addEventListener('click', showCreateCommunityForm);
            }
        }
    }
    
    async function loadGroupsPage() {
        appElement.innerHTML = '<div class="text-center py-5"><div class="spinner-border" role="status"><span class="visually-hidden">加载中...</span></div></div>';
        
        const groupsData = await apiRequest('/groups');
        
        let html = `
            <div class="row mb-4">
                <div class="col">
                    <h3>群组列表</h3>
                    <p class="text-muted">发现并加入感兴趣的群组讨论</p>
                </div>
                ${currentUser ? '<div class="col-auto"><button id="create-group-btn" class="btn btn-primary">创建群组</button></div>' : ''}
            </div>
            
            <div class="groups-list">
        `;
        
        if (groupsData && groupsData.groups && groupsData.groups.length > 0) {
            groupsData.groups.forEach(group => {
                html += createGroupItemHtml(group);
            });
        } else {
            html += '<div class="text-center py-5 text-muted">暂无群组，快来创建第一个群组吧！</div>';
        }
        
        html += `
            </div>
        `;
        
        appElement.innerHTML = html;
        
        // 如果用户已登录，添加创建群组按钮事件监听
        if (currentUser) {
            const createBtn = document.getElementById('create-group-btn');
            if (createBtn) {
                createBtn.addEventListener('click', showCreateGroupForm);
            }
        }
    }
