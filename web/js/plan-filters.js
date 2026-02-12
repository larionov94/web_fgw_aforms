//--------------------------------------------------------//
//              УСТАНОВКА ФОРМАТА ДАТЫ                    //
//--------------------------------------------------------//

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

//--------------------------------------------------------//
//              ОБНОВЛЕНИЕ СЧЕТЧИКА СТРОК                 //
//--------------------------------------------------------//

// Функция для обновления счетчика строк
function updateRowCounter() {
    const visibleRows = document.querySelectorAll('#productionTableBody tr:not([style*="display: none"])').length;
    const totalRows = document.querySelectorAll('#productionTableBody tr').length;
    const noDataRow = document.querySelector('#productionTableBody tr td.text-center[colspan]');

    // Если есть строка "нет данных", не считаем ее
    let actualVisibleRows = visibleRows;
    let actualTotalRows = totalRows;

    if (noDataRow) {
        const isNoDataVisible = !noDataRow.closest('tr').style.display ||
            noDataRow.closest('tr').style.display === '';
        actualTotalRows = totalRows - 1; // Исключаем строку "нет данных" из общего числа
        if (isNoDataVisible) {
            actualVisibleRows = 0; // Если видна строка "нет данных", значит обычных строк 0
        }
    }

    // Обновляем счетчики
    const visibleCounter = document.getElementById('visibleRowCount');
    const totalCounter = document.getElementById('totalRowCount');

    if (visibleCounter) {
        visibleCounter.textContent = actualVisibleRows;
    }
    if (totalCounter) {
        totalCounter.textContent = actualTotalRows;
    }
}

//--------------------------------------------------------//
//              ОБРАБОТКА ПОИСКА                          //
//--------------------------------------------------------//

// Функция поиска по таблице
function initProductionSearch(searchElementId, tableBodyId) {
    const searchElement = document.getElementById(searchElementId);
    if (!searchElement) return;

    searchElement.addEventListener('input', function () {
        const searchTerm = this.value.toLowerCase();
        const rows = document.querySelectorAll(`#${tableBodyId} tr`);

        rows.forEach(row => {
            // Пропускаем строку "нет данных"
            if (row.querySelector('td.text-center[colspan]')) {
                // Показываем только если поиск пустой И нет других строк
                const hasOtherVisibleRows = Array.from(rows).some(r => {
                    if (r === row) return false;
                    const noDataCell = r.querySelector('td.text-center[colspan]');
                    return !noDataCell &&
                        (!r.style.display || r.style.display === '');
                });
                row.style.display = (searchTerm === '' && !hasOtherVisibleRows) ? '' : 'none';
                return;
            }

            const articleCell = row.querySelector('td:nth-child(1)'); // Артикул
            const nameCell1 = row.querySelector('td:nth-child(2)'); // Наименование для упаковщика
            const nameCell2 = row.querySelector('td:nth-child(3)'); // Наименование для этикетки
            const colorCell = row.querySelector('td:nth-child(4)'); // Цвет
            const codeCell = row.querySelector('td:nth-child(7)'); // Код продукции

            let matches = false;

            // Проверяем все поля
            if (articleCell && articleCell.textContent.toLowerCase().includes(searchTerm)) matches = true;
            if (nameCell1 && nameCell1.textContent.toLowerCase().includes(searchTerm)) matches = true;
            if (nameCell2 && nameCell2.textContent.toLowerCase().includes(searchTerm)) matches = true;
            if (colorCell && colorCell.textContent.toLowerCase().includes(searchTerm)) matches = true;
            if (codeCell && codeCell.textContent.toLowerCase().includes(searchTerm)) matches = true;

            row.style.display = matches ? '' : 'none';
        });
        updateRowCounter();
    });
}

//--------------------------------------------------------//
//              ИНИЦИАЛИЗАЦИЯ ПРИ ЗАГРУЗКЕ                //
//--------------------------------------------------------//

