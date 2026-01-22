document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('productForm');
    const requiredFields = form.querySelectorAll('[required]');

    // Добавляем поле артикула в валидацию
    const articleField = document.getElementById('PrArticle');

    // Флаг для отслеживания изменений формы
    let formChanged = false;

    // ========== ВАЛИДАЦИЯ ФОРМЫ ==========

    // Функция проверки всех обязательных полей
    function validateForm() {
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
    function validateTab(tabPane) {
        const requiredInTab = tabPane.querySelectorAll('[required]');
        let hasErrors = false;
        let firstInvalidField = null;

        requiredInTab.forEach(field => {
            field.classList.remove('is-invalid');
            const oldError = field.parentNode.querySelector('.invalid-feedback');
            if (oldError) oldError.remove();

            const value = field.value.trim();

            // Специальная проверка для артикула
            if (field.id === 'PrArticle') {
                if (value.length !== 2) {
                    hasErrors = true;
                    field.classList.add('is-invalid');

                    const errorDiv = document.createElement('div');
                    errorDiv.className = 'invalid-feedback p-0';
                    errorDiv.textContent = `Артикул должен содержать ровно 2 символа`;
                    field.parentNode.appendChild(errorDiv);

                    if (!firstInvalidField) {
                        firstInvalidField = field;
                    }
                }
                return; // Выходим из обработки
            }

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
    form.addEventListener('submit', function (e) {
        const {isValid, firstInvalidField} = validateForm();

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
            showAlert('Пожалуйста, заполните все обязательные поля', 'danger');

            return false;
        }

        // Сбрасываем флаг изменений при успешной отправке
        formChanged = false;
    });

    form.addEventListener('keypress', function(e) {
        if (e.key === 'Enter' && e.target.matches('input:not([type="button"]):not([type="submit"])')) {
            e.preventDefault();

            if (e.target.tagName === 'TEXTAREA') {
                e.target.value += '\n';
            }

            const activeTab = document.querySelector('.tab-pane.active');
            if (!activeTab) return;

            // Проверяем валидацию текущей вкладки
            const { hasErrors } = validateTab(activeTab);

            if (!hasErrors) {
                if (activeTab.id === 'additionally') {
                    document.getElementById('saveProductionBtn1')?.click();
                } else {
                    activeTab.querySelector('.next-tab')?.click();
                }
            } else {
                // Если есть ошибки, показываем сообщение
                showAlert('Заполните обязательные поля перед переходом дальше', 'warning');
            }
        }
    });

    // Убираем ошибки при вводе
    form.querySelectorAll('input, textarea, select').forEach(field => {
        field.addEventListener('input', function () {
            this.classList.remove('is-invalid');
            const errorDiv = this.parentNode.querySelector('.invalid-feedback');
            if (errorDiv) errorDiv.remove();

            // Автоматически скрываем алерт при изменении
            hideAlert();

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
    const nextButtons = document.querySelectorAll('.next-tab');
    nextButtons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault();

            const currentTabPane = this.closest('.tab-pane');
            const {hasErrors, firstInvalidField} = validateTab(currentTabPane);

            if (hasErrors) {
                // Показываем сообщение об ошибке
                showAlert('Пожалуйста, заполните все обязательные поля в текущей вкладке.', 'warning');

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
    const prevButtons = document.querySelectorAll('.prev-tab');
    prevButtons.forEach(button => {
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

    // Функция показа сообщений
    function showAlert(message, type) {
        // Удаляем старые алерты
        const oldAlerts = document.querySelectorAll('.alert-dismissible:not(.alert-info)');
        oldAlerts.forEach(alert => alert.remove());

        // Создаем новый алерт
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show m-1 p-1 small`;
        alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close btn-sm m-1 p-1" data-bs-dismiss="alert" aria-label="Close"></button>
    `;

        // Сохраняем ссылку на текущий алерт
        currentAlert = alertDiv;

        // Вставляем после заголовка
        const cardHeader = document.querySelector('.card-header');
        if (cardHeader) {
            cardHeader.parentNode.insertBefore(alertDiv, cardHeader.nextSibling);
        }
    }

    // Функция скрытия алерта
    function hideAlert() {
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

    // ========== МОДАЛЬНЫЕ ОКНА ==========

    // Обработка выбора конструкторского наименования
    document.querySelectorAll('.select-design-name').forEach(button => {
        button.addEventListener('click', function () {
            const designName = this.getAttribute('data-name');
            const designId = this.getAttribute('data-id');

            // Заполняем поля
            document.getElementById('PrName').value = designName;
            document.getElementById('PrPackName').value = designName;
            document.getElementById('PrShortName').value = designName;

            // Убираем валидационные ошибки
            ['PrName', 'PrPackName', 'PrShortName'].forEach(id => {
                const field = document.getElementById(id);
                if (field) {
                    field.classList.remove('is-invalid');
                    const errorDiv = field.parentNode.querySelector('.invalid-feedback');
                    if (errorDiv) errorDiv.remove();
                }
            });

            // Сохраняем ID в скрытом поле
            let hiddenIdField = document.getElementById('DesignNameId');
            if (!hiddenIdField) {
                hiddenIdField = document.createElement('input');
                hiddenIdField.type = 'hidden';
                hiddenIdField.id = 'DesignNameId';
                hiddenIdField.name = 'DesignNameId';
                document.querySelector('form').appendChild(hiddenIdField);
            }
            hiddenIdField.value = designId;

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('designNameModal'));
            if (modal) modal.hide();

            // Отмечаем изменение формы
            formChanged = true;
        });
    });

    // Обработка выбора цвета стекла
    document.querySelectorAll('.select-color').forEach(button => {
        button.addEventListener('click', function () {
            const colorName = this.getAttribute('data-name');
            const colorGL = this.getAttribute('data-gl');

            // Заполняем поля
            document.getElementById('PrColor').value = colorName;
            document.getElementById('PrGL').value = colorGL;

            // Убираем валидационные ошибки
            ['PrColor', 'PrGL'].forEach(id => {
                const field = document.getElementById(id);
                if (field) {
                    field.classList.remove('is-invalid');
                    const errorDiv = field.parentNode.querySelector('.invalid-feedback');
                    if (errorDiv) errorDiv.remove();
                }
            });

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('colorModal'));
            if (modal) modal.hide();

            // Отмечаем изменение формы
            formChanged = true;
        });
    });

    // ========== ПОИСК В МОДАЛЬНЫХ ОКНАХ ==========

    // Обработка поиска конструкторского наименования
    const designNameSearch = document.getElementById('designNameSearch');
    if (designNameSearch) {
        designNameSearch.addEventListener('input', function () {
            const searchTerm = this.value.toLowerCase();
            const rows = document.querySelectorAll('#designNameModal tbody tr');

            rows.forEach(row => {
                // Пропускаем строку "нет данных"
                if (row.querySelector('td.text-center')) {
                    return;
                }

                const nameCell = row.querySelector('td:nth-child(2)');
                const idCell = row.querySelector('td:nth-child(1)');

                if (nameCell && idCell) {
                    const name = nameCell.textContent.toLowerCase();
                    const id = idCell.textContent.toLowerCase();
                    const matches = name.includes(searchTerm) || id.includes(searchTerm);
                    row.style.display = matches ? '' : 'none';
                }
            });
        });
    }

    // Сброс поиска при открытии/закрытии модальных окон
    const designNameModal = document.getElementById('designNameModal');
    if (designNameModal) {
        designNameModal.addEventListener('shown.bs.modal', function () {
            if (designNameSearch) {
                designNameSearch.value = '';
                designNameSearch.dispatchEvent(new Event('input'));
                designNameSearch.focus();
            }
        });

        designNameModal.addEventListener('hidden.bs.modal', function () {
            if (designNameSearch) {
                designNameSearch.value = '';
                designNameSearch.dispatchEvent(new Event('input'));
            }
        });
    }

    // ========== РАСЧЕТ ВЕСА ==========

    // Функция расчета
    function calculateTotal() {
        const prCount = document.getElementById('PrCount');
        const prRows = document.getElementById('PrRows');
        const weightBottle = document.getElementById('WeightBottle');
        const weightPallet = document.getElementById('WeightPallet');
        const weightTotal = document.getElementById('WeightTotal');
        const weightCommon = document.getElementById('PrWeight');

        if (!prCount || !prRows || !weightBottle || !weightPallet || !weightTotal) {
            return; // Если поля не найдены, выходим
        }

        const count = parseFloat(prCount.value) || 0;
        const rows = parseFloat(prRows.value) || 0;
        const bottleWeight = parseFloat(weightBottle.value) || 0;
        const palletWeight = parseFloat(weightPallet.value) || 45;

        const total = (count * rows * bottleWeight) + palletWeight;
        weightTotal.value = total.toFixed(3);
        weightCommon.value = total.toFixed(0)
    }

    // Вешаем обработчики на все поля для расчета
    ['PrCount', 'PrRows', 'WeightBottle', 'WeightPallet'].forEach(id => {
        const field = document.getElementById(id);
        if (field) {
            field.addEventListener('input', calculateTotal);
            field.addEventListener('change', calculateTotal);
        }
    });

    // Вызываем первый расчет
    calculateTotal();

    // Функция сборки HWD
    function updateHWD() {
        const height = document.getElementById('PrHeight').value.trim();
        const width = document.getElementById('PrWidth').value.trim();
        const depth = document.getElementById('PrDepth').value.trim();

        const hwdField = document.getElementById('PrHWD');

        // Собираем только если все поля заполнены
        if (height && width && depth) {
            hwdField.value = `${height}x${width}x${depth}`;
        } else {
            hwdField.value = ''; // Очищаем если не все заполнено
        }
    }

    // Вешаем обработчики на все 3 поля
    ['PrHeight', 'PrWidth', 'PrDepth'].forEach(id => {
        const field = document.getElementById(id);
        if (field) {
            field.addEventListener('input', updateHWD);
            field.addEventListener('change', updateHWD);
        }
    });

    // Инициализация при загрузке
    updateHWD();
});