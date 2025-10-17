// Admin Panel JavaScript
// This file handles all admin panel functionality

const API_BASE_URL = window.location.origin + '/api';

// State management
const adminState = {
    apiKey: '',
    currentTab: 'submissions',
    submissions: {
        data: [],
        page: 1,
        limit: 50,
        total: 0,
        sort: 'date'
    },
    users: {
        data: [],
        page: 1,
        limit: 50,
        total: 0,
        roleFilter: '',
        statusFilter: ''
    },
    stats: null,
    selectedSubmission: null,
    selectedUser: null
};

// Initialize admin panel when app.js calls it
function initAdminPanel(apiKey) {
    adminState.apiKey = apiKey;
    setupAdminEventListeners();
    fetchAdminStats();
    loadSubmissions();
}

// Set up all event listeners
function setupAdminEventListeners() {
    // Tab switching
    document.getElementById('tab-submissions').addEventListener('click', () => switchTab('submissions'));
    document.getElementById('tab-users').addEventListener('click', () => switchTab('users'));

    // Submissions controls
    document.getElementById('submissions-sort').addEventListener('change', (e) => {
        adminState.submissions.sort = e.target.value;
        adminState.submissions.page = 1;
        loadSubmissions();
    });
    document.getElementById('refresh-submissions-btn').addEventListener('click', () => {
        adminState.submissions.page = 1;
        loadSubmissions();
    });
    document.getElementById('submissions-prev-btn').addEventListener('click', () => {
        if (adminState.submissions.page > 1) {
            adminState.submissions.page--;
            loadSubmissions();
        }
    });
    document.getElementById('submissions-next-btn').addEventListener('click', () => {
        const maxPage = Math.ceil(adminState.submissions.total / adminState.submissions.limit);
        if (adminState.submissions.page < maxPage) {
            adminState.submissions.page++;
            loadSubmissions();
        }
    });

    // Users controls
    document.getElementById('users-role-filter').addEventListener('change', (e) => {
        adminState.users.roleFilter = e.target.value;
        adminState.users.page = 1;
        loadUsers();
    });
    document.getElementById('users-status-filter').addEventListener('change', (e) => {
        adminState.users.statusFilter = e.target.value;
        adminState.users.page = 1;
        loadUsers();
    });
    document.getElementById('refresh-users-btn').addEventListener('click', () => {
        adminState.users.page = 1;
        loadUsers();
    });
    document.getElementById('users-prev-btn').addEventListener('click', () => {
        if (adminState.users.page > 1) {
            adminState.users.page--;
            loadUsers();
        }
    });
    document.getElementById('users-next-btn').addEventListener('click', () => {
        const maxPage = Math.ceil(adminState.users.total / adminState.users.limit);
        if (adminState.users.page < maxPage) {
            adminState.users.page++;
            loadUsers();
        }
    });

    // Delete submission modal
    document.getElementById('cancel-delete-submission').addEventListener('click', hideDeleteSubmissionModal);
    document.getElementById('confirm-delete-submission').addEventListener('click', confirmDeleteSubmission);

    // Change role modal
    document.getElementById('cancel-change-role').addEventListener('click', hideChangeRoleModal);
    document.getElementById('confirm-change-role').addEventListener('click', confirmChangeRole);

    // Change status modal
    document.getElementById('cancel-change-status').addEventListener('click', hideChangeStatusModal);
    document.getElementById('confirm-change-status').addEventListener('click', confirmChangeStatus);

    // Close modals on overlay click
    document.querySelectorAll('.modal-overlay').forEach(modal => {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.add('hidden');
            }
        });
    });
}

// Tab switching
function switchTab(tab) {
    adminState.currentTab = tab;

    // Update tab buttons
    document.querySelectorAll('.admin-tab').forEach(btn => btn.classList.remove('active'));
    if (tab === 'submissions') {
        document.getElementById('tab-submissions').classList.add('active');
        document.getElementById('admin-submissions-tab').classList.remove('hidden');
        document.getElementById('admin-users-tab').classList.add('hidden');
        if (adminState.submissions.data.length === 0) {
            loadSubmissions();
        }
    } else {
        document.getElementById('tab-users').classList.add('active');
        document.getElementById('admin-submissions-tab').classList.add('hidden');
        document.getElementById('admin-users-tab').classList.remove('hidden');
        if (adminState.users.data.length === 0) {
            loadUsers();
        }
    }
}

