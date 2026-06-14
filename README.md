# Steam Workshop Downloader

[![Go Version](https://img.shields.io/badge/Go-1.23.4-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails Version](https://img.shields.io/badge/Wails-v2.12.0-blue?style=flat)](https://wails.io/)
[![Svelte Version](https://img.shields.io/badge/Svelte-v5-FF3E00?style=flat&logo=svelte)](https://svelte.dev/)
[![Platform Support](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey?style=flat)](https://github.com/)
[![Build](https://github.com/chadlrnsn/workshop-downloader/actions/workflows/deploy.yml/badge.svg)](https://github.com/chadlrnsn/workshop-downloader/actions/workflows/deploy.yml)

[English](#english) | [Русский](#русский)

![img](https://i.imgur.com/CHYqPgO.png)

## English

A professional desktop application built using Wails, Go, and Svelte 5 to simplify downloading items from the Steam Workshop. The application drives the SteamCMD client under the hood, managing asynchronous queues, active downloading states, and persistent Steam Guard logins.

### Features

* Internationalization: Full English and Russian localization.
* Automatic Dependency Resolution: System automatically downloads, unpacks, and configures the latest SteamCMD build if it is not present at the specified path.
* Steam Guard Integration: Interactive prompt handling in the user interface for Email and Mobile Authenticator (2FA) verification during account credentials setup.
* Persistent Login Support: Credentials caching and secure token validation through local Steam Sentry/Sentry (SSFN) tokens.
* Live Log Streaming: Execution logs and download progress updates are piped from SteamCMD directly to the interactive console view.
* Control Panel and Settings: Directory controls, active sessions overview, queue retry facilities, and a secure authentication reset mechanism.

### Prerequisites

* Go version 1.23.4 or higher
* Node.js and npm
* Wails CLI

### Dev Environment Setup

1. Clone the repository to your local directory.
2. Install the frontend dependencies:
   ```bash
   cd frontend
   npm install
   cd ..
   ```
3. Run the application in hot-reload development mode:
   ```bash
   wails dev
   ```

### Building for Production

Compile the compiled binary for your host operating system:
```bash
wails build
```

---

## Русский

Профессиональное настольное приложение, разработанное с использованием Wails, Go и Svelte 5, призванное упростить процесс скачивания модов и карт из Steam Workshop. Программа управляет утилитой SteamCMD под капотом, обеспечивая асинхронные очереди загрузок, отслеживание прогресса и поддержку авторизации с двухфакторной проверкой Steam Guard.

### Возможности

* Интернационализация: Полная поддержка английского и русского языков.
* Автоматическая установка зависимостей: Приложение самостоятельно загрузит, распакует и настроит исполняемый файл SteamCMD при его отсутствии по указанному пути.
* Интеграция с Steam Guard: Перехват интерактивных запросов двухфакторной авторизации (код по электронной почте или мобильный аутентификатор) и отправка ответов прямо из графического интерфейса.
* Сохранение сессии авторизации: Поддержка безопасного кэширования токенов входа (Sentry/SSFN) для последующих скачиваний без ввода пароля.
* Потоковая передача логов: Вывод консольных событий SteamCMD и логов загрузки в реальном времени.
* Панель управления и настроек: Изменение путей сохранения, просмотр статистики сессий, кнопка принудительного сброса авторизации и очистки сохраненных учетных данных.

### Требования к окружению

* Go версии 1.23.4 или выше
* Node.js и npm
* Интерфейс командной строки Wails CLI

### Запуск в режиме разработки

1. Склонируйте репозиторий в локальную папку.
2. Установите зависимости клиентской части:
   ```bash
   cd frontend
   npm install
   cd ..
   ```
3. Запустите приложение в режиме разработки с возможностью "горячей перезагрузки":
   ```bash
   wails dev
   ```

### Сборка готового приложения

Выполните команду сборки исполняемого файла для вашей операционной системы:
```bash
wails build
```
