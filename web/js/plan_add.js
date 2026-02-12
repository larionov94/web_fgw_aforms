document.addEventListener('DOMContentLoaded', function () {
    const formPlan = document.getElementById('planForm');
    if (!formPlan) return; // Защита от ошибки

    const requiredFields = formPlan.querySelectorAll('[required]');

    // Флаг для отслеживания изменений формы
    let formChanged = false;

    // ========== ВАЛИДАЦИЯ ФОРМЫ ==========

    // Функция проверки всех обязательных полей
    function validatePlanForm() {
        let isValid = true;
        let firstInvalidField = null;

        requiredFields.forEach(field => {
            // Удаляем старые сообщения об ошибках
            field.classList.remove('is-invalid');
            const oldError = field.parentNode.querySelector('.invalid-feedback');
            if (oldError) oldError.remove();

            // Проверяем поле
            if (!field.value.trim()) {
                isValid = false;
                field.classList.add('is-invalid');

                // Создаем сообщение об ошибке
                const errorDiv = document.createElement('div');
                errorDiv.className = 'invalid-feedback p-0';
                errorDiv.textContent = 'Это поле обязательно для заполнения';
                field.parentNode.appendChild(errorDiv);

                if (!firstInvalidField) {
                    firstInvalidField = field;
                }
            }
        });

        return {isValid, firstInvalidField};
    }

    // Функция проверки конкретной вкладки
    function validatePlanTab(tabPane) {
        const requiredInTab = tabPane.querySelectorAll('[required]');
        let hasErrors = false;
        let firstInvalidField = null;

        requiredInTab.forEach(field => {
            field.classList.remove('is-invalid');
            const oldError = field.parentNode.querySelector('.invalid-feedback');
            if (oldError) oldError.remove();

            const value = field.value.trim();

            if (!value) {
                hasErrors = true;
                field.classList.add('is-invalid');

                const errorDiv = document.createElement('div');
                errorDiv.className = 'invalid-feedback p-0';
                errorDiv.textContent = 'Это поле обязательно для заполнения';
                field.parentNode.appendChild(errorDiv);

                if (!firstInvalidField) {
                    firstInvalidField = field;
                }
            }
        });

        return {hasErrors, firstInvalidField};
    }

    // Валидация при отправке формы
    formPlan.addEventListener('submit', function (e) {
        const {isValid, firstInvalidField} = validatePlanForm();

        if (!isValid) {
            e.preventDefault();

            // Показываем первую вкладку с ошибкой
            if (firstInvalidField) {
                const tabPane = firstInvalidField.closest('.tab-pane');
                if (tabPane) {
                    const tabId = tabPane.id;
                    const tabButton = document.querySelector(`[data-bs-target="#${tabId}"]`);
                    if (tabButton) {
                        const tab = new bootstrap.Tab(tabButton);
                        tab.show();
                    }
                }

                firstInvalidField.scrollIntoView({behavior: 'smooth', block: 'center'});
                firstInvalidField.focus();
            }

            // Показываем общее сообщение об ошибке
            showPlanAlert('Пожалуйста, заполните все обязательные поля', 'danger');

            return false;
        }

        // Сбрасываем флаг изменений при успешной отправке
        formChanged = false;
    });

    formPlan.addEventListener('keypress', function (e) {
        if (e.key === 'Enter' && e.target.matches('input:not([type="button"]):not([type="submit"])')) {
            e.preventDefault();

            if (e.target.tagName === 'TEXTAREA') {
                e.target.value += '\n';
            }

            const activeTab = document.querySelector('.tab-pane.active');
            if (!activeTab) return;

            // Проверяем валидацию текущей вкладки
            const {hasErrors} = validatePlanTab(activeTab);

            if (!hasErrors) {
                if (activeTab.id === 'info') {
                    document.getElementById('savePlanBtn1')?.click();
                } else {
                    activeTab.querySelector('.next-tab')?.click();
                }
            } else {
                // Если есть ошибки, показываем сообщение
                showPlanAlert('Заполните обязательные поля перед переходом дальше', 'warning');
            }
        }
    });

    // Убираем ошибки при вводе
    formPlan.querySelectorAll('input, textarea, select').forEach(field => {
        field.addEventListener('input', function () {
            this.classList.remove('is-invalid');
            const errorDiv = this.parentNode.querySelector('.invalid-feedback');
            if (errorDiv) errorDiv.remove();

            // Автоматически скрываем алерт при изменении
            hidePlanAlert();

            // Отмечаем, что форма была изменена
            formChanged = true;
        });

        field.addEventListener('change', function () {
            this.classList.remove('is-invalid');
            const errorDiv = this.parentNode.querySelector('.invalid-feedback');
            if (errorDiv) errorDiv.remove();

            // Отмечаем, что форма была изменена
            formChanged = true;
        });
    });

    // ========== НАВИГАЦИЯ ПО ВКЛАДКАМ ==========

    // Обработка кнопок "Далее"
    const nextPlanButtons = document.querySelectorAll('.next-tab');
    nextPlanButtons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault();

            const currentTabPane = this.closest('.tab-pane');
            if (!currentTabPane) return;

            const {hasErrors, firstInvalidField} = validatePlanTab(currentTabPane);

            if (hasErrors) {
                // Показываем сообщение об ошибке
                showPlanAlert('Пожалуйста, заполните все обязательные поля в текущей вкладке.', 'warning');

                // Фокус на первое незаполненное поле
                if (firstInvalidField) {
                    firstInvalidField.focus();
                }

                return false;
            }

            // Переключаем на следующую вкладку
            const nextTabId = this.getAttribute('data-next-tab');
            const nextTabButton = document.getElementById(nextTabId);
            if (nextTabButton) {
                const tab = new bootstrap.Tab(nextTabButton);
                tab.show();
            }
        });
    });

    // Обработка кнопок "Назад"
    const prevPlanButtons = document.querySelectorAll('.prev-tab');
    prevPlanButtons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault();

            const prevTabId = this.getAttribute('data-prev-tab');
            const prevTabButton = document.getElementById(prevTabId);
            if (prevTabButton) {
                const tab = new bootstrap.Tab(prevTabButton);
                tab.show();
            }
        });
    });

    // ========== ФУНКЦИИ ВСПОМОГАТЕЛЬНЫЕ ==========
    let currentAlert = null;

    // Функция показа сообщений - ИСПРАВЛЕНО
    function showPlanAlert(message, type) {
        // Удаляем старые алерты
        const oldAlerts = document.querySelectorAll('.alert-dismissible:not(.alert-info)');
        oldAlerts.forEach(alert => alert.remove());

        // Создаем новый алерт
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show m-2 p-2 small`;
        alertDiv.setAttribute('role', 'alert');
        alertDiv.innerHTML = `
            <div class="d-flex justify-content-between align-items-center">
                <span>${message}</span>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
        `;

        // Сохраняем ссылку на текущий алерт
        currentAlert = alertDiv;

        // ИСПРАВЛЕНО: Ищем .card вместо .card-header
        const card = document.querySelector('.card');
        if (card) {
            // Вставляем алерт в начало карточки
            card.prepend(alertDiv);
        } else {
            // Если нет карточки, вставляем после формы
            formPlan.parentNode.insertBefore(alertDiv, formPlan);
        }

        // Автоматически скрываем через 5 секунд
        setTimeout(() => {
            if (currentAlert === alertDiv) {
                hidePlanAlert();
            }
        }, 5000);
    }

    // Функция скрытия алерта
    function hidePlanAlert() {
        if (currentAlert) {
            const closeButton = currentAlert.querySelector('.btn-close');
            if (closeButton) {
                closeButton.click();
            }
            currentAlert = null;
        }
    }

    // ========== ПОДТВЕРЖДЕНИЕ ПРИ ЗАКРЫТИИ СТРАНИЦЫ ==========
    window.addEventListener('beforeunload', function (e) {
        if (formChanged) {
            e.preventDefault();
            e.returnValue = 'У вас есть несохраненные изменения. Вы уверены, что хотите уйти?';
        }
    });



    // ========== ОБРАБОТЧИКИ ДЛЯ МОДАЛЬНЫХ ОКОН ==========

    // Выбор участка
    document.querySelectorAll('.select-sector-add').forEach(button => {
        button.addEventListener('click', function() {
            const sectorName = this.getAttribute('data-name');
            const sectorId = this.getAttribute('data-id');

            const sectorField = document.getElementById('SectorName');
            const extSector = document.getElementById('extSector');

            if (sectorField) {
                sectorField.value = sectorName;
                // Убираем красную рамку и сообщение об ошибке
                sectorField.classList.remove('is-invalid');
                const errorDiv = sectorField.parentNode.querySelector('.invalid-feedback');
                if (errorDiv) errorDiv.remove();
            }

            if (extSector) extSector.value = sectorId;

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('sectorAddModal'));
            if (modal) modal.hide();

            // Отмечаем изменение формы
            formChanged = true;
        });
    });

    // Выбор продукции
    document.querySelectorAll('.select-production-name-add').forEach(button => {
        button.addEventListener('click', function() {
            const prName = this.getAttribute('data-prname');
            const productionId = this.getAttribute('data-id');

            const prNameField = document.getElementById('PrName');
            const extProduction = document.getElementById('extProduction');

            if (prNameField) {
                prNameField.value = prName;
                // Убираем красную рамку и сообщение об ошибке
                prNameField.classList.remove('is-invalid');
                const errorDiv = prNameField.parentNode.querySelector('.invalid-feedback');
                if (errorDiv) errorDiv.remove();
            }

            if (extProduction) extProduction.value = productionId;

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('productionAddModal'));
            if (modal) modal.hide();

            // Отмечаем изменение формы
            formChanged = true;
        });
    });

    // ========== УСТАНОВКА ДАТЫ ПО УМОЛЧАНИЮ ==========
    const planDateField = document.getElementById('PlanDate');
    if (planDateField && !planDateField.value) {
        const today = new Date();
        const yyyy = today.getFullYear();
        const mm = String(today.getMonth() + 1).padStart(2, '0');
        const dd = String(today.getDate()).padStart(2, '0');
        planDateField.value = `${yyyy}-${mm}-${dd}`;
    }

});