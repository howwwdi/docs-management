package tgbot

import (
	corefuncs "docs-managment/internal/core-functions"
	pinata_api "docs-managment/internal/pinata-api"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	contract_interactions "docs-managment/internal/contract-interactions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func Init() error {
	godotenv.Load()

	bot_token := os.Getenv("TG_ACCESS_KEY")

	rpcUrl := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	hexPK := os.Getenv("PRIVATE_KEY")
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	contractABI, err := contract_interactions.LoadABI("ABI.json")
	if err != nil {
		log.Fatalf("Failed to load ABI: %v", err)
	}
	bot, err := tgbotapi.NewBotAPI(bot_token)
	if err != nil {
		return fmt.Errorf("failed to initialize bot: %v", err)
	} else {
		log.Printf("Bot initialized: %s", bot.Self.UserName)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("failed to get updates channel: %v", err)
	}

	waitingForFile := false
	var chatID int64

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				chatID = update.Message.Chat.ID
				commandsMessage := `
				<a href="https://sepolia.arbiscan.io/address/0xa693b16bb3bdbadc6001c653b13d36f65c9849db#code"><b>Смарт-контракт для документооборота на блокчейне</b></a>

<b>Доступные команды:</b>
1. <b>/start</b> - Информация о боте
2. <b>/upload</b> - Загрузить в IPFS хранилище и добавить информацию о файле в блокчейн
3. <b>/download [Doc ID]</b> - Скачать файл из IPFS хранилища (по id документа)
4. <b>/delete [Doc ID]</b> - Удалить файл из IPFS хранилища и стереть данные из контракта (по id документа)
5. <b>/help</b> - Получить список команд`
				msg := tgbotapi.NewMessage(chatID, commandsMessage)
				msg.ParseMode = tgbotapi.ModeHTML
				bot.Send(msg)
				continue
			case "help":
				chatID = update.Message.Chat.ID
				commandsMessage := `<b>Доступные команды:</b>
1. <b>/start</b> - Информация о боте
2. <b>/upload</b> - Загрузить в IPFS хранилище и добавить информацию о файле в блокчейн
3. <b>/download [Doc ID]</b> - Скачать файл из IPFS хранилища (по id)
4. <b>/delete [Doc ID]</b> - Удалить файл из IPFS хранилища и стереть данные из контракта
5. <b>/help</b> - Получить список команд`
				msg := tgbotapi.NewMessage(chatID, commandsMessage)
				msg.ParseMode = tgbotapi.ModeHTML
				bot.Send(msg)
				continue
			case "upload":
				waitingForFile = true
				chatID = update.Message.Chat.ID
				msg := tgbotapi.NewMessage(chatID, "Отправьте файл")
				bot.Send(msg)
				continue
			case "download":
				chatID = update.Message.Chat.ID
				parts := strings.Split(update.Message.Text, " ")
				if len(parts) != 2 {
					msg := tgbotapi.NewMessage(chatID, "Неверный формат команды")
					bot.Send(msg)
					continue
				}
				userID := parts[1]
				var filename string
				var ipfshash string
				var fileBytes []byte

				if docId, err := strconv.ParseUint(userID, 10, 64); err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ID."))
					continue
				} else {
					ipfshash, filename, err = corefuncs.GetDocument(client, common.HexToAddress(contractAddress), contractABI, docId)
					if err != nil {
						log.Fatalf("Failed to get document: %v", err)
					} else {
						fileBytes, _ = pinata_api.DownloadFromPinata(ipfshash, filename)
					}
				}

				file := tgbotapi.FileBytes{
					Name:  filename,
					Bytes: fileBytes,
				}
				msg := tgbotapi.NewMessage(chatID, "Скачиваю файл из IPFS хранилища")
				bot.Send(msg)

				bot.Send(tgbotapi.NewDocumentUpload(chatID, file))
				continue
			case "delete":
				chatID = update.Message.Chat.ID
				parts := strings.Split(update.Message.Text, " ")
				if len(parts) != 2 {
					msg := tgbotapi.NewMessage(chatID, "Неверный формат команды")
					bot.Send(msg)
					continue
				}
				userID := parts[1]
				var ipfshash string


				if docId, err := strconv.ParseUint(userID, 10, 64); err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ID."))
					continue
				} else {
					ipfshash, _, err = corefuncs.GetDocument(client, common.HexToAddress(contractAddress), contractABI, docId)
					if err != nil {
						log.Fatalf("Failed to get document: %v", err)
					} else {
						corefuncs.DeleteFile(client, common.HexToAddress(contractAddress), hexPK, contractABI, docId, ipfshash)
						commandsMessage := `<b>Файл успешно удален</b>`
						msg := tgbotapi.NewMessage(chatID, commandsMessage)
						msg.ParseMode = tgbotapi.ModeHTML
						bot.Send(msg)
					}
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
				bot.Send(msg)
				continue
			}
		}

		if waitingForFile {
			if update.Message.Document != nil {
				fmt.Println(update.Message.Document.FileName)
				fileID := update.Message.Document.FileID
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
				if err != nil {
					log.Println("error while recieving file", err)
					msg := tgbotapi.NewMessage(chatID, "Ошибка получения файла")
					bot.Send(msg)
					continue
				}

				fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot_token, file.FilePath)
				log.Println("Downloading file:", fileURL)

				resp, err := http.DefaultClient.Get(fileURL)
				if err != nil {
					return fmt.Errorf("failed to download file: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					data, _ := io.ReadAll(resp.Body)
					docID, txHash, ipfsHash := corefuncs.UploadFile(client, update.Message.Document.FileName, common.HexToAddress(contractAddress), hexPK, contractABI, data)
					chatID = update.Message.Chat.ID
					if docID == big.NewInt(0) || txHash == "" || ipfsHash == "" {
						msg := tgbotapi.NewMessage(chatID, "Ошибка загрузки файла в блокчейн")
						bot.Send(msg)
					} else {
						commandsMessage := fmt.Sprintf(`<b>Файл успешно загружен</b>
<b>Document ID:</b> %d
<b>Transaction Hash:</b> <a href="https://sepolia.arbiscan.io/tx/%s">%s</a>
<b>IPFS Hash:</b> %s`, docID, txHash, txHash, ipfsHash)
						msg := tgbotapi.NewMessage(chatID, commandsMessage)
						msg.ParseMode = tgbotapi.ModeHTML
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					msg := tgbotapi.NewMessage(chatID, "Ошибка загрузки файла в IPFS")
					bot.Send(msg)
				}
				waitingForFile = false
			} else {
				msg := tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте файл")
				bot.Send(msg)
			}
		}

	}

	return nil
}
