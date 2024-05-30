import psycopg2
import json
import os


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
            INSERT INTO postal_codes (country_code, zipcode, place, state, state_code, province, province_code, community, community_code, latitude, longitude)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """

        # Open the JSON file
        with open(data_file, 'r') as file:
            data = file.read()
            data_dict = json.loads(data)
            # Loop through each JSON object in the file
            for item in data_dict:
                print(item)

                # Execute the insertion with parameter values
                cursor.execute(sql, (item['country_code'], item['zipcode'], item['place'],
                                     item['state'], item['state_code'], item['province'],
                                     item['province_code'], item['community'], item['community_code'],
                                     item['latitude'], item['longitude']))

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


def get_json_files(folder_path):
  """
  This function retrieves all filenames with the .json extension within a directory.

  Args:
      folder_path: The path to the directory containing the JSON files.

  Returns:
      A list of filenames (including extension) that are JSON files.
  """
  json_files = []
  for filename in os.listdir(folder_path):
    if filename.endswith(".json"):
      json_files.append(f'{folder_path}/{filename}')
  return json_files

def main():
    # Example usage
    folder_path = "./postal-codes-json-xml-csv/data"  # Replace with your actual folder path
    json_files = get_json_files(folder_path)

    if json_files:
      print("Found JSON files:")
    if len(json_files):
        print(json_files)
    else:
        print("No JSON files found in the specified folder.")
#    Replace placeholders with your actual database credentials and table name
    db_name = "beehome"
    db_user = "hiro"
    db_password = "1"
    db_host = "localhost"
    db_port = 5432  # Default PostgreSQL port

    for file in json_files:
       insert_json_data_from_file(file, db_name, db_user, db_password, db_host, db_port)

main()