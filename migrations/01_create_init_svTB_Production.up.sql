-- СОЗДАТЬ ТАБЛИЦУ ПРОДУКЦИИ.
CREATE TABLE dbo.svTB_Production
(
    idProduction   INT IDENTITY (1,1)
        CONSTRAINT PK_skTB_Production
            PRIMARY KEY NONCLUSTERED,                         -- idProduction - ид продукции.
    PrName         VARCHAR(300)   DEFAULT ''        NOT NULL, -- PrName - наименование варианта упаковки продукции для упаковщика.
    PrShortName    VARCHAR(100)   DEFAULT ''        NOT NULL, -- PrShortName - короткое наименование продукции для этикетки.
    PrPackName     VARCHAR(300)   DEFAULT ''        NOT NULL, -- PrPackName - вариант упаковки.
    PrType         VARCHAR(100)   DEFAULT '',                 -- PrType - декларированная или нет.
    PrArticle      VARCHAR(5)     DEFAULT ''        NOT NULL  -- PrArticle - артикул варианта упаковки.
        CONSTRAINT IX_svTB_Production
            UNIQUE,
    PrColor        VARCHAR(20)    DEFAULT ''        NOT NULL, -- PrColor - цвет продукции.
    PrBarCode      VARCHAR(13)    DEFAULT '',                 -- PrBarCode - бар-код.
    PrCount        INT            DEFAULT 0         NOT NULL, -- PrCount - количество продукции в ряду.
    PrRows         INT            DEFAULT 0         NOT NULL, -- PrRows - количество рядов.
    PrWeight       DECIMAL(19, 3) DEFAULT 0         NOT NULL, -- PrWeight - вес п\п (кг).
    PrHWD          VARCHAR(100)   DEFAULT ''        NOT NULL, -- PrHWD - габариты (мм) 1000(высота)х1200(ширина)х1000(глубина).
    PrInfo         VARCHAR(1024),                             -- PrInfo - информация о продукции\комментарий.
    PrStatus       BIT            DEFAULT 1         NOT NULL, -- PrStatus - статус продукции.
    PrEditDate     DATETIME       DEFAULT GETDATE(),          -- PrEditDate - дата и время изменения записи.
    PrEditUser     INT            DEFAULT 1,                  -- PrEditUser - роль сотрудника. По умолчанию 1 - администратор, 5 - оператор.
    PrPart         INT            DEFAULT 0         NOT NULL, -- PrPart - номер текущей партии, номер партии и дата указываются вручную и не будут изменяться автоматически с течением времени.
    PrPartLastDate DATETIME       DEFAULT GETDATE() NOT NULL, -- PrPartLastDate - дата выпуска партии.
    PrPartAutoInc  SMALLINT       DEFAULT 1         NOT NULL, -- PrPartAutoInc - нумерация партии и даты! Ручная(0), Автоматическая(1), С указанной даты(2).
    PrPartRealDate DATETIME,                                  -- PrPartRealDate - дата продукции пока неизвестное поле.
    PrArchive      BIT            DEFAULT 0         NOT NULL, -- PrArchive - архивная запись или нет.
    PrPerGodn      SMALLINT       DEFAULT 0,                  -- PrPerGodn - срок годности в месяцах.
    PrSAP          VARCHAR(15),                               -- PrSAP - сап-код.
    PrProdType     BIT            DEFAULT 1         NOT NULL, -- PrProdType - тип продукции пищевая\не пищевая.
    PrUmbrella     BIT            DEFAULT 1         NOT NULL, -- PrUmbrella - беречь от влаги.
    PrSun          BIT            DEFAULT 1         NOT NULL, -- PrSun - беречь от солнца.
    PrDecl         BIT            DEFAULT 0         NOT NULL, -- PrDecl - декларирования или нет.
    PrParty        BIT            DEFAULT 0         NOT NULL, -- PrParty - партионная или нет.
    PrGL           SMALLINT       DEFAULT 0         NOT NULL, -- PrGL - петля Мёбиуса.
    PrVP           SMALLINT       DEFAULT 0         NOT NULL, -- PrVP - ванная печь.
    PrML           SMALLINT       DEFAULT 0         NOT NULL, -- PrML - машинная линия на печи.
    Created_at     DATETIME       DEFAULT GETDATE(),          -- Created_at - дата создания записи.
    Created_by     INT            DEFAULT 0         NOT NULL, -- Created_by - табельный номер сотрудника.
    Updated_at     DATETIME       DEFAULT GETDATE(),          -- Updated_at - дата изменения записи.
    Updated_by     INT            DEFAULT 0         NOT NULL, -- Updated_by - табельный номер сотрудника изменивший запись.

);
-- 1. Кластерный индекс (сейчас используется NONCLUSTERED на PK).
CREATE
CLUSTERED INDEX IX_svTB_Production_Created_at
    ON dbo.svTB_Production (Created_at DESC)
    WITH (DROP_EXISTING = OFF);

