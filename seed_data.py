import psycopg2


def execute_sql_file(file_path, db_config):
    """
    Execute SQL commands from a file.

    :param file_path: Path to the SQL file.
    :param db_config: Database configuration dictionary with keys 'dbname', 'user', 'password', 'host', and 'port'.
    """
    # Read the SQL file
    with open(file_path, "r") as file:
        sql = file.read()

    # Connect to the PostgreSQL database
    conn = psycopg2.connect(
        dbname=db_config["dbname"],
        user=db_config["user"],
        password=db_config["password"],
        host=db_config["host"],
        port=db_config["port"],
    )

    try:
        # Create a new database session and return a new cursor
        with conn.cursor() as cursor:
            # Execute the SQL command
            cursor.execute(sql)
            # Commit the transaction
            conn.commit()
            print("SQL file executed successfully")

    except Exception as e:
        print(f"Error executing SQL file: {e}")
        # Rollback the transaction in case of error
        conn.rollback()

    finally:
        # Close the database connection
        conn.close()


def main():
    # Database configuration
    db_config = {
        "dbname": "beehome",
        "user": "hiro",
        "password": "1",
        "host": "localhost",
        "port": "5432",
    }

    # Path to your SQL file
    sql_file_path = ["CreateTables_vn_units.sql", "ImportData_vn_units.sql","Import_data_be.sql"]

    # Execute the SQL file
    for path in sql_file_path:
        execute_sql_file(path, db_config)


if __name__ == "__main__":
    main()
