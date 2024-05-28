import psycopg2
import json

def insert_json_data_from_file(data_file, db_name, db_user, db_password, db_host, db_port):
    """Inserts JSON data from a file into a PostgreSQL database.

    Args:
        data_file (str): The path to the JSON file containing the data.
        db_name (str): The name of the PostgreSQL database.
        db_user (str): The username to connect to the database.
        db_password (str): The password to connect to the database.
        db_host (str): The hostname or IP address of the database server.
        db_port (int): The port number of the database server.

    Returns:
        None
    """

    try:
        # Connect to the PostgreSQL database
        connection = psycopg2.connect(
            dbname=db_name,
            user=db_user,
            password=db_password,
            host=db_host,
            port=db_port
        )

        # Create a cursor object
        cursor = connection.cursor()

        # Prepare the SQL statement with parameter placeholder
        sql = """
            INSERT INTO your_table_name (country_code, zipcode, place, state, state_code, province, province_code, community, community_code, latitude, longitude)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """

        # Open the JSON file
        with open(data_file, 'r') as file:
            # Loop through each JSON object in the file
            for line in file:
                data_dict = json.loads(line.strip())

                # Execute the insertion with parameter values
                cursor.execute(sql, (data_dict['country_code'], data_dict['zipcode'], data_dict['place'],
                                     data_dict['state'], data_dict['state_code'], data_dict['province'],
                                     data_dict['province_code'], data_dict['community'], data_dict['community_code'],
                                     data_dict['latitude'], data_dict['longitude']))

        # Commit the changes to the database
        connection.commit()

        print("Data from file inserted successfully!")

    except (Exception, psycopg2.Error) as error:
        print("Error while inserting data:", error)

    finally:
        # Close the cursor and connection objects, even in case of errors
        if cursor:
            cursor.close()
        if connection:
            connection.close()

# Replace placeholders with your actual database credentials and table name
data_file = "zip_code.json"  # Path to your JSON file
db_name = "your_database_name"
db_user = "your_username"
db_password = "your_password"
db_host = "your_database_host"
db_port = 5432  # Default PostgreSQL port

insert_json_data_from_file(data_file, db_name, db_user, db_password, db_host, db_port)