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
    PrArticle      VARCHAR(10)    DEFAULT ''        NOT NULL  -- PrArticle - артикул варианта упаковки.
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
    PrCreationDate DATETIME       DEFAULT GETDATE(),          -- PrCreationDate - дата создания продукции.
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
    PrPerfumery    BIT            DEFAULT 0         NOT NULL, -- PrPerfumery - тип парфимерия.
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

CREATE PROCEDURE dbo.svAFormsProductionAdd -- ХП добавляет продукцию.
-- Продукция
    @PrName VARCHAR(300),
    @PrShortName VARCHAR(100),
    @PrPackName VARCHAR(300),
    @PrDecl BIT,
    @PrSun BIT,
    @PrProdType BIT,
    @PrParty BIT,
    @PrUmbrella BIT,
    @PrPerfumery BIT,
    @PrColor VARCHAR(20),
    @PrGL SMALLINT,
    @PrArticle VARCHAR(2), -- вводим только 2 цифры ВП и МЛ
    @PrSAP VARCHAR(15),
    -- Упаковка
    @PrCount INT,
    @PrRows INT,
    @PrWeight DECIMAL(19, 3),
    @PrHWD VARCHAR(100),
    -- Комментарии
    @PrInfo VARCHAR(1024),
    @PrPart INT,
    @PrPartLastDate DATETIME,
    @PrPartAutoInc SMALLINT,
    @PrPerGond SMALLINT,
    -- Аудит
    @CreatedBy INT,
    @UpdatedBy INT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @NewPrArticle VARCHAR(5);
    DECLARE @MaxSequence INT;
    DECLARE @NewSequence VARCHAR(3);
    DECLARE @PrVP SMALLINT;
    DECLARE @PrML SMALLINT;
    DECLARE @PrType VARCHAR(100);

    -- 1. Проверяем, что артикул состоит из 2 цифр
    IF LEN(@PrArticle) <> 2 OR ISNUMERIC(@PrArticle) = 0
BEGIN
            RAISERROR (N'Артикул должен состоять из 2 цифр (1-ВП, 2-МЛ "Например:12")', 16, 1);
END;

    -- 2. Извлекаем VP и ML из первых 2 цифр
    SET @PrVP = CAST(SUBSTRING(@PrArticle, 1, 1) AS SMALLINT);
    SET @PrML = CAST(SUBSTRING(@PrArticle, 2, 1) AS SMALLINT);

    SET @PrType = IIF(@PrDecl = 1, N'Декларированная', N'');

    -- 3. Находим максимальные последние 3 цифры для этого префикса
SELECT @MaxSequence = ISNULL(MAX(
                                     CAST(SUBSTRING(PrArticle, 3, 3) AS INT)
                             ), -1)
FROM dbo.svTB_Production
WHERE PrArticle LIKE @PrArticle + '%' -- Ищем по первым двум цифрам
  AND ISNUMERIC(PrArticle) = 1
  AND LEN(PrArticle) = 5;

-- 4. Увеличиваем на 1 (если нет записей, -1 + 1 = 0)
SET @MaxSequence = @MaxSequence + 1;

    -- 5. Проверяем, не превышает ли 999
    IF @MaxSequence > 999
BEGIN
            RAISERROR (N'Достигнут максимальный номер последовательности (999) для префикса %s', 16, 1, @PrArticle);
END;

    -- 6. Форматируем в 3 цифры с ведущими нулями
    SET @NewSequence = RIGHT('000' + CAST(@MaxSequence AS VARCHAR(3)), 3);

    -- 7. Формируем полный артикул из 5 цифр
    SET @NewPrArticle = @PrArticle + @NewSequence;

INSERT INTO dbo.svTB_Production
(PrName,
 PrShortName,
 PrPackName,
 PrType,
 PrDecl,
 PrSun,
 PrProdType,
 PrParty,
 PrUmbrella,
 PrPerfumery,
 PrColor,
 PrGL,
 PrArticle,
 PrSAP,
 PrCount,
 PrRows,
 PrWeight,
 PrHWD,
 PrInfo,
 PrPart,
 PrPartLastDate,
 PrPartAutoInc,
 PrPerGodn,
 PrVP,
 PrML,
 Created_by,
 Updated_by)
VALUES (@PrName,
        @PrShortName,
        @PrPackName,
        @PrType,
        @PrDecl,
        @PrSun,
        @PrProdType,
        @PrParty,
        @PrUmbrella,
        @PrPerfumery,
        @PrColor,
        @PrGL,
        @NewPrArticle,
        @PrSAP,
        @PrCount,
        @PrRows,
        @PrWeight,
        @PrHWD,
        @PrInfo,
        @PrPart,
        @PrPartLastDate,
        @PrPartAutoInc,
        @PrPerGond,
        @PrVP,
        @PrML,
        @CreatedBy,
        @UpdatedBy);

-- 9. Возвращаем сгенерированный артикул
--     SELECT @NewPrArticle AS GeneratedArticle;

-- 10. Вставляем нумератор, если его не существует
IF (NOT EXISTS(SELECT NumName FROM dbo.svTB_Numerator WHERE NumName = @NewPrArticle))
BEGIN
INSERT INTO dbo.svTB_Numerator(NumName, Created_by, Updated_by)
VALUES (@NewPrArticle, @CreatedBy, @UpdatedBy);
END

