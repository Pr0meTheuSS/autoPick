import requests
from bs4 import BeautifulSoup

# 1. EXIST.RU

def parse_exist_ru():
    url = "https://exist.ru/Catalog/TO/"
    response = requests.get(url)
    soup = BeautifulSoup(response.content, 'html.parser')
    
    items = []
    for card in soup.select('.catalog-item-class'):  # Замените на реальный класс карточек
        name = card.select_one('.item-title-class').text.strip()  # Обновите на реальный селектор
        price = card.select_one('.price-class').text.strip()  # Обновите на реальный селектор
        link = card.select_one('a').get('href')
        image = card.select_one('img').get('src')
        items.append({
            'name': name,
            'price': price,
            'link': f"https://exist.ru{link}",
            'image': image
        })
    return items

# 2. AUTODOC.RU

def parse_autodoc_ru():
    url = "https://www.autodoc.ru"
    response = requests.get(url)
    soup = BeautifulSoup(response.content, 'html.parser')
    
    items = []
    for card in soup.select('.product-card-class'):  # Замените на реальный класс карточек
        name = card.select_one('.product-title-class').text.strip()  # Обновите на реальный селектор
        price = card.select_one('.price-class').text.strip()  # Обновите на реальный селектор
        link = card.select_one('a').get('href')
        image = card.select_one('img').get('src')
        items.append({
            'name': name,
            'price': price,
            'link': f"https://www.autodoc.ru{link}",
            'image': image
        })
    return items

# 3. AMRY.RU

def parse_amry_ru():
    url = "https://www.amry.ru"
    response = requests.get(url)
    soup = BeautifulSoup(response.content, 'html.parser')
    
    items = []
    for card in soup.select('.item-card-class'):  # Замените на реальный класс карточек
        name = card.select_one('.title-class').text.strip()  # Обновите на реальный селектор
        price = card.select_one('.price-class').text.strip()  # Обновите на реальный селектор
        link = card.select_one('a').get('href')
        image = card.select_one('img').get('src')
        items.append({
            'name': name,
            'price': price,
            'link': f"https://www.amry.ru{link}",
            'image': image
        })
    return items

# Остальные парсеры аналогичны по структуре, но с различными селекторами для карточек объявлений.

# Важно: Обновите селекторы для каждого сайта в соответствии с их текущей структурой HTML.

# Пример вызова одного из парсеров
if __name__ == "__main__":
    result_exist = parse_exist_ru()
    for item in result_exist:
        print(item)
