

function getAuthHeaders() {
    const token = localStorage.getItem('token') || '';
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

async function fetchTasks({ search = '', subject = '', type = '', page = 1, limit = 10 } = {}) {
    const params = new URLSearchParams();
    if (search)  params.append('search', search);
    if (subject) params.append('subject', subject);
    if (type)    params.append('type', type);
    params.append('page', String(page));
    params.append('limit', String(limit));

    const response = await fetch(`/api/v1/admin/tasks?${params.toString()}`, {
        headers: getAuthHeaders()
    });
    const json = await response.json();
    if (!response.ok) {
        throw new Error(json.message || 'Ошибка загрузки заданий');
    }
    return {
        tasks: (json.data && json.data.tasks) || [],
        total: (json.data && json.data.total) || 0
    };

}

async function createTask(taskData) {
    const response = await fetch('/api/v1/admin/tasks', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
            subject: taskData.subject,
            taskType: taskData.taskType,
            condition: taskData.condition,
            conditionImageUrl: taskData.conditionImageUrl || '',
            answer: taskData.answer,
            solution: taskData.solution || '',
            solutionImageUrl: taskData.solutionImageUrl || ''
        })
    });
    const json = await response.json();
    if (!response.ok) {
        throw new Error(json.message || 'Ошибка создания задания');
    }
    return json.data;

}

async function updateTask(id, taskData) {
    const response = await fetch(`/api/v1/admin/tasks/${encodeURIComponent(id)}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({
            subject: taskData.subject,
            taskType: taskData.taskType,
            condition: taskData.condition,
            conditionImageUrl: taskData.conditionImageUrl || '',
            answer: taskData.answer,
            solution: taskData.solution || '',
            solutionImageUrl: taskData.solutionImageUrl || ''
        })
    });
    const json = await response.json();
    if (!response.ok) {
        throw new Error(json.message || 'Ошибка обновления задания');
    }
    return json.data;

}

async function deleteTask(id) {
    const response = await fetch(`/api/v1/admin/tasks/${encodeURIComponent(id)}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
    });
    const json = await response.json();
    if (!response.ok) {
        throw new Error(json.message || 'Ошибка удаления задания');
    }
    return json.data;

}

const state = {
    search: '',
    subject: '',
    type: '',
    page: 1,
    limit: 10,
    total: 0,
    tasks: []
};

let debounceTimer = null;
let currentEditingId = null;
let currentDeletingId = null;

const searchInput = document.getElementById('searchInput');
const subjectFilter = document.getElementById('subjectFilter');
const typeFilter = document.getElementById('typeFilter');
const resetFiltersBtn = document.getElementById('resetFiltersBtn');
const tasksStats = document.getElementById('tasksStats');
const tasksTableBody = document.getElementById('tasksTableBody');
const paginationInfo = document.getElementById('paginationInfo');
const paginationControls = document.getElementById('paginationControls');

const openAddModalBtn = document.getElementById('openAddModalBtn');
const taskModalBackdrop = document.getElementById('taskModalBackdrop');
const closeModalBtn = document.getElementById('closeModalBtn');
const cancelModalBtn = document.getElementById('cancelModalBtn');
const modalTitle = document.getElementById('modalTitle');
const taskForm = document.getElementById('taskForm');
const formErrorBanner = document.getElementById('formErrorBanner');

const modalSubject = document.getElementById('modalSubject');
const modalTaskType = document.getElementById('modalTaskType');
const modalCondition = document.getElementById('modalCondition');
const conditionPreview = document.getElementById('conditionPreview');

const modalConditionImageUrl = document.getElementById('modalConditionImageUrl');
const modalConditionImagePreview = document.getElementById('modalConditionImagePreview');

const modalAnswer = document.getElementById('modalAnswer');
const modalSolution = document.getElementById('modalSolution');
const solutionPreview = document.getElementById('solutionPreview');

const modalSolutionImageUrl = document.getElementById('modalSolutionImageUrl');
const modalSolutionImagePreview = document.getElementById('modalSolutionImagePreview');

const deleteModalBackdrop = document.getElementById('deleteModalBackdrop');
const closeDeleteModalBtn = document.getElementById('closeDeleteModalBtn');
const cancelDeleteBtn = document.getElementById('cancelDeleteBtn');
const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');

const logoutBtn = document.getElementById('logoutBtn');

function populateFilterOptions() {
    subjectFilter.innerHTML = '<option value="">Все предметы</option>';
    Object.entries(SUBJECTS_DATA).forEach(([key, subject]) => {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = subject.label;
        subjectFilter.appendChild(option);
    });

    populateFilterTypeOptions('');
}