END
GO;



-- exec dbo.svAFormsProductionAdd N'TEST', N'TEST', N'TEST', 1, 0, 1, 0, 1, N'RED', 500, N'12', '', 10, 2, 523,
--      N'100x100x111', '', 100, '20251223 00:00:00.000', 1, 50, 1, 1;



CREATE PROCEDURE dbo.svAFormsProductionAll -- ХП выводит список продукции.
    AS
BEGIN
    SET
NOCOUNT ON;

CREATE PROCEDURE dbo.svAFormsProductionAll @SortField NVARCHAR(50) = 'idProduction', -- ХП возвращает список продукции с сортировкой.
                                           @SortOrder NVARCHAR(4) = 'DESC'
AS
BEGIN
    SET
NOCOUNT ON;

    -- Валидация параметров.
    SET
@SortField = LTRIM(RTRIM(ISNULL(@SortField, 'idProduction')));
    SET
@SortOrder = UPPER(LTRIM(RTRIM(ISNULL(@SortOrder, 'DESC'))));

    -- Безопасный список полей.
    IF
@SortOrder NOT IN ('ASC', 'DESC')
        SET @SortOrder = 'DESC';

    DECLARE
@SQL NVARCHAR(MAX);

    IF
@SortField NOT IN ('idProduction', 'PrArticle', 'PrPackName', 'PrShortName',
                          'PrColor', 'PrCount', 'PrRows', 'PrEditDate')
        SET @SortField = 'idProduction';

    -- Безопасное формирование.
    DECLARE
@OrderBy NVARCHAR(100) = QUOTENAME(@SortField) + ' ' + @SortOrder;

    SET
@SQL = N'
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

    SET
NOCOUNT ON;

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

CREATE PROCEDURE dbo.svAFormsProductionUpd -- ХП обновляет продукцию.
    @IdProduction INT,
-- Продукция
    @PrName VARCHAR(300),
    @PrShortName VARCHAR(100),
    @PrPackName VARCHAR(300),
    @PrDecl BIT,
    @PrSun BIT,
    @PrProdType BIT,
    @PrParty BIT,
    @PrUmbrella BIT,
    @PrPerfumery BIT,
    @PrColor VARCHAR(20),
    @PrGL SMALLINT,
    @PrArticle VARCHAR(5),
    @PrSAP VARCHAR(15),
    -- Упаковка
    @PrCount INT,
    @PrRows INT,
    @PrWeight DECIMAL(19, 3),
    @PrHWD VARCHAR(100),
    -- Комментарии
    @PrInfo VARCHAR(1024),
    @PrPart INT,
    @PrPartLastDate DATETIME,
    @PrPartAutoInc SMALLINT,
    @PrPerGond SMALLINT,

    -- Аудит
    @CreatedBy INT,
    @UpdatedBy INT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @PrType VARCHAR(100);
    SET @PrType = IIF(@PrDecl = 1, N'Декларированная', N'');

UPDATE svTB_Production
SET PrName         = @PrName,
    PrShortName    = @PrShortName,
    PrPackName     = @PrPackName,
    PrType         = @PrType,
    PrDecl         = @PrDecl,
    PrSun          = @PrSun,
    PrProdType     = @PrProdType,
    PrParty        = @PrParty,
    PrUmbrella     = @PrUmbrella,
    PrPerfumery    = @PrPerfumery,
    PrColor        = @PrColor,
    PrGL           = @PrGL,
    PrArticle      = @PrArticle,
    PrSAP          = @PrSAP,
    PrCount        = @PrCount,
    PrRows         = @PrRows,
    PrWeight       = @PrWeight,
    PrHWD          = @PrHWD,
    PrInfo         = @PrInfo,
    PrEditDate     = GETDATE(),
    PrPart         = @PrPart,
    PrPartLastDate = @PrPartLastDate,
    PrPartAutoInc  = @PrPartAutoInc,
    PrPerGodn      = @PrPerGond,
    Created_by     = @CreatedBy,
    Updated_by     = @UpdatedBy,
    Updated_at     = GETDATE()
WHERE idProduction = @IdProduction;

END
GO;

CREATE PROCEDURE dbo.svAFormsProductionFindById -- ХП - ищет запись по ИД.
    @IdProduction INT
AS
BEGIN
    SET NOCOUNT ON;

SELECT idProduction,
       PrName,
       PrShortName,
       PrPackName,
       PrType,
       PrArticle,
       PrColor,
       PrBarCode,
       PrCount,
       PrRows,
       PrWeight,
       PrHWD,
       PrInfo,
       PrStatus,
       PrEditDate,
       PrEditUser,
       PrPart,
       PrPartLastDate,
       PrPartAutoInc,
       COALESCE(PrPartRealDate, '') as PrPartRealDate, -- преобразуем NULL в пустую строку
       PrArchive,
       PrPerGodn,
       PrSAP,
       PrProdType,
       PrUmbrella,
       PrSun,
       PrDecl,
       PrParty,
       PrGL,
       PrVP,
       PrML,
       Created_at,
       Created_by,
       Updated_at,
       Updated_by,
       COALESCE(PrCreationDate, '') as PrCreationDate,
       PrPerfumery
FROM dbo.svTB_Production
WHERE idProduction = @IdProduction;

END
GO;
