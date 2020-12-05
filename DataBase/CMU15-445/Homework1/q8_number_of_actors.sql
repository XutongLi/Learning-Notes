WITH hamill_titles AS (
  SELECT DISTINCT(crew.title_id)
    FROM people
    JOIN crew
    ON crew.person_id == people.person_id AND people.name == "Mark Hamill" AND people.born == 1951
)
SELECT COUNT(DISTINCT(crew.person_id))
  FROM crew
  WHERE (crew.category == "actor" OR crew.category == "actress") AND crew.title_id in hamill_titles
;
