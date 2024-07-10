import pandas as pd


# Chuyển đổi dữ liệu thành DataFrame
df = pd.read_csv("/home/hiro/Downloads/group_services.csv")

# Tạo bảng tên
table_name = "services"  # Thay đổi tên bảng nếu cần thiết

# Lấy danh sách tên cột
column_names = df.columns.tolist()

# Tạo chuỗi INSERT query
insert_query = f"""
INSERT INTO {table_name} ({','.join(column_names)})
VALUES
"""

# Thêm giá trị cho mỗi bản ghi
for index, row in df.iterrows():
    values = [f"'{x}'" if not pd.isna(x) else "NULL" for x in row.tolist()]
    insert_query += f"({','.join(values)}),\n"

# Loại bỏ dấu phẩy cuối cùng và thêm dấu ngoặc đơn
insert_query = insert_query[:-2] + ";"

# In INSERT query
with open("text.txt","w") as fs:
    fs.write(insert_query)