document.addEventListener('DOMContentLoaded', function () {
    // --- Установка дат ---
    const startDateInput = document.getElementById('startDate');
    const endDateInput = document.getElementById('endDate');
    const planDateInput = document.getElementById('PlanDate');

    if (planDateInput) {
        planDateInput.value = formatDate(new Date());
    }

    if (startDateInput && endDateInput && !startDateInput.value && !endDateInput.value) {
        setDateRange(180);
    }

    // --- Обработка выбора производства (основной и дополнительный) ---
    document.querySelectorAll('.select-production-name, .select-production-name-add').forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault();

            const prName = this.getAttribute('data-prname');
            const idProd = this.getAttribute('data-id');

            // Определяем, какие поля заполнять
            if (this.classList.contains('select-production-name')) {
                document.getElementById('PrName').value = prName;
                document.getElementById('idProduction').value = idProd;

                // Закрываем основное модальное окно
                const modal = bootstrap.Modal.getInstance(document.getElementById('productionModal'));
                if (modal) modal.hide();
            } else if (this.classList.contains('select-production-name-add')) {
                document.getElementById('PrName').value = prName;
                document.getElementById('extProduction').value = idProd;

                // Закрываем дополнительное модальное окно
                const modal = bootstrap.Modal.getInstance(document.getElementById('productionAddModal'));
                if (modal) modal.hide();
            }
        });
    });

    // --- Обработка выбора сектора (основной и дополнительный) ---
    document.querySelectorAll('.select-sector, .select-sector-add').forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault();

            const sectorName = this.getAttribute('data-name');
            const idSec = this.getAttribute('data-id');

            // Определяем, какие поля заполнять
            if (this.classList.contains('select-sector')) {
                document.getElementById('SectorName').value = sectorName;
                document.getElementById('idSector').value = idSec;

                // Закрываем основное модальное окно
                const modal = bootstrap.Modal.getInstance(document.getElementById('sectorModal'));
                if (modal) modal.hide();
            } else if (this.classList.contains('select-sector-add')) {
                document.getElementById('SectorName').value = sectorName;
                document.getElementById('extSector').value = idSec;

                // Закрываем дополнительное модальное окно
                const modal = bootstrap.Modal.getInstance(document.getElementById('sectorAddModal'));
                if (modal) modal.hide();
            }
        });
    });

    // --- Инициализация поиска для основного модального окна ---
    initProductionSearch('productionSearch', 'productionTableBody');

    // --- Инициализация поиска для дополнительного модального окна ---
    initProductionSearch('productionSearch', 'productionTableBody');

    // --- Обработка модальных окон (сброс поиска) ---

    // Основное модальное окно производства
    const productionModal = document.getElementById('productionModal');
    if (productionModal) {
        productionModal.addEventListener('shown.bs.modal', function () {
            const searchInput = document.getElementById('productionSearch');
            if (searchInput) {
                searchInput.value = '';
                searchInput.placeholder = "Поиск по артикулу, наименованию...";
                searchInput.dispatchEvent(new Event('input'));
                searchInput.focus();
            }
            updateRowCounter();
        });

        productionModal.addEventListener('hidden.bs.modal', function () {
            const searchInput = document.getElementById('productionSearch');
            if (searchInput) {
                searchInput.value = '';
                searchInput.dispatchEvent(new Event('input'));
            }
        });
    }

    // Дополнительное модальное окно производства
    const productionAddModal = document.getElementById('productionAddModal');
    if (productionAddModal) {
        productionAddModal.addEventListener('shown.bs.modal', function () {
            const searchInput = document.getElementById('productionSearch');
            if (searchInput) {
                searchInput.value = '';
                searchInput.placeholder = "Поиск по артикулу, наименованию...";
                searchInput.dispatchEvent(new Event('input'));
                searchInput.focus();
            }
            updateRowCounter();
        });

        productionAddModal.addEventListener('hidden.bs.modal', function () {
            const searchInput = document.getElementById('productionSearch');
            if (searchInput) {
                searchInput.value = '';
                searchInput.dispatchEvent(new Event('input'));
            }
        });
    }

    // --- Первоначальное обновление счетчика ---
    updateRowCounter();
});

// Альтернативный вариант: если кнопки добавляются динамически,
// оставьте один глобальный обработчик для делегирования
document.addEventListener('click', function (e) {
    // Обработка выбора производства (дополнительный вариант для динамических элементов)
    if (e.target.classList.contains('select-production-name-add')) {
        e.preventDefault();

        const prName = e.target.getAttribute('data-prname');
        const idProd = e.target.getAttribute('data-id');

        document.getElementById('PrName').value = prName;
        document.getElementById('extProduction').value = idProd;

        const modal = bootstrap.Modal.getInstance(document.getElementById('productionAddModal'));
        if (modal) modal.hide();
    }

    // Обработка выбора сектора (дополнительный вариант для динамических элементов)
    if (e.target.classList.contains('select-sector-add')) {
        e.preventDefault();

        const sectorName = e.target.getAttribute('data-name');
        const sectorId = e.target.getAttribute('data-id');

        document.getElementById('SectorName').value = sectorName;
        document.getElementById('extSector').value = sectorId;

        const modal = bootstrap.Modal.getInstance(document.getElementById('sectorAddModal'));
        if (modal) modal.hide();
    }
});