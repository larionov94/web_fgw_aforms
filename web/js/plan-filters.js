//--------------------------------------------------------//
//              УСТАНОВКА ФОРМАТА ДАТЫ                    //
//--------------------------------------------------------//
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

// Установка дат по умолчанию при загрузке, если не заданы
document.addEventListener('DOMContentLoaded', function () {
    const startDateInput = document.getElementById('startDate');
    const endDateInput = document.getElementById('endDate');

    const planDateInput = document.getElementById('PlanDate');

    planDateInput.value = formatDate(new Date());

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
    });

    document.querySelectorAll('.select-sector').forEach(button => {
        button.addEventListener('click', function () {
            const SectorName = this.getAttribute('data-name');
            const idSec = this.getAttribute('data-id');

            document.getElementById('SectorName').value = SectorName
            document.getElementById('idSector').value = idSec

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('sectorModal'));
            if (modal) modal.hide();
        });
    });

    //--------------------------------------------------------//
    //     ОБРАБОТКА ПОИСКА ВАРИАНТА УПАКОВКИ ПРОДУКЦИИ       //
    //                  МОДАЛЬНОЕ ОКНО                        //
    //--------------------------------------------------------//

    // Обработка поиска конструкторского наименования
    const productionSearch = document.getElementById('productionSearch');
    if (productionSearch) {
        productionSearch.addEventListener('input', function () {
            const searchTerm = this.value.toLowerCase();
            const rows = document.querySelectorAll('#productionTableBody tr');

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

    // Сброс поиска при открытии/закрытии модальных окон
    const productionModal = document.getElementById('productionModal');
    if (productionModal) {
        productionModal.addEventListener('shown.bs.modal', function () {
            if (productionSearch) {
                productionSearch.value = '';
                productionSearch.placeholder = "Поиск по артикулу, наименованию...";
                productionSearch.dispatchEvent(new Event('input'));
                productionSearch.focus();
            }
            updateRowCounter();
        });

        productionModal.addEventListener('hidden.bs.modal', function () {
            if (productionSearch) {
                productionSearch.value = '';
                productionSearch.dispatchEvent(new Event('input'));
            }
        });
    }
    updateRowCounter();

    //--------------------------------------------------------//
    //              ОБНОВЛЕНИЕ СЧЕТЧИКА СТРОК В               //
    //                    МОДАЛЬНОМ ОКНЕ                      //
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
});

