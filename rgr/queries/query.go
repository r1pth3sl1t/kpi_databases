package queries

import "strconv"

func PrepareInsertQuery(table string, data map[string]string) (string, []any) {
	query := "INSERT INTO public." + table + "("
	var values []any
	for column := range data {
		query += column + ","
		values = append(values, data[column])
	}
	query = query[:len(query)-1]
	query += ") VALUES ("
	for i := 1; i <= len(data); i++ {
		query += "$"
		query += strconv.Itoa(i)
		query += ","
	}
	query = query[:len(query)-1]
	query += ")"
	return query, values
}

func PrepareUpdateQuery(table string, data map[string]string, pkey map[string]string) (string, []any) {
	query := "UPDATE public." + table
	query += " SET "
	i := 1
	var values []any
	for column := range data {
		query += column + " = $" + strconv.Itoa(i) + ","
		values = append(values, data[column])
		i++
	}
	query = query[:len(query)-1]

	query += " WHERE "
	for column := range pkey {
		query += column + " = $" + strconv.Itoa(i) + " AND "
		values = append(values, pkey[column])
		i++
	}
	query = query[:len(query)-4]
	return query, values
}

func PrepareDeleteQuery(table string, pkey map[string]string) (string, []any) {
	var values []any
	query := "DELETE FROM " + table + " WHERE "
	i := 1
	for key := range pkey {
		query += key + " = $" + strconv.Itoa(i) + " AND "
		values = append(values, pkey[key])
		i++
	}

	query = query[:len(query)-4]

	return query, values
}

func GetUserGeneratingQuery() string {
	return `
	insert into public.user(user_id, first_name, last_name, email)
	select * from (select row_number() over() as user_id,
			   first_name,
			   last_name, 
			   lower(first_name || '_' || last_name) || '@mail.com' as email
from unnest(array['Pavlo', 'Ruslan', 'Anna',
				  'Ivan', 'Yulia', 'Anastasia',
				  'Volodymyr', 'Sofia', 'Anton',
				  'Olexiy', 'Olexandr', 'Ivanna',
				  'Bohdan', 'Oleh', 'Fedir',
				  'Hryhorii', 'Serhii', 'Dmytro',
				  'Polina', 'Kateryna', 'Svitlana',
				  'Stepan', 'Ostap', 'Halyna',
				  'Viktor', 'Kyrylo', 'Roman']) as first_name
cross join (select generate_series(1, 10000)) as user_id
cross join unnest(array['Ivanenko', 'Petrenko', 'Shevchenko',
						'Kuznets', 'Sobol', 'Mazur',
						'Kvitka', 'Sydorenko', 'Koval',
					   'Zhuk', 'Tkach', 'Tkachuk',
						'Vlasenko', 'Tymoshenko', 'Ostapenko',
					   'Rudenko','Moroz','Petrenko',
						'Pavlenko','Vasilenko', 'Kryvenko',
					   'Shpak', 'Los', 'Shvets',
						'Bondar', 'Savchenko', 'Korol']) as last_name) as users
where user_id not in (select user_id from public.user)
order by random()
limit $1;

`
}

func GenerateSkillsQuery() string {
	return `
	insert into users_to_skills(user_id, skill_id)
	select user_id, skill_id from "user"
	cross join skill
	where not(user_id in (select user_id from "users_to_skills") 
				  and 
			  skill_id in (select skill_id from "users_to_skills"))
	order by random()
	limit $1;
`
}

func GenerateConnectionQuery() string {
	return `
	insert into "connection"(u1, u2)
	select u1t.user_id as u1, u2t.user_id as u2 from "user" as u1t
	cross join "user" as u2t
	where u1t.user_id != u2t.user_id
	and not(
	u1t.user_id in (select u1 from connection) and
	u2t.user_id in (select u2 from connection)
	)
	order by random()
	limit $1;
`
}

func GetFetchingPrimaryKeysQuery() string {
	return `
	select distinct constraint_column_usage.table_name, constraint_column_usage.column_name, constraint_type
	from information_schema.constraint_column_usage 
	inner join information_schema.table_constraints 
	on constraint_column_usage.table_name = information_schema.table_constraints.table_name 
	where constraint_column_usage.table_schema = 'public' and constraint_type = 'PRIMARY KEY';
`
}

func GetFetchingTablesDataQuery() string {
	return `
	select table_name, column_name 
	from information_schema.columns 
	where table_schema = 'public' 
	order by table_name;
	`
}
