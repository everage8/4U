

const SUBJECTS_DATA = {
    mathematics: {
        label: 'Математика (профиль)',
        types: [
            { id: 'probability_numbers', label: 'Теория вер./чисел' },
            { id: 'transformations',     label: 'Преобразования' },
            { id: 'planimetry',          label: 'Планиметрия' },
            { id: 'stereometry',         label: 'Стереометрия' },
            { id: 'equations',           label: 'Уравнения' },
            { id: 'word_problems',       label: 'Текстовое задание' },
            { id: 'function_graphs',     label: 'Графики функций' },
            { id: 'inequalities',        label: 'Неравенства' },
            { id: 'derivatives',         label: 'Производные/первообразные' }
        ]
    },
    physics: {
        label: 'Физика',
        types: [
            { id: 'magnetism',           label: 'Магнетизм' },
            { id: 'molecular_physics',   label: 'Молекулярная физика и термодинамика' },
            { id: 'si_units',            label: 'Перевод в СИ' },
            { id: 'quantum_nuclear',     label: 'Квантовая/Ядерная физика' },
            { id: 'mechanics',           label: 'Механика' },
            { id: 'oscillations_waves',  label: 'Механические колебания и волны' },
            { id: 'electrodynamics',     label: 'Электродинамика' },
            { id: 'optics',              label: 'Оптика' }
        ]
    }
};

const MOCK_TASKS_STORAGE_KEY = 'mock_tasks_db';

function getMockTasksDB() {
    try {
        const stored = localStorage.getItem(MOCK_TASKS_STORAGE_KEY);
        if (!stored) {
            saveMockTasksDB(DEFAULT_MOCK_TASKS);
            return [...DEFAULT_MOCK_TASKS];
        }
        return JSON.parse(stored);
    } catch (e) {
        console.error('Ошибка чтения mock_tasks_db из localStorage:', e);
        return [...DEFAULT_MOCK_TASKS];
    }
}

function saveMockTasksDB(tasks) {
    try {
        localStorage.setItem(MOCK_TASKS_STORAGE_KEY, JSON.stringify(tasks));
    } catch (e) {
        console.error('Ошибка записи mock_tasks_db в localStorage:', e);
    }
}

function getSubjectLabel(subjectKey) {
    return SUBJECTS_DATA[subjectKey]?.label || subjectKey || 'Предмет';
}

function getTypeLabel(subjectKey, typeKey) {
    const subject = SUBJECTS_DATA[subjectKey];
    if (!subject) return typeKey || '—';
    const typeObj = subject.types.find(t => t.id === typeKey);
    return typeObj ? typeObj.label : typeKey;
}

function renderLatexIn(container) {
    if (!container) return;
    if (window.renderMathInElement) {
        renderMathInElement(container, {
            delimiters: [
                { left: '$$', right: '$$', display: true },
                { left: '$', right: '$', display: false }
            ],
            throwOnError: false
        });
    }
}
