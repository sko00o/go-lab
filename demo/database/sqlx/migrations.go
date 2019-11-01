package database

func up(driverName string) []string {
	switch driverName {
	case "postgres":
		return []string{`
CREATE TABLE IF NOT EXISTS peoples (
	id bigserial PRIMARY KEY,
	created_at TIMESTAMP with time zone NOT NULL, 
	updated_at TIMESTAMP with time zone NOT NULL,
	deleted_at TIMESTAMP with time zone NULL,
	name VARCHAR(255) NULL, 
	age INT NULL
);`, `
CREATE INDEX idx_peoples_age ON peoples(age);
`,
		}
	case "mysql":
		return []string{`
CREATE TABLE IF NOT EXISTS peoples (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	created_at TIMESTAMP(6) NOT NULL, 
	updated_at TIMESTAMP(6) NOT NULL,
	deleted_at TIMESTAMP(6),
	name VARCHAR(255) NULL, 
	age INT NULL
);
`, `
CREATE INDEX idx_peoples_age ON peoples(age);
`,
		}
	default:
		return nil
	}
}

func down(driverName string) []string {
	switch driverName {
	case "postgres":
		return []string{`
DROP INDEX idx_peoples_age;
`, `
DROP TABLE IF EXISTS peoples;
`}
	case "mysql":
		return []string{`
DROP INDEX idx_peoples_age ON peoples;
`, `
DROP TABLE IF EXISTS peoples;
`,
		}
	default:
		return nil
	}
}
