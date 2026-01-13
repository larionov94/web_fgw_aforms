// Валидация формы и навигация по вкладкам
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('productForm');
    const tabs = document.querySelectorAll('.nav-link');
    const nextButtons = document.querySelectorAll('.next-tab');
    const prevButtons = document.querySelectorAll('.prev-tab');

    // Навигация по вкладкам
    nextButtons.forEach(button => {
        button.addEventListener('click', function() {
            const nextTabId = this.getAttribute('data-next-tab');
            const nextTab = document.querySelector(`#${nextTabId}`);

            if (nextTab) {
                const tab = new bootstrap.Tab(nextTab);
                tab.show();
            }
        });
    });

    prevButtons.forEach(button => {
        button.addEventListener('click', function() {
            const prevTabId = this.getAttribute('data-prev-tab');
            const prevTab = document.querySelector(`#${prevTabId}`);

            if (prevTab) {
                const tab = new bootstrap.Tab(prevTab);
                tab.show();
            }
        });
    });

    // Валидация при отправке формы
    form.addEventListener('submit', function(e) {
        // Проверка обязательных полей в первой вкладке
        const requiredFields = document.querySelectorAll('#product [required]');
        let isValid = true;
        let firstInvalidField = null;

        requiredFields.forEach(field => {
            if (!field.value.trim()) {
                isValid = false;
                field.classList.add('is-invalid');

                // Добавляем сообщение об ошибке
                if (!field.nextElementSibling || !field.nextElementSibling.classList.contains('invalid-feedback')) {
                    const errorDiv = document.createElement('div');
                    errorDiv.className = 'invalid-feedback';
                    errorDiv.textContent = 'Это поле обязательно для заполнения';
                    field.parentNode.appendChild(errorDiv);
                }

                // Запоминаем первое невалидное поле
                if (!firstInvalidField) {
                    firstInvalidField = field;
                }
            } else {
                field.classList.remove('is-invalid');
                const errorDiv = field.parentNode.querySelector('.invalid-feedback');
                if (errorDiv) {
                    errorDiv.remove();
                }
            }
        });

        if (!isValid) {
            e.preventDefault();

            // Показываем вкладку с ошибкой
            const productTab = document.querySelector('#product-tab');
            const tab = new bootstrap.Tab(productTab);
            tab.show();

            // Скроллим к первому невалидному полю
            if (firstInvalidField) {
                firstInvalidField.scrollIntoView({ behavior: 'smooth', block: 'center' });
                firstInvalidField.focus();
            }

            // Показываем общее сообщение об ошибке
            showAlert('Пожалуйста, заполните все обязательные поля', 'danger');
        }
    });

    // Убираем ошибки при вводе
    form.querySelectorAll('input, textarea').forEach(field => {
        field.addEventListener('input', function() {
            this.classList.remove('is-invalid');
            const errorDiv = this.parentNode.querySelector('.invalid-feedback');
            if (errorDiv) {
                errorDiv.remove();
            }
        });
    });

    // Функция показа сообщений
    function showAlert(message, type) {
        // Удаляем старые алерты
        const oldAlerts = document.querySelectorAll('.alert-dismissible:not(.alert-info)');
        oldAlerts.forEach(alert => alert.remove());

        // Создаем новый алерт
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show m-3`;
        alertDiv.innerHTML = `
                    ${message}
                    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
                `;

        // Вставляем после заголовка
        const cardHeader = document.querySelector('.card-header');
        cardHeader.parentNode.insertBefore(alertDiv, cardHeader.nextSibling);
    }

    // Подтверждение при закрытии страницы с несохраненными данными
    let formChanged = false;

    form.querySelectorAll('input, textarea, select').forEach(field => {
        field.addEventListener('input', function() {
            formChanged = true;
        });

        field.addEventListener('change', function() {
            formChanged = true;
        });
    });

    window.addEventListener('beforeunload', function(e) {
        if (formChanged) {
            e.preventDefault();
            e.returnValue = 'У вас есть несохраненные изменения. Вы уверены, что хотите уйти?';
        }
    });

    // Сбрасываем флаг изменений при отправке формы
    form.addEventListener('submit', function() {
        formChanged = false;
    });

    //-------------------------------------------------------------------------------------------------------------------
    // МОДАЛЬНОЕ ОКНО----------------------------------------------------------------------------------------------------
    //-------------------------------------------------------------------------------------------------------------------
    // Обработка выбора конструкторского наименования
    document.querySelectorAll('.select-design-name').forEach(button => {
        button.addEventListener('click', function() {
            const designName = this.getAttribute('data-name');
            const designId = this.getAttribute('data-id');

            // Заполняем поле конструкторского наименования
            document.getElementById('PrName').value = designName;
            document.getElementById('PrPackName').value = designName;
            document.getElementById('PrShortName').value = designName;

            // Если нужно сохранить ID в скрытом поле
            const hiddenIdField = document.getElementById('DesignNameId');
            if (!hiddenIdField) {
                // Создаем скрытое поле для ID, если его нет
                const hiddenInput = document.createElement('input');
                hiddenInput.type = 'hidden';
                hiddenInput.id = 'DesignNameId';
                hiddenInput.name = 'DesignNameId';
                hiddenInput.value = designId;
                document.querySelector('form').appendChild(hiddenInput);
            } else {
                hiddenIdField.value = designId;
            }

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('designNameModal'));
            modal.hide();
        });
    });

    // Обработка поиска
    const searchInput = document.getElementById('designNameSearch');
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();
            const rows = document.querySelectorAll('#designNameModal tbody tr');

            rows.forEach(row => {
                // Пропускаем строку "нет данных"
                if (row.querySelector('td.text-center')) {
                    return;
                }

                const nameCell = row.querySelector('td:nth-child(2)'); // Вторая колонка (Наименование)
                const idCell = row.querySelector('td:nth-child(1)');   // Первая колонка (ID)

                if (nameCell && idCell) {
                    const name = nameCell.textContent.toLowerCase();
                    const id = idCell.textContent.toLowerCase();
                    // Ищем как по названию, так и по ID
                    const matches = name.includes(searchTerm) || id.includes(searchTerm);
                    row.style.display = matches ? '' : 'none';
                }
            });
        });
    }

    // Сброс поиска при открытии модального окна
    document.getElementById('designNameModal').addEventListener('shown.bs.modal', function() {
        if (searchInput) {
            searchInput.value = '';
            searchInput.dispatchEvent(new Event('input')); // Сбрасываем фильтрацию
            searchInput.focus(); // Фокус на поле поиска
        }
    });

    // Сброс поиска при закрытии модального окна
    document.getElementById('designNameModal').addEventListener('hidden.bs.modal', function() {
        if (searchInput) {
            searchInput.value = '';
            searchInput.dispatchEvent(new Event('input')); // Сбрасываем фильтрацию
        }
    });



    // ========== ОБРАБОТКА ЦВЕТА СТЕКЛА ==========

    // Обработка выбора цвета
    document.querySelectorAll('.select-color').forEach(button => {
        button.addEventListener('click', function() {
            const colorName = this.getAttribute('data-name');
            const colorGL = this.getAttribute('data-gl');


            // Заполняем поле "Цвет стекла"
            document.getElementById('PrColor').value = colorName;
            document.getElementById('PrGL').value = colorGL;

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('colorModal'));
            modal.hide();
        });
    });
});
