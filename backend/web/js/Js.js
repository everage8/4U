

const subjectSelect = document.getElementById('subjectSelect');
const typeSelect = document.getElementById('typeSelect');
const showButton = document.getElementById('showButton');
const resultsContainer = document.getElementById('resultsContainer');

async function fetchPublicTasks({ subject = '', type = '' } = {}) {
    const params = new URLSearchParams();
    if (subject) params.append('subject', subject);
    if (type) params.append('type', type);

    const response = await fetch(`/api/v1/tasks?${params.toString()}`);
    const json = await response.json();
    if (!response.ok) {
        throw new Error(json.message || 'Ошибка загрузки заданий');
    }
    return json.data || [];

}

function buildSubjectOptions() {
    subjectSelect.innerHTML = '<option value="">Выберите предмет</option>';
    Object.entries(SUBJECTS_DATA).forEach(([key, subject]) => {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = subject.label;
        subjectSelect.appendChild(option);
    });
}

function clearTypeOptions() {
    typeSelect.innerHTML = '<option value="">Сначала выберите предмет</option>';
    typeSelect.disabled = true;
}

function populateTypeOptions(subjectKey) {
    const subject = SUBJECTS_DATA[subjectKey];
    if (!subject) {
        clearTypeOptions();
        return;
    }

    typeSelect.innerHTML = '<option value="">Выберите тип задания</option>';
    subject.types.forEach((type) => {
        const option = document.createElement('option');
        option.value = type.id;
        option.textContent = type.label;
        typeSelect.appendChild(option);
    });
    typeSelect.disabled = false;
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

function createTaskCard(task) {
    const card = document.createElement('article');
    card.className = 'task-card';

    const subjectName = getSubjectLabel(task.subject);
    const typeName = getTypeLabel(task.subject, task.taskType);

    const conditionImageHtml = task.conditionImageUrl ? `
        <img class="task-card-image" src="${escapeHtml(task.conditionImageUrl)}" alt="Иллюстрация к условию задачи №${task.number}">
    ` : '';

    const solutionImageHtml = task.solutionImageUrl ? `
        <img class="task-card-image" src="${escapeHtml(task.solutionImageUrl)}" alt="Иллюстрация к решению задачи №${task.number}">
    ` : '';

    const answerHtml = task.answer ? `
        <div class="answer-highlight">
            <strong>Ответ:</strong> ${escapeHtml(task.answer)}
        </div>
    ` : '';

    const solutionHtml = task.solution ? `
        <div class="solution-text" style="${task.answer ? 'margin-top: 10px;' : ''}">
            ${escapeHtml(task.solution)}
        </div>
    ` : '';

    card.innerHTML = `
        <div class="task-card-header">
            <div>
                <p class="task-meta">${subjectName} · ${typeName}</p>
                <h2 class="task-title">Задача №${task.number}</h2>
            </div>
            <div class="task-controls">
                <button class="toggle-answer" type="button">Показать решение</button>
            </div>
        </div>
        <div style="padding: 0 24px 16px;">
            ${conditionImageHtml}
            <div style="font-size: 15px; line-height: 1.6; color: #374151;">
                ${escapeHtml(task.condition)}
            </div>
        </div>
        <div class="task-answer">
            <div class="answer-box">
                ${answerHtml}
                ${solutionHtml}
                ${solutionImageHtml}
            </div>
        </div>
    `;

    const toggleButton = card.querySelector('.toggle-answer');
    const answerPanel = card.querySelector('.task-answer');

    toggleButton.addEventListener('click', () => {
        const isOpen = card.classList.toggle('open');
        toggleButton.textContent = isOpen ? 'Скрыть решение' : 'Показать решение';

        if (isOpen) {
            answerPanel.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
    });

    return card;
}

function renderResults(tasks) {
    resultsContainer.innerHTML = '';

    if (!tasks || tasks.length === 0) {
        resultsContainer.innerHTML = '<div class="empty-state">По выбранному предмету и типу заданий пока нет.</div>';
        return;
    }

    tasks.forEach((task) => {
        resultsContainer.appendChild(createTaskCard(task));
    });

    renderLatexIn(resultsContainer);
}

subjectSelect.addEventListener('change', () => {
    const selected = subjectSelect.value;
    populateTypeOptions(selected);
});

showButton.addEventListener('click', async () => {
    const subjectKey = subjectSelect.value;
    const typeKey = typeSelect.value;

    if (!subjectKey || !typeKey) {
        resultsContainer.innerHTML = '<div class="empty-state">Выберите предмет и тип, затем нажмите «Показать задания».</div>';
        return;
    }

    resultsContainer.innerHTML = '<div class="empty-state">Загрузка заданий...</div>';

    try {
        const tasks = await fetchPublicTasks({ subject: subjectKey, type: typeKey });
        renderResults(tasks);
    } catch (err) {
        resultsContainer.innerHTML = `<div class="empty-state">Ошибка загрузки: ${escapeHtml(err.message)}</div>`;
    }
});

buildSubjectOptions();
