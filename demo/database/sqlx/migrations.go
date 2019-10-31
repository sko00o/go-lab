package database

func up(driverName string) string {
	switch driverName {
	case "postgres":
		return `
CREATE TABLE IF NOT EXISTS peoples (
	id bigserial PRIMARY KEY,
	created_at TIMESTAMP with time zone NOT NULL, 
	updated_at TIMESTAMP with time zone NOT NULL,
	deleted_at TIMESTAMP with time zone NULL,
	name VARCHAR(255) NULL, 
	age INT NULL
);
`
	case "mysql":
		return `
CREATE TABLE IF NOT EXISTS peoples (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP(6) NOT NULL, 
	updated_at TIMESTAMP(6) NOT NULL,
	deleted_at TIMESTAMP(6),
	name VARCHAR(255) NULL, 
	age INT NULL
);
`
	default:
		return ""
	}
}

func down(_ string) string {
	return `
DROP TABLE IF EXISTS peoples;
`
}
