import requests 
import random
import time

while True:
    random_ticker = random.choice(["TWLO","DDOG"])
    response = requests.get("http://localhost:1313/stock/"+random_ticker)
    time.sleep(1)