-- 2. Индекс для часто используемых фильтров.
CREATE
NONCLUSTERED INDEX IX_svTB_Production_Status_Archive
    ON dbo.svTB_Production (PrStatus, PrArchive)
    INCLUDE (PrArticle, PrName, PrShortName);

-- 3. Индекс для фильтрации по VP и ML (используется в процедуре генерации артикула).
CREATE
NONCLUSTERED INDEX IX_svTB_Production_VP_ML
    ON dbo.svTB_Production (PrVP, PrML)
    INCLUDE (PrArticle, PrStatus, PrArchive);


CREATE PROCEDURE dbo.svAFormsProductionAll -- ХП выводит список продукции.
    AS
BEGIN
    SET
NOCOUNT ON;

CREATE PROCEDURE dbo.svAFormsProductionAll @SortField NVARCHAR(50) = 'idProduction', -- ХП возвращает список продукции с сортировкой.
                                           @SortOrder NVARCHAR(4) = 'DESC'
AS
BEGIN
    SET NOCOUNT ON;

    -- Валидация параметров.
    SET @SortField = LTRIM(RTRIM(ISNULL(@SortField, 'idProduction')));
    SET @SortOrder = UPPER(LTRIM(RTRIM(ISNULL(@SortOrder, 'DESC'))));

    -- Безопасный список полей.
    IF @SortOrder NOT IN ('ASC', 'DESC')
        SET @SortOrder = 'DESC';

    DECLARE @SQL NVARCHAR(MAX);

    IF @SortField NOT IN ('idProduction', 'PrArticle', 'PrPackName', 'PrShortName',
                          'PrColor', 'PrCount', 'PrRows', 'PrEditDate')
        SET @SortField = 'idProduction';

    -- Безопасное формирование.
    DECLARE @OrderBy NVARCHAR(100) = QUOTENAME(@SortField) + ' ' + @SortOrder;

    SET @SQL = N'
        SELECT idProduction,
               PrShortName,
               PrPackName,
               PrArticle,
               PrColor,
               PrCount,
               PrRows,
               PrHWD,
               PrEditDate
        FROM dbo.svTB_Production
        ORDER BY ' + @OrderBy

    EXEC sp_executesql @SQL

END
GO;

CREATE PROCEDURE dbo.svAFormsProductionFilterById -- ХП ищет продукцию по артиклю и наименованию и коду продукции.
    @ArticlePattern VARCHAR(7) = N'', -- Паттерн поиска (например, "1", "12")
    @NamePattern NVARCHAR(100) = N'', -- Паттерн поиска по имени\названию продукции
    @IdPattern VARCHAR(7) = N'' -- Паттерн поиска по коду продукции.
AS
BEGIN

    SET NOCOUNT ON;

SELECT idProduction,
       PrShortName,
       PrPackName,
       PrArticle,
       PrColor,
       PrCount,
       PrRows,
       PrHWD,
       PrEditDate
FROM dbo.svTB_Production
WHERE PrArticle LIKE '%' + @ArticlePattern + '%'
  AND (@NamePattern = '' OR PrPackName LIKE '%' + @NamePattern + '%')
  AND idProduction LIKE '%' + @IdPattern + '%'
ORDER BY idProduction;

END
GO;
