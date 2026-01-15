document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('productForm');
    const requiredFields = form.querySelectorAll('[required]');

    // Функция для проверки и подсветки полей
    function validateForm() {
        let isValid = true;

        // Сначала убираем все подсветки
        requiredFields.forEach(field => {
            field.classList.remove('is-invalid');
        });

        // Проверяем каждое обязательное поле
        requiredFields.forEach(field => {
            if (!field.value.trim()) {
                field.classList.add('is-invalid');
                isValid = false;

                // Показываем вкладку с ошибкой
                const tabPane = field.closest('.tab-pane');
                if (tabPane && !tabPane.classList.contains('active')) {
                    const tabId = tabPane.id;
                    const tabButton = document.querySelector(`[data-bs-target="#${tabId}"]`);
                    if (tabButton) {
                        tabButton.click();
                    }
                }
            }
        });

        return isValid;
    }

    // Обработчик отправки формы
    form.addEventListener('submit', function (e) {
        if (!validateForm()) {
            e.preventDefault();

            // Показываем уведомление
            alert('Пожалуйста, заполните все обязательные поля, отмеченные красным цветом.');
        }
    });

    // Убираем красную подсветку при вводе
    requiredFields.forEach(field => {
        field.addEventListener('input', function () {
            if (this.value.trim()) {
                this.classList.remove('is-invalid');
            }
        });

        field.addEventListener('change', function () {
            if (this.value.trim()) {
                this.classList.remove('is-invalid');
            }
        });
    });

    // Навигация по вкладкам с проверкой
    const nextButtons = document.querySelectorAll('.next-tab');
    nextButtons.forEach(button => {
        button.addEventListener('click', function (e) {
            // Сразу предотвращаем все стандартные действия
            e.preventDefault();
            e.stopPropagation();
            e.stopImmediatePropagation(); // Полностью останавливаем событие

            const nextTabId = this.getAttribute('data-next-tab');
            const currentTabPane = this.closest('.tab-pane');

            // Проверяем обязательные поля в текущей вкладке
            const requiredInCurrentTab = currentTabPane.querySelectorAll('[required]');
            let hasErrors = false;

            requiredInCurrentTab.forEach(field => {
                if (!field.value.trim()) {
                    field.classList.add('is-invalid');
                    hasErrors = true;
                }
            });

            if (hasErrors) {
                alert('Пожалуйста, заполните все обязательные поля в текущей вкладке.');

                // Фокус на первое незаполненное поле
                const firstInvalid = currentTabPane.querySelector('.is-invalid');
                if (firstInvalid) {
                    firstInvalid.focus();
                }
                // НЕ переключаем вкладку - просто выходим
                return false;
            } else {
                // Только если нет ошибок - переключаем вкладку
                const nextTabButton = document.getElementById(nextTabId);
                if (nextTabButton) {
                    const tab = new bootstrap.Tab(nextTabButton);
                    tab.show();
                }
            }
        }, true); // Используем capture phase для приоритета
    });

    // Обработка кнопок "Назад" - убираем подсветку
    const prevButtons = document.querySelectorAll('.prev-tab');
    prevButtons.forEach(button => {
        button.addEventListener('click', function () {
            const prevTabId = this.getAttribute('data-prev-tab');
            const prevTab = document.getElementById(prevTabId);
            if (prevTab) {
                const tabInstance = new bootstrap.Tab(prevTab);
                tabInstance.show();
            }
        });
    });

    // Модификация модальных окон для работы с required полями
    const designNameButtons = document.querySelectorAll('.select-design-name');
    designNameButtons.forEach(button => {
        button.addEventListener('click', function () {
            const designNameInput = document.getElementById('PrPackName');
            designNameInput.value = this.getAttribute('data-name');
            designNameInput.classList.remove('is-invalid');

            const modal = bootstrap.Modal.getInstance(document.getElementById('designNameModal'));
            if (modal) modal.hide();
        });
    });

    const colorButtons = document.querySelectorAll('.select-color');
    colorButtons.forEach(button => {
        button.addEventListener('click', function () {
            const colorInput = document.getElementById('PrColor');
            const glInput = document.getElementById('PrGL');

            colorInput.value = this.getAttribute('data-name');
            glInput.value = this.getAttribute('data-gl');

            colorInput.classList.remove('is-invalid');
            glInput.classList.remove('is-invalid');

            const modal = bootstrap.Modal.getInstance(document.getElementById('colorModal'));
            if (modal) modal.hide();
        });
    });
});