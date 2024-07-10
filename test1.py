import pandas as pd

def replace_provider_id(users_df, providers_df) -> pd.DataFrame:
  """
  Hàm thay thế trường provider* của dữ liệu user bằng id của dữ liệu provider sao cho trường email tương ứng với trường user_id

  Args:
      users_df (pandas.DataFrame): DataFrame chứa dữ liệu người dùng.
      providers_df (pandas.DataFrame): DataFrame chứa dữ liệu nhà cung cấp.

  Returns:
      pandas.DataFrame: DataFrame sau khi thay thế trường provider*.
  """

  # Tạo dictionary lưu trữ mapping giữa email và provider_id
  email_to_provider_id = providers_df.set_index('user_id')['id'].to_dict()

  # Thay thế giá trị provider* bằng provider_id tương ứng
  for index, row in users_df.iterrows():
    provider_field = row['provider_id']  # Lấy tên trường provider*
    email = row['email']  # Lấy email
    email = email.removesuffix('@example.com')
    provider_id = email_to_provider_id.get(email)  # Tìm provider_id tương ứng

    if provider_id:
      users_df.iloc[index, 11] = provider_id
    else:
      users_df.iloc[index, 11] =None

  return users_df

def replace_user_id(users_df, providers_df) -> pd.DataFrame:
  """
  Hàm thay thế trường provider* của dữ liệu user bằng id của dữ liệu provider sao cho trường email tương ứng với trường user_id

  Args:
      users_df (pandas.DataFrame): DataFrame chứa dữ liệu người dùng.
      providers_df (pandas.DataFrame): DataFrame chứa dữ liệu nhà cung cấp.

  Returns:
      pandas.DataFrame: DataFrame sau khi thay thế trường provider*.
  """

  # Tạo dictionary lưu trữ mapping giữa email và provider_id
  filtered_df = users_df.dropna(subset=['provider_id']) 
  user_id_to_provider_id = filtered_df.set_index('provider_id')['id'].to_dict()

  # Thay thế giá trị provider* bằng provider_id tương ứng
  for index, row in providers_df.iterrows():
    id = row['id']  # Lấy email
    user_id = user_id_to_provider_id.get(id)  # Tìm provider_id tương ứng

    if user_id:
      providers_df.iloc[index, 11] = user_id
    else:
      providers_df.iloc[index, 11] =None

  return providers_df

# Đọc dữ liệu vào DataFrame
users_df = pd.read_csv("/home/hiro/Downloads/users.csv")

providers_df = pd.read_csv('/home/hiro/Downloads/group_services.csv')

# Thay thế provider_id
providers_df = replace_user_id(users_df.copy(), providers_df.copy())

providers_df.to_csv("/home/hiro/Downloads/users1.csv",index=False)