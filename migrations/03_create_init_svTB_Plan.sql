-- СОЗДАТЬ ТАБЛИЦУ ПЛАН ПРОДУКЦИИ
CREATE TABLE dbo.svTB_Plan
(
    idPlan        INT IDENTITY (1, 1)
        CONSTRAINT PK_svTb_Plan PRIMARY KEY NONCLUSTERED,
    PlanShift     INT      DEFAULT (-1) NOT NULL, -- PlanShift - номер смены
    extProduction INT      DEFAULT 0    NOT NULL, -- extProduction - внешний ключ с таблицей svTB_Production.idProduction
    extSector     INT                   NOT NULL, -- extSector - внешний ключ с таблицей svTB_Sector.idSector
    PlanCount     INT      DEFAULT 0    NOT NULL, -- PlanCount - плановое кол-во продукции
    PlanDate      DATETIME DEFAULT GETDATE(),     -- PlanDate - дата плана создания
    PlanInfo      VARCHAR(1024),                  -- PlanInfo - комментарий к плану
    PlEditDate    DATETIME DEFAULT GETDATE(),     -- PlEditDate - дата редактирования плана
    Created_at    DATETIME DEFAULT GETDATE(),     -- Created_at - дата создания записи
    Created_by    INT      DEFAULT 0    NOT NULL, -- Created_by - табельный номер сотрудника
    Updated_at    DATETIME DEFAULT GETDATE(),     -- Updated_at - дата изменения записи
    Updated_by    INT      DEFAULT 0    NOT NULL, -- Updated_by - табельный номер сотрудника изменивший запись
);

CREATE INDEX IX_svTB_Plan_extProduction ON dbo.svTB_Plan(extProduction);
CREATE INDEX IX_svTB_Plan_extSector ON dbo.svTB_Plan(extSector);
CREATE INDEX IX_svTB_Plan_PlanDate ON dbo.svTB_Plan(PlanDate);

CREATE PROCEDURE dbo.svAFormsPlansWithSortingAndFiltering -- ХП возвращает список планов с сортировкой и фильтрацией
    @SortField NVARCHAR(50) = 'idPlan', -- поле для сортировки
    @SortOrder NVARCHAR(4) = 'DESC', -- порядок сортировки
    @StartDate DATE = NULL, -- дата начала периода
    @EndDate DATE = NULL -- дата конца периода
AS
BEGIN
    SET NOCOUNT ON;

    -- Валидация параметров
    SET @SortField = LTRIM(RTRIM(ISNULL(@SortField, 'idPlan')));
    SET @SortOrder = UPPER(LTRIM(RTRIM(ISNULL(@SortOrder, 'DESC'))));

    -- Безопасный список полей
    IF @SortOrder NOT IN ('ASC', 'DESC')
        SET @SortOrder = 'DESC';

    -- Полный список полей для сортировки из всех таблиц
    IF @SortField NOT IN (
        -- Поля из svTB_Plan
                          'idPlan', 'PlanShift', 'extProduction', 'extSector',
                          'PlanCount', 'PlanDate', 'PlEditDate',
        -- Поля из svTB_Production
                          'PrShortName', 'PrArticle', 'PrPackName', 'PrColor',
                          'PrCount', 'PrRows', 'PrHWD', 'PrType', 'PrWeight',
        -- Поля из svTB_Sector
                          'SectorName'
        )
        SET @SortField = 'idPlan';

    -- Формируем ORDER BY в зависимости от выбранного поля
    DECLARE @OrderBy NVARCHAR(100) = QUOTENAME(@SortField) + ' ' + @SortOrder;

    DECLARE @SQL NVARCHAR(MAX);

    SET @SQL = N'
        SELECT
            p.idPlan,
            p.PlanShift,

            -- Информация о продукции
            p.extProduction,
            ISNULL(pr.PrName, '''') AS KDName,
            ISNULL(pr.PrShortName, '''') AS KDNameShort,
            ISNULL(pr.PrType, '''') AS KDNameType,
            ISNULL(pr.PrArticle, '''') AS KDNameArticle,
            ISNULL(pr.PrColor, '''') AS KDNameColor,
            ISNULL(pr.PrCount, 0) AS KDNameCount,
            ISNULL(pr.PrRows, 0) AS KDNameRows,
            ISNULL(pr.PrHWD, '''') AS KDNameHWD,
            ISNULL(pr.PrWeight, 0) AS KDNameWeight,

            -- Информация о секторе
            p.extSector,
            s.SectorName,

            -- Основная информация плана
            p.PlanCount,
            p.PlanDate,
            ISNULL(p.PlanInfo, '''') AS PlanInfo,
            p.PlEditDate

        FROM dbo.svTB_Plan p
        LEFT JOIN dbo.svTB_Production pr ON p.extProduction = pr.idProduction
        LEFT JOIN dbo.svTB_Sector s ON p.extSector = s.idSector
        WHERE p.PlanDate >= @StartDateParam AND p.PlanDate <= @EndDateParam
        ORDER BY ' + @OrderBy;

    -- Используем другие имена параметров в динамическом SQL, чтобы избежать конфликта
EXEC sp_executesql @SQL,
         N'@StartDateParam DATE, @EndDateParam DATE',
         @StartDateParam = @StartDate,
         @EndDateParam = @EndDate;
END
go

