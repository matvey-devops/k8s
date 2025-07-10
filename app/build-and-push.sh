#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔨 Начинаем сборку Docker образа...${NC}"

# Имя образа
IMAGE_NAME="matvey5686/k8s-app"
TAG="latest"

# Сборка образа
echo -e "${GREEN}📦 Собираем образ ${IMAGE_NAME}:${TAG}${NC}"
docker build -t ${IMAGE_NAME}:${TAG} ~/study/k8s/app

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Образ успешно собран${NC}"
else
    echo -e "${RED}❌ Ошибка при сборке образа${NC}"
    exit 1
fi

# Добавляем тег с версией
VERSION=$(date +%Y%m%d%H%M%S)
docker tag ${IMAGE_NAME}:${TAG} ${IMAGE_NAME}:${VERSION}

# Пуш в Docker Hub
echo -e "${GREEN}📤 Отправляем образ в Docker Hub...${NC}"
echo -e "${GREEN}   Убедитесь, что вы вошли в Docker Hub: docker login${NC}"

docker push ${IMAGE_NAME}:${TAG}
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Образ ${IMAGE_NAME}:${TAG} успешно отправлен${NC}"
else
    echo -e "${RED}❌ Ошибка при отправке образа${NC}"
    exit 1
fi

docker push ${IMAGE_NAME}:${VERSION}
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Образ ${IMAGE_NAME}:${VERSION} успешно отправлен${NC}"
else
    echo -e "${RED}❌ Ошибка при отправке образа с версией${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Все готово!${NC}"
echo -e "${GREEN}   Образ доступен как:${NC}"
echo -e "${GREEN}   - ${IMAGE_NAME}:${TAG}${NC}"
echo -e "${GREEN}   - ${IMAGE_NAME}:${VERSION}${NC}" 