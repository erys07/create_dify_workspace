package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	tenantID := os.Getenv("TENANT_ID")
	accountID := os.Getenv("ACCOUNT_ID")

	var workspaceName string
	fmt.Print("what's the name of your workspace?")
	fmt.Scanln(&workspaceName)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error when opening:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("error when connecting:", err)
	}
	fmt.Println("successfully connected!")

	query := fmt.Sprintf(`
    INSERT INTO tenants (
        id,
        name,
        encrypt_public_key,
        plan,
        status,
        created_at,
        updated_at,
        custom_config
    )
    SELECT
        gen_random_uuid(),
        '%s',
        encrypt_public_key,
        plan,
        status,
        NOW(),
        NOW(),
        custom_config
    FROM tenants
    WHERE id = '%s'
    RETURNING id;`, workspaceName, tenantID)

	var newTenantID string
	err = db.QueryRow(query).Scan(&newTenantID)
	if err != nil {
		log.Fatal("error to insert tenant:", err)
	}

	fmt.Println("New tenant inserted:")

	joinQuery := `
    INSERT INTO tenant_account_joins (
        id,
        tenant_id,
        account_id,
        role,
        invited_by,
        created_at,
        updated_at,
        current
    ) VALUES (
        gen_random_uuid(),
        $1,
        $2,
        'owner',
        NULL,
        NOW(),
        NOW(),
        TRUE
    );`

	_, err = db.Exec(joinQuery, newTenantID, accountID)
	if err != nil {
		log.Fatal("error when inserting join:", err)
	}

	fmt.Println("successfully inserted into tenant_account_joins table!")
}