// Fetch admin statistics
async function fetchAdminStats() {
    try {
        const response = await fetch(`${API_BASE_URL}/admin/stats`, {
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch stats');
        }

        const stats = await response.json();
        adminState.stats = stats;
        renderStats(stats);
    } catch (error) {
        console.error('Error fetching stats:', error);
        showToast('Failed to load statistics', 'error');
    }
}

// Render statistics
function renderStats(stats) {
    document.getElementById('stat-total-users').textContent = stats.total_users || 0;
    document.getElementById('stat-active-users').textContent = stats.active_users || 0;
    document.getElementById('stat-submissions').textContent = stats.total_submissions || 0;
    document.getElementById('stat-votes').textContent = stats.total_votes || 0;
}

// Load submissions
async function loadSubmissions() {
    const loading = document.getElementById('submissions-loading');
    const list = document.getElementById('submissions-list');
    const pagination = document.getElementById('submissions-pagination');

    loading.classList.remove('hidden');
    list.classList.add('hidden');
    pagination.classList.add('hidden');

    try {
        const params = new URLSearchParams({
            page: adminState.submissions.page,
            limit: adminState.submissions.limit,
            sort: adminState.submissions.sort
        });

        const response = await fetch(`${API_BASE_URL}/admin/submissions?${params}`, {
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load submissions');
        }

        const data = await response.json();
        adminState.submissions.data = data.submissions || [];
        adminState.submissions.total = data.total || 0;

        renderSubmissions();
        updateSubmissionsPagination();
    } catch (error) {
        console.error('Error loading submissions:', error);
        showToast('Failed to load submissions', 'error');
        loading.textContent = 'Failed to load submissions';
    }
}

// Render submissions table
function renderSubmissions() {
    const loading = document.getElementById('submissions-loading');
    const list = document.getElementById('submissions-list');

    if (adminState.submissions.data.length === 0) {
        loading.textContent = 'No submissions found';
        loading.classList.remove('hidden');
        list.classList.add('hidden');
        return;
    }

    loading.classList.add('hidden');
    list.classList.remove('hidden');
    list.innerHTML = '';

    adminState.submissions.data.forEach(submission => {
        const row = document.createElement('div');
        row.className = 'admin-table-row flex items-center gap-4 p-4 bg-bg-secondary hover:bg-bg-hover rounded-lg border border-border transition-colors';

        row.innerHTML = `
            <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium truncate">${escapeHtml(submission.filename)}</span>
                    ${getRoleBadge(submission.user_role)}
                </div>
                <div class="text-sm text-text-secondary">
                    <span>by ${escapeHtml(submission.username)}</span>
                    <span class="mx-2">•</span>
                    <span>${formatDate(submission.created_at)}</span>
                    <span class="mx-2">•</span>
                    <span>👍 ${submission.upvotes} 👎 ${submission.downvotes}</span>
                </div>
            </div>
            <div class="flex gap-2">
                <button onclick="deleteSubmission(${submission.id})" class="btn btn-secondary btn-sm text-error hover:border-error">
                    Delete
                </button>
            </div>
        `;

        list.appendChild(row);
    });
}

// Update submissions pagination
function updateSubmissionsPagination() {
    const pagination = document.getElementById('submissions-pagination');
    const showing = document.getElementById('submissions-showing');
    const prevBtn = document.getElementById('submissions-prev-btn');
    const nextBtn = document.getElementById('submissions-next-btn');

    const start = (adminState.submissions.page - 1) * adminState.submissions.limit + 1;
    const end = Math.min(adminState.submissions.page * adminState.submissions.limit, adminState.submissions.total);

    showing.textContent = `Showing ${start}-${end} of ${adminState.submissions.total}`;
    prevBtn.disabled = adminState.submissions.page === 1;
    nextBtn.disabled = adminState.submissions.page >= Math.ceil(adminState.submissions.total / adminState.submissions.limit);

    pagination.classList.remove('hidden');
}

// Load users
async function loadUsers() {
    const loading = document.getElementById('users-loading');
    const list = document.getElementById('users-list');
    const pagination = document.getElementById('users-pagination');

    loading.classList.remove('hidden');
    list.classList.add('hidden');
    pagination.classList.add('hidden');

    try {
        const params = new URLSearchParams({
            page: adminState.users.page,
            limit: adminState.users.limit
        });

        if (adminState.users.roleFilter) {
            params.append('role', adminState.users.roleFilter);
        }
        if (adminState.users.statusFilter) {
            params.append('status', adminState.users.statusFilter);
        }

        const response = await fetch(`${API_BASE_URL}/admin/users?${params}`, {
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load users');
        }

        const data = await response.json();
        adminState.users.data = data.users || [];
        adminState.users.total = data.total || 0;

        renderUsers();
        updateUsersPagination();
    } catch (error) {
        console.error('Error loading users:', error);
        showToast('Failed to load users', 'error');
        loading.textContent = 'Failed to load users';
    }
}

// Render users table
function renderUsers() {
    const loading = document.getElementById('users-loading');
    const list = document.getElementById('users-list');

    if (adminState.users.data.length === 0) {
        loading.textContent = 'No users found';
        loading.classList.remove('hidden');
        list.classList.add('hidden');
        return;
    }

    loading.classList.add('hidden');
    list.classList.remove('hidden');
    list.innerHTML = '';

    adminState.users.data.forEach(user => {
        const row = document.createElement('div');
        row.className = 'admin-table-row flex items-center gap-4 p-4 bg-bg-secondary hover:bg-bg-hover rounded-lg border border-border transition-colors';

        const statusBadge = user.is_active
            ? '<span class="badge-active text-xs px-2 py-1 rounded-full bg-success/20 text-success border border-success/30">Active</span>'
            : '<span class="badge-suspended text-xs px-2 py-1 rounded-full bg-error/20 text-error border border-error/30">Suspended</span>';

        row.innerHTML = `
            <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium">${escapeHtml(user.username)}</span>
                    ${getRoleBadge(user.role)}
                    ${statusBadge}
                </div>
                <div class="text-sm text-text-secondary">
                    <span>${user.submission_count} submission${user.submission_count !== 1 ? 's' : ''}</span>
                    <span class="mx-2">•</span>
                    <span>Joined ${formatDate(user.created_at)}</span>
                </div>
            </div>
            <div class="flex gap-2">
                <button onclick="showChangeRoleModal(${user.id}, '${escapeHtml(user.username)}', '${user.role}')" class="btn btn-secondary btn-sm">
                    Change Role
                </button>
                <button onclick="showChangeStatusModal(${user.id}, '${escapeHtml(user.username)}', ${user.is_active})" class="btn btn-secondary btn-sm ${user.is_active ? 'text-error hover:border-error' : 'text-success hover:border-success'}">
                    ${user.is_active ? 'Suspend' : 'Activate'}
                </button>
            </div>
        `;

        list.appendChild(row);
    });
}

// Update users pagination
function updateUsersPagination() {
    const pagination = document.getElementById('users-pagination');
    const showing = document.getElementById('users-showing');
    const prevBtn = document.getElementById('users-prev-btn');
    const nextBtn = document.getElementById('users-next-btn');

    const start = (adminState.users.page - 1) * adminState.users.limit + 1;
    const end = Math.min(adminState.users.page * adminState.users.limit, adminState.users.total);

    showing.textContent = `Showing ${start}-${end} of ${adminState.users.total}`;
    prevBtn.disabled = adminState.users.page === 1;
    nextBtn.disabled = adminState.users.page >= Math.ceil(adminState.users.total / adminState.users.limit);

    pagination.classList.remove('hidden');
}

// Delete submission
function deleteSubmission(submissionId) {
    adminState.selectedSubmission = submissionId;
    document.getElementById('delete-submission-modal').classList.remove('hidden');
    document.getElementById('delete-submission-reason').value = '';
}

function hideDeleteSubmissionModal() {
    document.getElementById('delete-submission-modal').classList.add('hidden');
    adminState.selectedSubmission = null;
}

async function confirmDeleteSubmission() {
    if (!adminState.selectedSubmission) return;

    const reason = document.getElementById('delete-submission-reason').value.trim();
    const body = reason ? JSON.stringify({ reason }) : undefined;

    try {
        const response = await fetch(`${API_BASE_URL}/admin/submissions/delete?id=${adminState.selectedSubmission}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`,
                'Content-Type': 'application/json'
            },
            body: body
        });

        if (!response.ok) {
            throw new Error('Failed to delete submission');
        }

        showToast('Submission deleted successfully', 'success');
        hideDeleteSubmissionModal();
        loadSubmissions();
        fetchAdminStats(); // Update stats
    } catch (error) {
        console.error('Error deleting submission:', error);
        showToast('Failed to delete submission', 'error');
    }
}

// Change user role
function showChangeRoleModal(userId, username, currentRole) {
    adminState.selectedUser = userId;
    document.getElementById('change-role-username').textContent = username;
    document.getElementById('change-role-select').value = currentRole;
    document.getElementById('change-role-modal').classList.remove('hidden');
}

function hideChangeRoleModal() {
    document.getElementById('change-role-modal').classList.add('hidden');
    adminState.selectedUser = null;
}

async function confirmChangeRole() {
    if (!adminState.selectedUser) return;

    const newRole = document.getElementById('change-role-select').value;

    try {
        const response = await fetch(`${API_BASE_URL}/admin/users/role?id=${adminState.selectedUser}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ role: newRole })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to change role');
        }

        showToast('User role updated successfully', 'success');
        hideChangeRoleModal();
        loadUsers();
    } catch (error) {
        console.error('Error changing role:', error);
        showToast(error.message || 'Failed to change user role', 'error');
    }
}

// Change user status
function showChangeStatusModal(userId, username, isActive) {
    adminState.selectedUser = { id: userId, isActive };
    document.getElementById('change-status-username').textContent = username;
    document.getElementById('change-status-action').textContent = isActive ? 'suspend' : 'activate';
    document.getElementById('confirm-change-status').textContent = isActive ? 'Suspend' : 'Activate';
    document.getElementById('change-status-reason').value = '';
    document.getElementById('change-status-modal').classList.remove('hidden');
}

function hideChangeStatusModal() {
    document.getElementById('change-status-modal').classList.add('hidden');
    adminState.selectedUser = null;
}

async function confirmChangeStatus() {
    if (!adminState.selectedUser) return;

    const newStatus = !adminState.selectedUser.isActive;
    const reason = document.getElementById('change-status-reason').value.trim();
    const body = {
        is_active: newStatus
    };
    if (reason) {
        body.reason = reason;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/admin/users/status?id=${adminState.selectedUser.id}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${adminState.apiKey}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to change status');
        }

        showToast(`User ${newStatus ? 'activated' : 'suspended'} successfully`, 'success');
        hideChangeStatusModal();
        loadUsers();
        fetchAdminStats(); // Update stats
    } catch (error) {
        console.error('Error changing status:', error);
        showToast(error.message || 'Failed to change user status', 'error');
    }
}

// Toast notification
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    const icon = document.getElementById('toast-icon');
    const messageEl = document.getElementById('toast-message');

    if (type === 'success') {
        icon.textContent = '✓';
        icon.className = 'text-success text-xl';
    } else if (type === 'error') {
        icon.textContent = '✗';
        icon.className = 'text-error text-xl';
    } else {
        icon.textContent = 'ℹ';
        icon.className = 'text-primary text-xl';
    }

    messageEl.textContent = message;
    toast.classList.remove('hidden');

    setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}

// Helper functions
function getRoleBadge(role) {
    const badges = {
        'admin': '<span class="badge-admin text-xs px-2 py-1 rounded-full bg-success/20 text-success border border-success/30 font-medium">Admin</span>',
        'moderator': '<span class="badge-moderator text-xs px-2 py-1 rounded-full bg-primary/20 text-primary border border-primary/30 font-medium">Moderator</span>',
        'user': '<span class="badge-user text-xs px-2 py-1 rounded-full bg-text-muted/20 text-text-muted border border-text-muted/30 font-medium">User</span>'
    };
    return badges[role] || badges['user'];
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
        return 'Today';
    } else if (diffDays === 1) {
        return 'Yesterday';
    } else if (diffDays < 7) {
        return `${diffDays} days ago`;
    } else {
        return date.toLocaleDateString();
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Expose functions to global scope for onclick handlers
window.deleteSubmission = deleteSubmission;
window.showChangeRoleModal = showChangeRoleModal;
window.showChangeStatusModal = showChangeStatusModal;