function populateFilterTypeOptions(subjectKey) {
    typeFilter.innerHTML = '<option value="">Все типы</option>';
    if (!subjectKey || !SUBJECTS_DATA[subjectKey]) {
        typeFilter.disabled = true;
        return;
    }

    SUBJECTS_DATA[subjectKey].types.forEach(type => {
        const option = document.createElement('option');
        option.value = type.id;
        option.textContent = type.label;
        typeFilter.appendChild(option);
    });
    typeFilter.disabled = false;
}

function populateModalSubjectOptions() {
    modalSubject.innerHTML = '<option value="">Выберите предмет</option>';
    Object.entries(SUBJECTS_DATA).forEach(([key, subject]) => {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = subject.label;
        modalSubject.appendChild(option);
    });

    populateModalTypeOptions('');
}

function populateModalTypeOptions(subjectKey) {
    modalTaskType.innerHTML = subjectKey ? '<option value="">Выберите тип задания</option>' : '<option value="">Сначала выберите предмет</option>';
    if (!subjectKey || !SUBJECTS_DATA[subjectKey]) {
        modalTaskType.disabled = true;
        return;
    }

    SUBJECTS_DATA[subjectKey].types.forEach(type => {
        const option = document.createElement('option');
        option.value = type.id;
        option.textContent = type.label;
        modalTaskType.appendChild(option);
    });
    modalTaskType.disabled = false;
}

async function loadTasks() {
    tasksTableBody.innerHTML = `
        <tr>
            <td colspan="5" style="text-align: center; padding: 40px; color: #6b7280;">
                Загрузка данных...
            </td>
        </tr>
    `;

    try {
        const data = await fetchTasks({
            search: state.search,
            subject: state.subject,
            type: state.type,
            page: state.page,
            limit: state.limit
        });

        state.total = data.total;
        state.tasks = data.tasks;
        renderTasksTable(data.tasks);
        renderPagination();
    } catch (err) {
        tasksTableBody.innerHTML = `
            <tr>
                <td colspan="5" style="text-align: center; padding: 40px; color: #dc2626;">
                    Не удалось загрузить задания: ${escapeHtml(err.message)}
                </td>
            </tr>
        `;
        tasksStats.textContent = 'Ошибка загрузки';
    }
}

