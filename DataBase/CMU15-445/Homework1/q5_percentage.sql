SELECT
  CAST(premiered/10*10 AS TEXT) || 's' AS decade,
  ROUND(CAST(COUNT(*) AS REAL) / (SELECT COUNT(*) FROM titles) * 100.0, 4) as percentage
  FROM titles
  WHERE premiered is not null
  GROUP BY decade
  ORDER BY percentage DESC, decade ASC
  ;
