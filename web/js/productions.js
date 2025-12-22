document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const searchForm = document.getElementById('searchForm');
    let debounceTimer;

    // 1. Авто-поиск при вводе
    searchInput.addEventListener('input', function(e) {
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
    if (searchInput.value) {
        searchInput.focus();
        // 5. Помещаем курсор в конец
        searchInput.setSelectionRange(searchInput.value.length, searchInput.value.length);
    }
});