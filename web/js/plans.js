// Функция для установки диапазона дат
function setDateRange(days) {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(endDate.getDate() - days);

    // Форматируем даты в формат YYYY-MM-DD
    document.getElementById('startDate').value = formatDate(startDate);
    document.getElementById('endDate').value = formatDate(endDate);
}

// Функция форматирования даты
function formatDate(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}

// Установка дат по умолчанию при загрузке, если не заданы
document.addEventListener('DOMContentLoaded', function () {
    const startDateInput = document.getElementById('startDate');
    const endDateInput = document.getElementById('endDate');

    if (!startDateInput.value && !endDateInput.value) {
        setDateRange(180);
    }

    document.querySelectorAll('.select-production-name').forEach(button => {
        button.addEventListener('click', function () {
            const prName = this.getAttribute('data-prname');
            const idProd = this.getAttribute('data-id');

            document.getElementById('PrName').value = prName
            document.getElementById('idProduction').value = idProd

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('productionModal'));
            if (modal) modal.hide();
        });
    })
});

