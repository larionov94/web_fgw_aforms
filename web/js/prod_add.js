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
});