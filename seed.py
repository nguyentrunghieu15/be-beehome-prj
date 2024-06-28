import psycopg2
import json
import os
import argparse


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

        sql = "DELETE FROM postal_codes;"
        cursor.execute(sql)

        # Prepare the SQL statement with parameter placeholder
        sql = """
            INSERT INTO postal_codes (country_code, zipcode, place, state, state_code, province, province_code, community, community_code, latitude, longitude)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """

        # Open the JSON file
        with open(data_file, 'r') as file:
            data = file.read()
            data_dict = json.loads(data)
            # Create a dictionary to store unique zip codes
            unique_zipcodes = {}
            
            # Filter out items with duplicate zip codes
            filtered_data = []
            for item in data_dict:
                zipcode = item['zipcode']
                if zipcode not in unique_zipcodes:
                    unique_zipcodes[zipcode] = True
                    filtered_data.append(item)
                    # Loop through each JSON object in the file
            for item in filtered_data:
                item['latitude']=0
                item['longitude']=0

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


def insert_service(db_name, db_user, db_password, db_host, db_port):
# Connect to your PostgreSQL database
    conn = psycopg2.connect(
        dbname=db_name,
        user=db_user,
        password=db_password,
        host=db_host,
        port=db_port
    )

    # Create a cursor object
    cur = conn.cursor()
    
    # Define SQL queries to delete existing data
    delete_queries = [
        "DELETE FROM public.services",
        "DELETE FROM public.group_services"
    ]

    # Execute delete queries
    for query in delete_queries:
        cur.execute(query)


    # Define the SQL queries to insert data
    queries = [
        # Inserting records into group_services table
        """
        INSERT INTO public.group_services (created_at, name, detail)
        VALUES
            (NOW(), 'Home Cleaning', 'Group of services related to cleaning homes')
        """,

        # Retrieving the ID of the inserted group service
        """
        SELECT id FROM public.group_services WHERE name = 'Home Cleaning'
        """,

        # Inserting records into services table for Home Cleaning
        """
        INSERT INTO public.services (created_at, name, detail, price, unit_price, group_service_id)
        VALUES
            (NOW(), 'Window Cleaning (Interior)', 'Cleaning interior windows', 50.00, 'USD', %(group_service_id)s),
            (NOW(), 'Fridge Cleaning', 'Cleaning the interior of the fridge', 30.00, 'USD', %(group_service_id)s),
            (NOW(), 'Oven Cleaning', 'Cleaning the oven', 40.00, 'USD', %(group_service_id)s)
        """,

        # Inserting more records into services table for Home Cleaning services
        """
        INSERT INTO public.services (created_at, name, detail, price, unit_price, group_service_id)
        VALUES
            (NOW(), 'Carpet Cleaning', 'Cleaning carpets in the house', 60.00, 'USD', %(group_service_id)s),
            (NOW(), 'Bathroom Cleaning', 'Thorough cleaning of bathrooms', 45.00, 'USD', %(group_service_id)s),
            (NOW(), 'Kitchen Cleaning', 'Cleaning all surfaces in the kitchen', 55.00, 'USD', %(group_service_id)s),
            (NOW(), 'Bedroom Cleaning', 'Cleaning bedrooms including dusting and vacuuming', 50.00, 'USD', %(group_service_id)s),
            (NOW(), 'Laundry Service', 'Washing and folding laundry', 35.00, 'USD', %(group_service_id)s)
        """,

        # Inserting records into group_services table for Gardening Services
        """
        INSERT INTO public.group_services (created_at, name, detail)
        VALUES
            (NOW(), 'Gardening Services', 'Group of services related to gardening and outdoor maintenance')
        """,

        # Retrieving the ID of the inserted group service for Gardening Services
        """
        SELECT id FROM public.group_services WHERE name = 'Gardening Services'
        """,

        # Inserting records into services table for Gardening Services
        """
        INSERT INTO public.services (created_at, name, detail, price, unit_price, group_service_id)
        VALUES
            (NOW(), 'Lawn Mowing', 'Mowing and trimming grass in the garden', 40.00, 'USD', %(group_service_id)s),
            (NOW(), 'Hedge Trimming', 'Trimming and shaping hedges and bushes', 50.00, 'USD', %(group_service_id)s),
            (NOW(), 'Weed Control', 'Removing weeds from garden beds and pathways', 30.00, 'USD', %(group_service_id)s),
            (NOW(), 'Planting Services', 'Planting new flowers, shrubs, or trees', 45.00, 'USD', %(group_service_id)s)
        """,

        # Inserting records into group_services table for Home Repair
        """
        INSERT INTO public.group_services (created_at, name, detail)
        VALUES
            (NOW(), 'Home Repair', 'Group of services related to home repair and maintenance')
        """,

        # Retrieving the ID of the inserted group service for Home Repair
        """
        SELECT id FROM public.group_services WHERE name = 'Home Repair'
        """,

        # Inserting records into services table for Home Repair
        """
        INSERT INTO public.services (created_at, name, detail, price, unit_price, group_service_id)
        VALUES
            (NOW(), 'Plumbing Repair', 'Fixing plumbing issues such as leaks or clogs', 60.00, 'USD', %(group_service_id)s),
            (NOW(), 'Electrical Repair', 'Repairing electrical wiring or outlets', 70.00, 'USD', %(group_service_id)s),
            (NOW(), 'Painting Services', 'Interior or exterior painting of walls or surfaces', 80.00, 'USD', %(group_service_id)s),
            (NOW(), 'Roof Repair', 'Fixing leaks or damages on the roof', 90.00, 'USD', %(group_service_id)s)
        """
    ]
    group_service_id = ''
    # Execute the SQL queries
    for query in queries:
        if 'SELECT' in query:
            cur.execute(query)
            group_service_id = cur.fetchone()[0]
        else:
            cur.execute(query, {'group_service_id': group_service_id})

    # Commit the transaction
    conn.commit()
    print("Complete insert service")

    # Close the cursor and connection
    cur.close()
    conn.close()




def main(seed_zipcode=False, seed_services=False):
  

    # Replace placeholders with your actual database credentials and table name
    db_name = "bee-home"
    db_user = "hiro"
    db_password = "123456a@"
    db_host = "localhost"
    db_port = 5433  # Default PostgreSQL port

    if seed_zipcode:
          # Example usage
        folder_path = "./postal-codes-json-xml-csv/data"  # Replace with your actual folder path
        json_files = get_json_files(folder_path)

        if json_files:
            print("Found JSON files:")
            print(json_files)
        else:
            print("No JSON files found in the specified folder.")
        for file in json_files:
            insert_json_data_from_file(file, db_name, db_user, db_password, db_host, db_port)

    if seed_services:
        insert_service(db_name, db_user, db_password, db_host, db_port)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Process JSON files and optionally seed zipcodes or services into database.")
    parser.add_argument("--seed-zipcode", action="store_true", help="Whether to seed zipcodes into the database.")
    parser.add_argument("--seed-services", action="store_true", help="Whether to seed services into the database.")
    args = parser.parse_args()
    
    main(args.seed_zipcode, args.seed_services)

