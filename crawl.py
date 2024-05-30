import requests
from bs4 import BeautifulSoup
import json
import threading

def crawl(postal_code):
    # URL of the webpage
    url = "https://mabuuchinh.vn/default.aspx?page=SearchMBCEN"

    # Cookie string
    cookie = f'MBC=text_search={postal_code}&type_search=0'

    # Headers with the cookie
    headers = {
        "Cookie": cookie
    }

    # Send a GET request to the webpage with the cookie
    response = requests.get(url, headers=headers)

    # Check if the request was successful
    if response.status_code == 200:
        # Parse the HTML content
        soup = BeautifulSoup(response.content, "html.parser")

        # Find the element using XPath
        element = soup.find("span", {"id": "ctl00_ctl06_rptMBC_ctl00_MBC_POSTCODE"})
        
        # Create a dictionary from the provided data
        address_data = {
            "country_code": "VI",
            "zipcode": "AD300",
            "place": "Ordino",
            "state": "Ordino",
            "state_code": "",
            "province": "",
            "province_code": "",
            "community": "",
            "community_code": "",
            "latitude": "",
            "longitude": ""
        }
        
        # Print the element if found
        if element:
            address_data["zipcode"] = element.text
            if len(address_data["zipcode"]) < 4 or "-" in element.text:
                return None
        else:
            return None

        # Find the element using XPath
        element = soup.find("span", {"id": "ctl00_ctl06_rptMBC_ctl00_MBC_Name"})
        
        # Print the element if found
        if element:
            address_data["place"] = element.text
            address_data["state"] = element.text
        else:
            return None

        return address_data
    else:
        return None

def crawl_range(start, end, stride):
    data = []
    for i in range(start, end, 1):
        result = crawl(i)
        if result is not None:
            data.append(result)
            print(result)
    # Open the file in write mode ('w') 
    with open(f'./postal-codes-json-xml-csv/data/address_data_{start}_{end}_{stride}.json', 'w') as outfile:
    # Dump the list of dictionaries to the file as JSON
        json.dump(data, outfile, indent=4)  # Add indent for readability (optional)

        print("Data exported to address_data.json")

# Define the stride and the range
stride = 10000
start = 1
end = 100000

# Create threads for each stride
threads = []
for i in range(start, end, stride):
    t = threading.Thread(target=crawl_range, args=(i, min(i + stride, end), stride))
    threads.append(t)
    t.start()

# Wait for all threads to complete
for t in threads:
    t.join()


