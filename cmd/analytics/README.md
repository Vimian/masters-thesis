# sql queries

```
Count of files being or done being analyzed

SELECT
    COUNT(DISTINCT file_name) as files
FROM "analytics";
```

```
Different window lengths and dictionaries for compression size lesser than original size

SELECT
    COUNT(*) as succesfull_compressions
FROM "analytics"
WHERE "dictionary_limit_reached" = 0;
```

```
Total count of different window lengths and files

SELECT
    COUNT(*) as total_rows
FROM "analytics";
```
