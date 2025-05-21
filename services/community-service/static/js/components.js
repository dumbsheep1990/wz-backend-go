// 万知相同社区平台 - 组件渲染功能
// 创建社区卡片HTML
function createCommunityCardHtml(community) {
    const defaultImg = '/static/images/community-default.jpg';
    const tags = community.tags && community.tags.length > 0 
        ? community.tags.map(tag => `<span class="community-tag">${tag}</span>`).join('') 
        : '';
        
    return `
        <div class="col">
            <div class="card community-card">
                <img src="${defaultImg}" class="card-img-top" alt="${community.name}">
                <div class="community-body">
                    <h5 class="community-title">
                        <a href="/communities/${community.id}" class="text-decoration-none">${community.name}</a>
                    </h5>
                    <div class="community-meta">
                        <span>${community.member_count || 0}成员</span> • 
                        <span>${community.group_count || 0}群组</span> • 
                        <span>${community.post_count || 0}帖子</span>
                    </div>
                    <div class="community-desc">${community.description || '暂无描述'}</div>
                    <div class="community-tags">
                        ${tags}
                    </div>
                </div>
                <div class="card-footer bg-white">
                    <div class="d-flex justify-content-between align-items-center">
                        <small class="text-muted">创建于 ${formatDate(community.create_time)}</small>
                        <button class="btn btn-sm btn-outline-primary join-community-btn" data-id="${community.id}">加入</button>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// 创建群组项目HTML
function createGroupItemHtml(group) {
    const tags = group.tags && group.tags.length > 0 
        ? group.tags.map(tag => `<span class="community-tag">${tag}</span>`).join('') 
        : '';
        
    return `
        <div class="group-item">
            <div class="group-header">
                <h5 class="group-name">
                    <a href="/groups/${group.id}" class="text-decoration-none">${group.name}</a>
                </h5>
                <div class="group-meta">
                    <span>${group.member_count || 0}成员</span> • 
                    <span>${group.post_count || 0}帖子</span>
                </div>
            </div>
            <div class="group-content">${group.description || '暂无描述'}</div>
            <div class="group-tags">${tags}</div>
            <div class="group-footer">
                <div class="group-owner">创建者: ${group.owner_name}</div>
                <div class="group-actions">
                    <button class="btn btn-sm btn-outline-primary join-group-btn" data-id="${group.id}">加入群组</button>
                </div>
            </div>
        </div>
    `;
}

// 创建帖子项目HTML
function createPostItemHtml(post) {
    return `
        <div class="post-item">
            <div class="post-header">
                <img src="/static/images/avatar-default.png" alt="用户头像" class="post-avatar">
                <span class="post-author">${post.author_name}</span>
                <span class="post-meta">${formatDate(post.create_time)}</span>
            </div>
            <h5 class="post-title">
                <a href="/posts/${post.id}" class="text-decoration-none">${post.title}</a>
            </h5>
            <div class="post-content">${truncateText(post.content, 150)}</div>
            <div class="post-actions">
                <div class="post-action">
                    <i class="bi bi-hand-thumbs-up"></i>
                    <span>${post.like_count || 0} 点赞</span>
                </div>
                <div class="post-action">
                    <i class="bi bi-chat"></i>
                    <span>${post.comment_count || 0} 评论</span>
                </div>
                <div class="post-action">
                    <i class="bi bi-eye"></i>
                    <span>${post.view_count || 0} 查看</span>
                </div>
            </div>
        </div>
    `;
}

// 创建评论项目HTML
function createCommentItemHtml(comment) {
    return `
        <div class="comment-item">
            <div class="comment-header">
                <img src="/static/images/avatar-default.png" alt="用户头像" class="comment-avatar">
                <span class="comment-author">${comment.author_name}</span>
                <span class="comment-meta">${formatDate(comment.create_time)}</span>
            </div>
            <div class="comment-content">${comment.content}</div>
            <div class="comment-actions">
                <button class="btn btn-sm text-muted reply-btn" data-comment-id="${comment.id}">回复</button>
                <button class="btn btn-sm text-muted like-btn">
                    <i class="bi bi-hand-thumbs-up"></i>
                    <span>${comment.like_count || 0}</span>
                </button>
            </div>
        </div>
    `;
}

// 创建社区创建表单HTML
function createCommunityFormHtml() {
    return `
        <div class="modal fade" id="create-community-modal" tabindex="-1" aria-hidden="true">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">创建新社区</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div class="modal-body">
                        <form id="create-community-form">
                            <div class="mb-3">
                                <label for="community-name" class="form-label">社区名称</label>
                                <input type="text" class="form-control" id="community-name" required maxlength="50">
                            </div>
                            <div class="mb-3">
                                <label for="community-desc" class="form-label">社区描述</label>
                                <textarea class="form-control" id="community-desc" rows="3" maxlength="200"></textarea>
                            </div>
                            <div class="mb-3">
                                <label for="community-location" class="form-label">所在地区</label>
                                <input type="text" class="form-control" id="community-location" placeholder="例如: 江苏省-南京市">
                            </div>
                            <div class="mb-3">
                                <label for="community-tags" class="form-label">标签</label>
                                <input type="text" class="form-control" id="community-tags" placeholder="多个标签用逗号分隔">
                                <div class="form-text">最多添加5个标签</div>
                            </div>
                            <div class="text-danger mb-3" id="community-create-error"></div>
                            <button type="submit" class="btn btn-primary w-100">创建社区</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// 创建群组创建表单HTML
function createGroupFormHtml(communities) {
    let communitiesOptions = '';
    if (communities && communities.length > 0) {
        communitiesOptions = communities.map(community => 
            `<option value="${community.id}">${community.name}</option>`
        ).join('');
    }
    
    return `
        <div class="modal fade" id="create-group-modal" tabindex="-1" aria-hidden="true">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">创建新群组</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div class="modal-body">
                        <form id="create-group-form">
                            <div class="mb-3">
                                <label for="group-community" class="form-label">所属社区</label>
                                <select class="form-select" id="group-community" required>
                                    <option value="">请选择社区</option>
                                    ${communitiesOptions}
                                </select>
                            </div>
                            <div class="mb-3">
                                <label for="group-name" class="form-label">群组名称</label>
                                <input type="text" class="form-control" id="group-name" required maxlength="50">
                            </div>
                            <div class="mb-3">
                                <label for="group-desc" class="form-label">群组描述</label>
                                <textarea class="form-control" id="group-desc" rows="3" maxlength="200"></textarea>
                            </div>
                            <div class="mb-3">
                                <label for="group-tags" class="form-label">标签</label>
                                <input type="text" class="form-control" id="group-tags" placeholder="多个标签用逗号分隔">
                                <div class="form-text">最多添加5个标签</div>
                            </div>
                            <div class="text-danger mb-3" id="group-create-error"></div>
                            <button type="submit" class="btn btn-primary w-100">创建群组</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// 创建帖子创建表单HTML
function createPostFormHtml(communities, groups) {
    let communitiesOptions = '';
    if (communities && communities.length > 0) {
        communitiesOptions = communities.map(community => 
            `<option value="${community.id}">${community.name}</option>`
        ).join('');
    }
    
    return `
        <div class="modal fade" id="create-post-modal" tabindex="-1" aria-hidden="true">
            <div class="modal-dialog modal-lg">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">发布新帖子</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div class="modal-body">
                        <form id="create-post-form">
                            <div class="mb-3">
                                <label for="post-community" class="form-label">发布到社区</label>
                                <select class="form-select" id="post-community" required onchange="updateGroupSelect()">
                                    <option value="">请选择社区</option>
                                    ${communitiesOptions}
                                </select>
                            </div>
                            <div class="mb-3">
                                <label for="post-group" class="form-label">发布到群组 (可选)</label>
                                <select class="form-select" id="post-group" disabled>
                                    <option value="">请先选择社区</option>
                                </select>
                            </div>
                            <div class="mb-3">
                                <label for="post-title" class="form-label">标题</label>
                                <input type="text" class="form-control" id="post-title" required maxlength="100">
                            </div>
                            <div class="mb-3">
                                <label for="post-content" class="form-label">内容</label>
                                <textarea class="form-control" id="post-content" rows="6" required></textarea>
                            </div>
                            <div class="mb-3">
                                <label for="post-tags" class="form-label">标签</label>
                                <input type="text" class="form-control" id="post-tags" placeholder="多个标签用逗号分隔">
                            </div>
                            <div class="text-danger mb-3" id="post-create-error"></div>
                            <button type="submit" class="btn btn-primary">发布帖子</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// 格式化日期
function formatDate(dateString) {
    if (!dateString) return '未知时间';
    
    const date = new Date(dateString);
    if (isNaN(date)) return dateString;
    
    const now = new Date();
    const diff = now - date;
    
    // 小于1分钟
    if (diff < 60000) {
        return '刚刚';
    }
    // 小于1小时
    else if (diff < 3600000) {
        return Math.floor(diff / 60000) + '分钟前';
    }
    // 小于24小时
    else if (diff < 86400000) {
        return Math.floor(diff / 3600000) + '小时前';
    }
    // 小于30天
    else if (diff < 2592000000) {
        return Math.floor(diff / 86400000) + '天前';
    }
    // 否则显示完整日期
    else {
        return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`;
    }
}

// 截断文本
function truncateText(text, maxLength) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// 处理表单提交
async function handleCreateCommunity(event) {
    event.preventDefault();
    
    const nameInput = document.getElementById('community-name');
    const descInput = document.getElementById('community-desc');
    const locationInput = document.getElementById('community-location');
    const tagsInput = document.getElementById('community-tags');
    const errorElement = document.getElementById('community-create-error');
    
    errorElement.textContent = '';
    
    const name = nameInput.value.trim();
    const description = descInput.value.trim();
    const location = locationInput.value.trim();
    let tags = tagsInput.value.trim() ? tagsInput.value.split(',').map(tag => tag.trim()) : [];
    
    // 验证
    if (!name) {
        errorElement.textContent = '请输入社区名称';
        return;
    }
    
    if (tags.length > 5) {
        errorElement.textContent = '标签最多只能添加5个';
        return;
    }
    
    // 发送请求创建社区
    const response = await apiRequest('/communities', 'POST', {
        name,
        description,
        owner_id: currentUser.id,
        owner_name: currentUser.name,
        location,
        tags
    });
    
    if (response && response.id) {
        const modal = bootstrap.Modal.getInstance(document.getElementById('create-community-modal'));
        modal.hide();
        
        showNotice('社区创建成功！');
        
        // 重新加载社区列表
        loadCommunitiesPage();
    } else {
        errorElement.textContent = '创建社区失败，请稍后再试';
    }
}
