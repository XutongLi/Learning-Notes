WITH hamill_movies(title_id) AS (
  SELECT crew.title_id
    FROM crew
    JOIN people
    ON crew.person_id == people.person_id AND people.name == "Mark Hamill" AND people.born == 1951
)
SELECT titles.primary_title
  FROM crew
  JOIN people
  ON crew.person_id == people.person_id AND people.name == "George Lucas" AND people.born == 1944 AND crew.title_id IN hamill_movies
  JOIN titles
  ON crew.title_id == titles.title_id AND titles.type == "movie"
  ORDER BY titles.primary_title
;