function renderTasksTable(tasks) {
    if (tasks.length === 0) {
        tasksTableBody.innerHTML = `
            <tr>
                <td colspan="5" style="text-align: center; padding: 40px; color: #6b7280;">
                    Задания не найдены. Попробуйте изменить параметры поиска или фильтров.
                </td>
            </tr>
        `;
        tasksStats.textContent = 'Найдено заданий: 0';
        return;
    }

    const startItem = (state.page - 1) * state.limit + 1;
    const endItem = Math.min(state.page * state.limit, state.total);
    tasksStats.textContent = `Показано ${startItem}–${endItem} из ${state.total} заданий`;

    tasksTableBody.innerHTML = tasks.map(task => {
        const subjectClass = task.subject === 'mathematics' ? 'admin-badge--math' : 'admin-badge--physics';
        const subjectName = getSubjectLabel(task.subject);
        const typeName = getTypeLabel(task.subject, task.taskType);
        const hasImages = Boolean(task.conditionImageUrl || task.solutionImageUrl);

        return `
            <tr data-id="${task.id}">
                <td style="font-weight: 700; color: #6b7280;">#${task.number}</td>
                <td>
                    <div class="admin-badges-container">
                        <span class="admin-badge ${subjectClass}">${subjectName}</span>
                        <span class="admin-badge admin-badge--type">${typeName}</span>
                    </div>
                </td>
                <td>
                    <div class="admin-task-condition" title="${escapeHtml(task.condition)}">
                        ${escapeHtml(task.condition)}
                    </div>
                    ${hasImages ? `<div style="font-size: 12px; color: #3b82f6; margin-top: 4px;">📷 Прикреплены изображения</div>` : ''}
                </td>
                <td>
                    <span class="admin-task-answer" title="${escapeHtml(task.answer)}">
                        ${escapeHtml(task.answer)}
                    </span>
                </td>
                <td style="text-align: right;">
                    <div class="admin-actions" style="justify-content: flex-end;">
                        <button type="button" class="admin-btn admin-btn--secondary admin-btn--sm edit-btn" data-id="${task.id}">
                            ✏️ Редактировать
                        </button>
                        <button type="button" class="admin-btn admin-btn--danger admin-btn--sm delete-btn" data-id="${task.id}">
                            🗑️
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }).join('');

    renderLatexIn(tasksTableBody);
}

function renderPagination() {
    const totalPages = Math.ceil(state.total / state.limit) || 1;
    paginationInfo.textContent = `Страница ${state.page} из ${totalPages}`;

    let controlsHtml = '';

    controlsHtml += `
        <button type="button" class="admin-page-btn" id="prevPageBtn" ${state.page <= 1 ? 'disabled' : ''}>
            &larr;
        </button>
    `;

    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || (i >= state.page - 1 && i <= state.page + 1)) {
            controlsHtml += `
                <button type="button" class="admin-page-btn ${i === state.page ? 'active' : ''}" data-page="${i}">
                    ${i}
                </button>
            `;
        } else if (i === state.page - 2 || i === state.page + 2) {
            controlsHtml += `<span style="padding: 0 4px; color: #9ca3af;">...</span>`;
        }
    }

    controlsHtml += `
        <button type="button" class="admin-page-btn" id="nextPageBtn" ${state.page >= totalPages ? 'disabled' : ''}>
            &rarr;
        </button>
    `;

    paginationControls.innerHTML = controlsHtml;
}

function escapeHtml(str) {
    if (!str) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

function openTaskModal(task = null) {
    resetFormErrors();

    if (task) {
        currentEditingId = task.id;
        modalTitle.textContent = `Редактировать задание #${task.number}`;
        modalSubject.value = task.subject;
        populateModalTypeOptions(task.subject);
        modalTaskType.value = task.taskType;
        modalCondition.value = task.condition;
        modalConditionImageUrl.value = task.conditionImageUrl || '';
        modalAnswer.value = task.answer;
        modalSolution.value = task.solution || '';
        modalSolutionImageUrl.value = task.solutionImageUrl || '';

        updateImagePreview(task.conditionImageUrl || '', modalConditionImagePreview);
        updateImagePreview(task.solutionImageUrl || '', modalSolutionImagePreview);
        updateLatexPreview(modalCondition, conditionPreview);
        updateLatexPreview(modalSolution, solutionPreview);
    } else {
        currentEditingId = null;
        modalTitle.textContent = 'Добавить новое задание';
        taskForm.reset();
        populateModalTypeOptions('');
        updateImagePreview('', modalConditionImagePreview);
        updateImagePreview('', modalSolutionImagePreview);
        updateLatexPreview(modalCondition, conditionPreview);
        updateLatexPreview(modalSolution, solutionPreview);
    }

    taskModalBackdrop.classList.add('is-open');
}

function closeTaskModal() {
    taskModalBackdrop.classList.remove('is-open');
    currentEditingId = null;
    resetFormErrors();
    updateImagePreview('', modalConditionImagePreview);
    updateImagePreview('', modalSolutionImagePreview);
    updateLatexPreview(modalCondition, conditionPreview);
    updateLatexPreview(modalSolution, solutionPreview);
}

function openDeleteModal(id) {
    currentDeletingId = id;
    deleteModalBackdrop.classList.add('is-open');
}

function closeDeleteModal() {
    deleteModalBackdrop.classList.remove('is-open');
    currentDeletingId = null;
}

function resetFormErrors() {
    formErrorBanner.classList.remove('is-visible');
    [modalSubject, modalTaskType, modalCondition, modalAnswer].forEach(el => {
        el.classList.remove('is-invalid');
    });
}

function validateTaskForm() {
    resetFormErrors();
    let isValid = true;

    if (!modalSubject.value) {
        modalSubject.classList.add('is-invalid');
        isValid = false;
    }
    if (!modalTaskType.value) {
        modalTaskType.classList.add('is-invalid');
        isValid = false;
    }
    if (!modalCondition.value.trim()) {
        modalCondition.classList.add('is-invalid');
        isValid = false;
    }
    if (!modalAnswer.value.trim()) {
        modalAnswer.classList.add('is-invalid');
        isValid = false;
    }

    if (!isValid) {
        formErrorBanner.textContent = 'Заполните все обязательные поля (отмечены звездочкой)';
        formErrorBanner.classList.add('is-visible');
    }

    return isValid;
}

function updateImagePreview(url, imgElement) {
    if (!imgElement) return;
    if (url && (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('data:image/'))) {
        imgElement.src = url;
        imgElement.classList.add('is-visible');
    } else {
        imgElement.src = '';
        imgElement.classList.remove('is-visible');
    }
}

function updateLatexPreview(textareaElement, previewElement) {
    if (!textareaElement || !previewElement) return;
    const text = textareaElement.value.trim();
    if (text) {
        previewElement.textContent = text;
        previewElement.style.display = 'block';
        renderLatexIn(previewElement);
    } else {
        previewElement.textContent = '';
        previewElement.style.display = 'none';
    }
}

function setupEventListeners() {

    searchInput.addEventListener('input', () => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            state.search = searchInput.value;
            state.page = 1;
            loadTasks();
        }, 300);
    });

    subjectFilter.addEventListener('change', () => {
        state.subject = subjectFilter.value;
        state.type = '';
        populateFilterTypeOptions(state.subject);
        state.page = 1;
        loadTasks();
    });

    typeFilter.addEventListener('change', () => {
        state.type = typeFilter.value;
        state.page = 1;
        loadTasks();
    });

    resetFiltersBtn.addEventListener('click', () => {
        searchInput.value = '';
        subjectFilter.value = '';
        typeFilter.value = '';
        populateFilterTypeOptions('');
        state.search = '';
        state.subject = '';
        state.type = '';
        state.page = 1;
        loadTasks();
    });

    modalSubject.addEventListener('change', () => {
        populateModalTypeOptions(modalSubject.value);
    });

    modalConditionImageUrl.addEventListener('input', () => {
        updateImagePreview(modalConditionImageUrl.value.trim(), modalConditionImagePreview);
    });
    modalSolutionImageUrl.addEventListener('input', () => {
        updateImagePreview(modalSolutionImageUrl.value.trim(), modalSolutionImagePreview);
    });

    modalCondition.addEventListener('input', () => {
        updateLatexPreview(modalCondition, conditionPreview);
    });
    modalSolution.addEventListener('input', () => {
        updateLatexPreview(modalSolution, solutionPreview);
    });

    openAddModalBtn.addEventListener('click', () => {
        openTaskModal();
    });

    closeModalBtn.addEventListener('click', closeTaskModal);
    cancelModalBtn.addEventListener('click', closeTaskModal);
    taskModalBackdrop.addEventListener('click', (e) => {
        if (e.target === taskModalBackdrop) closeTaskModal();
    });

    closeDeleteModalBtn.addEventListener('click', closeDeleteModal);
    cancelDeleteBtn.addEventListener('click', closeDeleteModal);
    deleteModalBackdrop.addEventListener('click', (e) => {
        if (e.target === deleteModalBackdrop) closeDeleteModal();
    });

    confirmDeleteBtn.addEventListener('click', async () => {
        if (!currentDeletingId) return;

        try {
            await deleteTask(currentDeletingId);
            closeDeleteModal();

            const maxPage = Math.ceil((state.total - 1) / state.limit) || 1;
            if (state.page > maxPage) {
                state.page = maxPage;
            }
            loadTasks();
        } catch (err) {
            alert('Не удалось удалить задание: ' + err.message);
        }
    });

    taskForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (!validateTaskForm()) {
            return;
        }

        const taskData = {
            subject: modalSubject.value,
            taskType: modalTaskType.value,
            condition: modalCondition.value.trim(),
            conditionImageUrl: modalConditionImageUrl.value.trim(),
            answer: modalAnswer.value.trim(),
            solution: modalSolution.value.trim(),
            solutionImageUrl: modalSolutionImageUrl.value.trim()
        };

        try {
            if (currentEditingId) {
                await updateTask(currentEditingId, taskData);
            } else {
                await createTask(taskData);
                state.page = 1;
            }

            closeTaskModal();
            loadTasks();
        } catch (err) {
            formErrorBanner.textContent = err.message || 'Ошибка сохранения';
            formErrorBanner.classList.add('is-visible');
        }
    });

    tasksTableBody.addEventListener('click', (e) => {
        const editBtn = e.target.closest('.edit-btn');
        const deleteBtn = e.target.closest('.delete-btn');

        if (editBtn) {
            const taskId = editBtn.dataset.id;
            const task = state.tasks.find(t => t.id === taskId);
            if (task) {
                openTaskModal(task);
            }
        } else if (deleteBtn) {
            const taskId = deleteBtn.dataset.id;
            openDeleteModal(taskId);
        }
    });

    paginationControls.addEventListener('click', (e) => {
        const btn = e.target.closest('button');
        if (!btn || btn.disabled) return;

        if (btn.id === 'prevPageBtn') {
            if (state.page > 1) {
                state.page--;
                loadTasks();
            }
        } else if (btn.id === 'nextPageBtn') {
            const totalPages = Math.ceil(state.total / state.limit);
            if (state.page < totalPages) {
                state.page++;
                loadTasks();
            }
        } else if (btn.dataset.page) {
            state.page = Number(btn.dataset.page);
            loadTasks();
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            if (taskModalBackdrop.classList.contains('is-open')) {
                closeTaskModal();
            }
            if (deleteModalBackdrop.classList.contains('is-open')) {
                closeDeleteModal();
            }
        }
    });

    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            localStorage.removeItem('token');
            window.location.href = '/login';
        });
    }
}

document.addEventListener('DOMContentLoaded', () => {

    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login';
        return;
    }

    populateFilterOptions();
    populateModalSubjectOptions();
    setupEventListeners();
    loadTasks();
});
