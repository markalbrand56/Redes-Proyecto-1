# Project #1: XMPP Chat Client

[![Wails build](https://github.com/markalbrand56/Redes-Proyecto-1/actions/workflows/wails.yml/badge.svg)](https://github.com/markalbrand56/Redes-Proyecto-1/actions/workflows/wails.yml)

Universidad del Valle de Guatemala

## Author

**Mark Albrand** (21004)

## Table of Contents

- [Project #1: XMPP Chat Client](#project-1-xmpp-chat-client)
  - [Author](#author)
  - [Table of Contents](#table-of-contents)
  - [Project Description](#project-description)
    - [Tech Stack](#tech-stack)
      - [Installation of dependencies](#installation-of-dependencies)
        - [Go](#go)
        - [Wails](#wails)
        - [Node.js](#nodejs)
        - [Tailwind CSS](#tailwind-css)
  - [Functionalities](#functionalities)
    - [Project requirements](#project-requirements)
      - [Account Management](#account-management)
      - [Communication](#communication)
    - [Additional functionalities](#additional-functionalities)
  - [Project Structure](#project-structure)
  - [How to Run](#how-to-run)
  - [License](#license)

## Project Description

### Tech Stack

- **Frontend**: Vue.js
  - **Styling**: Tailwind CSS
- **Backend**: Go 1.21
- **Linker**: Wails

This project is built using the Wails framework, which allows to build desktop applications using Go and a JavaScript frontend framework. The frontend is built using Vue.js, and the backend is built using Go.

#### Installation of dependencies

To run this project, you need to have the following dependencies installed on your system:

##### Go

To install Go, follow the instructions on the [official website](https://golang.org/doc/install) and download the installer for your operating system. Make sure to get at least **Go 1.21**.

##### Wails

To install Wails, follow the instructions on the [official website](https://wails.io/docs/gettingstarted/installation). In summary, you can install Wails using the following command:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

This command will install the Wails CLI tool, which is used to build and run Wails applications.

##### Node.js

To install Node.js, follow the instructions on the [official website](https://nodejs.org/en/download/). Make sure to get at least **Node.js 14**.

##### Tailwind CSS

Tailwind does not require a global installation, but you can install it as a development dependency on the frontend project. To do so, run the following commands:

```bash
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

## Functionalities

This project aims to build a chat client using the XMPP protocol. Its requirements were defined by the course instructor, and additional functionalities were added to the project.

This implementation uses the XMPP protocol to communicate with a server, and it uses `gosrc.io/xmpp` library to handle the communication and connection. More information about the library can be found [here](https://pkg.go.dev/gosrc.io/xmpp).

### Project requirements

This project had to be built around the server given for this project, which is a custom server built by the course instructors. This means that this project will not be fully functional with other XMPP servers, as it uses custom implementations for some functionalities.

> **Note**: One of this functionalities is to be able to create an account. This was a challenge with the server provided, as it did not allow anonymous connections. This is why for this functionality this project has to connect to a pre-existing account in order to create a new account.

#### Account Management

- [x] Register a new account
- [x] Login with an existing account
- [x] Logout from the account
- [x] Delete the account

#### Communication

- [x] Show the list of contacts, and their status
- [x] Add a new contact
- [x] Show details of a contact
- [x] Send a message to a contact (one-to-one chat)
- [x] Participate in a group chat
- [x] Define a status and a status message
- [x] Send/receive notifications
- [x] Send/receive files

### Additional functionalities

- [x] Get archive of messages
- [x] Delete group chat
- [x] Invite a contact to a existing private group chat
- [x] Delete a contact from the contact list
- [x] Persistent user status and status message

## Project Structure

The project is structured as follows, the project itself is contained in the `cmd` directory. This directory contains the backend application, which is built using Go, and the frontend application, which is built using Vue.js.

```bash
├── cmd/
│   ├── backend/
│   │   ├── models/            # Models for the backend application
│   │   │    └── stanza/           # Structs to represent XMPP stanzas not covered by the library
│   │   └── chat/
│   │       ├── chat.go            # Core chat event handling functions
│   │       ├── handlers.go        # Event handlers for various chat operations
│   │       ├── register.go        # Functions to register events
│   │       └── events/            # Event definitions for the chat application
│   │
│   ├── frontend/
│   │   ├── index.html         # HTML template for the frontend
│   │   ├── wails/             # Wails files for linking the frontend and backend
│   │   ├── dist/              # Distribution folder for the frontend build
│   │   ├── src/               # Source code for the frontend
│   │   │   ├── assets/            # Static assets like images and fonts
│   │   │   ├── components/        # Vue components used across the application
│   │   │   ├── pages/             # Vue components representing pages/views
│   │   │   ├── App.vue            # Root Vue component
│   │   │   ├── main.js            # Entry point for the frontend application
│   │   │   ├── router.js          # Vue Router setup for page navigation
│   │   │   └── style.css          # Global CSS styles
│   │   ├── package.json               # Node.js package file
│   │   ├── tailwind.config.js         # Tailwind CSS configuration file
│   │   └── vite.config.js             # Vite configuration file
│   │
│   ├── wails.json                 # Wails configuration file
│   ├── app.go                     # Entry point for the Go backend application
│   ├── go.mod                     # Go module file
│   ├── go.sum                     # Go dependencies checksum file
│   └── main.go                    # Main Go application logic
│
├── .gitignore                 # Git ignore file
├─  README.md                  # Project README file 
└── LICENSE                    # Project license file


```

## How to Run

This project provides binaries for Windows, MacOS and Linux in the release section on GitHub. For Windows, it also provides a NSIS installer to install the application.

Additionally, you can run the project from source code. To do so, you need to have the dependencies listed on the [tech stack](#tech-stack) section.

Wails handles the building and running of the project, including the frontend and backend. So, to run the project from source code, follow these steps:

1. Clone the repository
2. Navigate to the project directory
3. Enter to the `cmd` directory
4. To run on development mode, run the following command:

```bash
wails dev
```

5. To build the project, run the following command:

```bash
wails build
```

6. The binaries will be available in the `build` directory

> As stated before on the [project functionalities](#functionalities) section, this project was built around a custom server. If you want to use a different server, you will need to modify the `address` variable on the `chat.go` file over on the `cmd/backend/chat` directory. 

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
