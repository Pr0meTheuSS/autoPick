from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager
from bs4 import BeautifulSoup
import time
import requests
import redis
import json

# Инициализация Redis клиента
redis_cache = redis.StrictRedis(host='localhost', port=6379, db=0, decode_responses=True)

app = FastAPI()

# Добавление middleware для CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Разрешить все источники
    allow_credentials=True,
    allow_methods=["*"],  # Разрешить все методы (GET, POST и т.д.)
    allow_headers=["*"],  # Разрешить все заголовки
)

# Модель для отображения информации о товаре
class CarPart(BaseModel):
    title: str
    link: str
    article: str
    brand: str
    country: str
    price: str
    stock: str

# Функция для парсинга данных с сайта Drom
def parse_drom(search_string: str, model: str, page: int = 1):
    url = f"https://baza.drom.ru/sell_spare_parts/+/{search_string}/{model}/?page={page}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124"
    }
    response = requests.get(url, headers=headers)
    response.raise_for_status()

    soup = BeautifulSoup(response.text, 'html.parser')
    items = []

    for card in soup.select('.bull-item__content-wrapper'):
        title_element = card.select_one('.bull-item__subject-container a.bulletinLink')
        title = title_element.text.strip() if title_element else 'Неизвестно'
        link = title_element['href'] if title_element else '#'

        image_element = card.select_one('.bull-image-container img')
        image_url = image_element['src'] if image_element else 'https://via.placeholder.com/150'

        brand_element = card.select_one('.bull-item__annotation-row.manufacturer')
        brand = brand_element.text.strip() if brand_element else 'Неизвестно'

        price_element = card.select_one('.price-block__price')
        price = price_element.get('data-price', '0').strip() if price_element else '0'

        location_element = card.select_one('.bull-delivery__city')
        location = location_element.text.strip() if location_element else 'Не указано'

        date_element = card.select_one('.date')
        date = date_element.text.strip() if date_element else 'Неизвестно'

        items.append({
            'title': title,
            'link': f"https://baza.drom.ru{link}",
            'image_url': image_url,
            'brand': brand,
            'price': f"{price} ₽",
            'location': location,
            'date': date
        })

    return items

# Функция для кэширования данных в Redis
def cache_drom_redis(search_string: str, model: str, page: int = 1):
    key = f"drom:{search_string}:{model}:{page}"

    # Проверка наличия в кэше
    cached_data = redis_cache.get(key)
    if cached_data:
        print(f"Returning cached data for key: {key}")
        return json.loads(cached_data)

    # Получение данных с сайта, если они не найдены в кэше
    items = parse_drom(search_string, model, page)

    # Кэширование данных в Redis без времени жизни
    redis_cache.set(key, json.dumps(items))
    print(f"Cached data for key: {key} without expiration")

    return items

# Маршрут для получения данных о запчастях с кэшированием
@app.get("/drom/{search_string}/{model}", response_model=List[dict])
async def get_drom_parts(search_string: str, model: str, page: int = 1):
    try:
        items = cache_drom_redis(search_string, model, page)
        if not items:
            raise HTTPException(status_code=404, detail="Запчасти не найдены")
        return items
    except Exception as e:
        print(str(e))
        raise HTTPException(status_code=500, detail=str(e))

# Маршрут для получения данных о запчастях
@app.get("/parts/{brand}/{model}", response_model=List[CarPart])
async def get_car_parts(brand: str, model: str):
    try:
        items = parse_autotrade_su(brand, model)
        if not items:
            raise HTTPException(status_code=404, detail="Запчасти не найдены")
        return items
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# Маршрут для получения списка марок автомобилей
@app.get("/brands", response_model=List[str])
async def get_car_brands():
    try:
        brands = parse_autotrade_brands()
        if not brands:
            raise HTTPException(status_code=404, detail="Марки не найдены")
        return brands
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
