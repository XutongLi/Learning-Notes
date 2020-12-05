SELECT 
  CAST(premiered/10*10 AS TEXT) || 's' AS decade,
  COUNT(*) AS num_movies
  FROM titles
  WHERE premiered is not null
  GROUP BY decade
  ORDER BY num_movies DESC
  ;
