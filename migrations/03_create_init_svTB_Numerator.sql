CREATE TABLE dbo.svTB_Numerator -- Таблица для хранения артикула и генератора чисел для этикетки бар-кода.
(
    idNum      INT IDENTITY (1,1)                  -- idNum - код.
        CONSTRAINT PK_skTB_Numerator
            PRIMARY KEY NONCLUSTERED,
    NumName    VARCHAR(5)                NOT NULL  -- NumName - нумератор/артикул.
        CONSTRAINT IX_svTB_Numerator
            UNIQUE,
    NumBegin   INT      default 1        NOT NULL, -- NumBegin - начальный номер.
    NumEnd     INT      default 99999999 NOT NULL, -- NumEnd - конечный номер.
    NumStep    INT      default 1        NOT NULL, -- NumStep - шаг итерации.
    NumValue   INT      default 0        NOT NULL, -- NumValue - последнее использованное значение.
    Created_at DATETIME DEFAULT GETDATE(),         -- Created_at - дата создания записи.
    Created_by INT      DEFAULT 0        NOT NULL, -- Created_by - табельный номер сотрудника.
    Updated_at DATETIME DEFAULT GETDATE(),         -- Updated_at - дата изменения записи.
    Updated_by INT      DEFAULT 0        NOT NULL, -- Updated_by - табельный номер сотрудника изменивший запись.
);

