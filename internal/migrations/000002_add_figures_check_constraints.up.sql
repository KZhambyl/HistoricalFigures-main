ALTER TABLE figures ADD CONSTRAINT figures_years_of_life_check CHECK (length(years_of_life)=9);
ALTER TABLE figures ADD CONSTRAINT figures_name_check CHECK (length(name) BETWEEN 8 AND 20);