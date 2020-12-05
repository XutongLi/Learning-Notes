SELECT type, count(*) AS title_count FROM titles GROUP BY type ORDER BY title_count ASC;
