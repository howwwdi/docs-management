# Document Management System / Система Управления Документами

This project is a Document Management System built on Ethereum blockchain and IPFS. It allows users to upload, download, and delete documents using a Telegram bot interface.

### Project Structure

```
.env
.gitignore
ABI.json
artifacts/
	build-info/
	contracts/
		DocumentManagement.sol/
			...
cache/
cmd/
	deploy/
		deploy.go
	docs-management/
		docs.go
contracts/
	DocumentManagement.sol
document.txt
go.mod
go.sum
hardhat.config.js
internal/
	contract-interactions/
		interactions.go
	core-functions/
		core-funcs.go
	pinata-api/
		pinata.go
	tg-bot/
		tg-bot.go
package.json
```

### Getting Started

#### Prerequisites

- Go 1.16+
- Node.js
- Hardhat
- Ethereum client (e.g., Infura)
- Pinata account
- Telegram Bot API token

#### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/document-management-system.git
    cd document-management-system
    ```

2. Install Go dependencies:
    ```sh
    go mod tidy
    ```

3. Install Node.js dependencies:
    ```sh
    npm install
    ```

4. Compile the smart contract:
    ```sh
    npx hardhat compile
    ```

#### Configuration

1. Create a 

.env

 file in the root directory and add the following environment variables:
    ```
    RPC_URL=<Your Ethereum RPC URL>
    PRIVATE_KEY=<Your Ethereum Private Key>
    PINATA_API_KEY=<Your Pinata API Key>
    PINATA_SECRET=<Your Pinata Secret Key>
    TG_ACCESS_KEY=<Your Telegram Bot API Token>
    CONTRACT_ADDRESS=<Deployed Contract Address>
    ```

#### Deployment

1. Deploy the smart contract:
    ```sh
    go run cmd/deploy/deploy.go
    ```

#### Running the Telegram Bot

1. Start the Telegram bot:
    ```sh
    go run cmd/docs-management/docs.go
    ```

### Usage

#### Telegram Bot Commands

- `/start` - Information about the bot
- `/upload` - Upload a file to IPFS and add its information to the blockchain
- `/download [Doc ID]` - Download a file from IPFS by document ID
- `/delete [Doc ID]` - Delete a file from IPFS and remove its data from the contract
- `/help` - Get a list of available commands

### License

This project is licensed under the MIT License. See the LICENSE file for details.

### Acknowledgements

- [Ethereum](https://ethereum.org/)
- [IPFS](https://ipfs.io/)
- [Pinata](https://pinata.cloud/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Hardhat](https://hardhat.org/)

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

---

Этот проект представляет собой Систему Управления Документами, построенную на блокчейне Ethereum и IPFS. Он позволяет пользователям загружать, скачивать и удалять документы с использованием интерфейса Telegram-бота.

### Структура Проекта

```
.env
.gitignore
ABI.json
artifacts/
	build-info/
	contracts/
		DocumentManagement.sol/
			...
cache/
cmd/
	deploy/
		deploy.go
	docs-management/
		docs.go
contracts/
	DocumentManagement.sol
document.txt
go.mod
go.sum
hardhat.config.js
internal/
	contract-interactions/
		interactions.go
	core-functions/
		core-funcs.go
	pinata-api/
		pinata.go
	tg-bot/
		tg-bot.go
package.json
```

### Начало Работы

#### Предварительные Требования

- Go 1.16+
- Node.js
- Hardhat
- Ethereum клиент (например, Infura)
- Аккаунт Pinata
- Токен API Telegram Bot

#### Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/document-management-system.git
    cd document-management-system
    ```

2. Установите зависимости Go:
    ```sh
    go mod tidy
    ```

3. Установите зависимости Node.js:
    ```sh
    npm install
    ```

4. Скомпилируйте смарт-контракт:
    ```sh
    npx hardhat compile
    ```

#### Конфигурация

1. Создайте файл 

.env

 в корневом каталоге и добавьте следующие переменные окружения:
    ```
    RPC_URL=<Ваш Ethereum RPC URL>
    PRIVATE_KEY=<Ваш Ethereum Приватный Ключ>
    PINATA_API_KEY=<Ваш Pinata API Ключ>
    PINATA_SECRET=<Ваш Pinata Секретный Ключ>
    TG_ACCESS_KEY=<Ваш Telegram Bot API Токен>
    CONTRACT_ADDRESS=<Адрес Развернутого Контракта>
    ```

#### Развертывание

1. Разверните смарт-контракт:
    ```sh
    go run cmd/deploy/deploy.go
    ```

#### Запуск Telegram Бота

1. Запустите Telegram бота:
    ```sh
    go run cmd/docs-management/docs.go
    ```

### Использование

#### Команды Telegram Бота

- `/start` - Информация о боте
- `/upload` - Загрузить файл в IPFS и добавить его информацию в блокчейн
- `/download [Doc ID]` - Скачать файл из IPFS по идентификатору документа
- `/delete [Doc ID]` - Удалить файл из IPFS и удалить его данные из контракта
- `/help` - Получить список доступных команд

### Лицензия

Этот проект лицензирован по лицензии MIT. См. файл LICENSE для подробностей.

### Благодарности

- [Ethereum](https://ethereum.org/)
- [IPFS](https://ipfs.io/)
- [Pinata](https://pinata.cloud/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Hardhat](https://hardhat.org/)

### Вклад

Вклады приветствуются! Пожалуйста, откройте issue или отправьте pull request для любых улучшений или исправлений ошибок.

---

Feel free to customize this README to better fit your project's needs.
