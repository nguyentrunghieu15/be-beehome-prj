import psycopg2 
from pymongo import MongoClient

dbname="beehome"
user="hiro"
password="1"
host="localhost"
port="5432"
uri = "mongodb://hiro:12345678@localhost:27017/"
dbmongo ="lms"

# Kết nối tới PostgreSQL
pg_conn = psycopg2.connect(
    dbname=dbname,
    user=user,
    password=password,
    host=host,
    port=port
)
pg_cursor = pg_conn.cursor()

# Kết nối tới MongoDB
mongo_client = MongoClient(uri)
mongo_db = mongo_client[dbmongo]

# Chuyển đổi dữ liệu User
pg_cursor.execute("SELECT id, provider_id, role FROM users")
users = pg_cursor.fetchall()
mongo_users = [{"user_id": user[0], "provider_id": user[1], "role": user[2]} for user in users]
mongo_db.users.insert_many(mongo_users)

# Chuyển đổi dữ liệu Provider
pg_cursor.execute("SELECT id, user_id FROM providers")
providers = pg_cursor.fetchall()
mongo_providers = [{"provider_id": provider[0], "user_id": provider[1]} for provider in providers]
mongo_db.providers.insert_many(mongo_providers)

# Chuyển đổi dữ liệu Hire
pg_cursor.execute("SELECT id, provider_id, user_id FROM hires")
hires = pg_cursor.fetchall()
mongo_hires = [{"hire_id": hire[0], "provider_id": hire[1], "user_id": hire[2]} for hire in hires]
mongo_db.hires.insert_many(mongo_hires)

# Chuyển đổi dữ liệu Review
pg_cursor.execute("SELECT id, hire_id, provider_id, user_id FROM reviews")
reviews = pg_cursor.fetchall()
mongo_reviews = [{"review_id": review[0], "hire_id": review[1], "provider_id": review[2], "user_id": review[3]} for review in reviews]
mongo_db.reviews.insert_many(mongo_reviews)

# Chuyển đổi dữ liệu GroupService
pg_cursor.execute("SELECT id FROM group_services")
group_services = pg_cursor.fetchall()
mongo_group_services = [{"group_service_id": group_service[0]} for group_service in group_services]
mongo_db.group_services.insert_many(mongo_group_services)

# Chuyển đổi dữ liệu Service
pg_cursor.execute("SELECT id, group_service_id FROM services")
services = pg_cursor.fetchall()
mongo_services = [{"service_id": service[0], "group_service_id": service[1]} for service in services]
mongo_db.services.insert_many(mongo_services)

# Chuyển đổi dữ liệu SocialMedia
pg_cursor.execute("SELECT id, provider_id FROM social_media")
social_medias = pg_cursor.fetchall()
mongo_social_medias = [{"social_media_id": social_media[0], "provider_id": social_media[1]} for social_media in social_medias]
mongo_db.social_media.insert_many(mongo_social_medias)

# Chuyển đổi dữ liệu PaymentMethod
pg_cursor.execute("SELECT id, provider_id FROM payment_methods")
payment_methods = pg_cursor.fetchall()
mongo_payment_methods = [{"payment_method_id": payment_method[0], "provider_id": payment_method[1]} for payment_method in payment_methods]
mongo_db.payment_methods.insert_many(mongo_payment_methods)

# Đóng kết nối
pg_cursor.close()
pg_conn.close()
mongo_client.close()
