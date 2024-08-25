# Project #1: XMPP Chat Client

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
  - [Functionalities](#functionalities)
    - [Project requirements](#project-requirements)
      - [Account Management](#account-management)
      - [Communication](#communication)
    - [Additional functionalities](#additional-functionalities)
  - [Project Structure](#project-structure)
  - [How to Run](#how-to-run)

## Project Description

### Tech Stack

- **Frontend**: Vue.js
- **Backend**: Go 1.21
- **Linker**: Wails

This project is built using the Wails framework, which allows to build desktop applications using Go and a JavaScript frontend framework. The frontend is built using Vue.js, and the backend is built using Go.

#### Installation of dependencies

#### Go

To install Go, follow the instructions on the [official website](https://golang.org/doc/install) and download the installer for your operating system. Make sure to get at least **Go 1.21**.

#### Wails

To install Wails, follow the instructions on the [official website](https://wails.io/docs/gettingstarted/installation). In summary, you can install Wails using the following command:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

This command will install the Wails CLI tool, which is used to build and run Wails applications.

## Functionalities

This project aims to build a chat client using the XMPP protocol. Its requirements were defined by the course instructor, and additional functionalities were added to the project.

This implementation uses the XMPP protocol to communicate with a server, and it uses `gosrc.io/xmpp` library to handle the communication and connection. More information about the library can be found [here](https://pkg.go.dev/gosrc.io/xmpp).

This project had to be built around the server given for this project, which is a custom server built by the course instructors. This means that this project will not be fully functional with other XMPP servers, as it uses custom implementations for some functionalities.

> **Note**: One of this functionalities is to be able to create an account. This was a challenge with the server provided, as it did not allow anonymous connections. This is why for this functionality this project has to connect to a pre-existing account in order to create a new account.

### Project requirements

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

## Project Structure

## How to Run

This project provides binaries for Windows, MacOS and Linux in the release section on GitHub.

Additionally, you can run the project from source code. To do so, you need to have the dependencies listed on the [tech stack](#tech-stack) section.

To run the project from source code, follow these steps:

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