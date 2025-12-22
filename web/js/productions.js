document.addEventListener('DOMContentLoaded', function() {
    const searchArticle = document.getElementById('searchArticle');
    const searchName = document.getElementById('searchName');
    const searchIdProduction = document.getElementById('searchIdProduction'); // Новое поле
    const searchForm = document.getElementById('searchForm');

    // Очистка при клике на крестик
    const clearSearchIconArticle = document.getElementById('clearSearchIconArticle');
    const clearSearchIconName = document.getElementById('clearSearchIconName');
    const clearSearchIconIdProduction = document.getElementById('clearSearchIconIdProduction'); // Новый крестик

    let debounceTimer;

    // Очистка при клике на крестик артикула
    if (clearSearchIconArticle) {
        clearSearchIconArticle.addEventListener('click', function() {
            searchArticle.value = '';
            searchForm.submit();
        });
    }

    // Очистка при клике на крестик наименования
    if (clearSearchIconName) {
        clearSearchIconName.addEventListener('click', function() {
            searchName.value = '';
            searchForm.submit();
        });
    }

    // Очистка при клике на крестик ID продукции (новое)
    if (clearSearchIconIdProduction) {
        clearSearchIconIdProduction.addEventListener('click', function() {
            searchIdProduction.value = '';
            searchForm.submit();
        });
    }

    // 1. Авто-поиск при вводе в артикул
    searchArticle.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });

    // 1. Авто-поиск при вводе в наименование
    searchName.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });

    // 1. Авто-поиск при вводе в ID продукции (новое)
    searchIdProduction.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });

    // 4. Фокус на поле поиска если есть поисковый запрос
    // Приоритет: ID продукции → Наименование → Артикул
    if (searchIdProduction && searchIdProduction.value) {
        searchIdProduction.focus();
        // 5. Помещаем курсор в конец
        searchIdProduction.setSelectionRange(searchIdProduction.value.length, searchIdProduction.value.length);
    } else if (searchName && searchName.value) {
        searchName.focus();
        // 5. Помещаем курсор в конец
        searchName.setSelectionRange(searchName.value.length, searchName.value.length);
    } else if (searchArticle && searchArticle.value) {
        searchArticle.focus();
        // 5. Помещаем курсор в конец
        searchArticle.setSelectionRange(searchArticle.value.length, searchArticle.value.length);
    } else if (searchArticle) {
        searchArticle.focus(); // По умолчанию фокус на первом поле
    }
});