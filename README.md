# Go-Dev-Chat

Go-Dev-Chat is a Command Line Interface (CLI) Chat Application developed with the Go programming language. It enables users to engage in real-time conversations through a terminal interface, utilizing websockets for instant message updates and leveraging the [tview](https://pkg.go.dev/github.com/rivo/tview#section-readme) library for an intuitive user interface.

## Features

- **Websockets for Real-Time Updates:** Go-Dev-Chat employs websockets to facilitate real-time message updates, ensuring a seamless and dynamic chat experience.

- **Kafka for Message Distribution:** The application utilizes producers and consumers with Kafka, enabling efficient handling and distribution of real-time messages across servers.

- **Postgres for Persistent Storage:** User information and messages are stored persistently in a Postgres database, ensuring data integrity and enabling seamless retrieval of historical messages.

- **RESTful Endpoints with Gin:** Go-Dev-Chat uses the Gin framework to build RESTful endpoints, providing a robust and scalable foundation for handling various functionalities.

- **Websocket Endpoint with Gorilla/Websocket:** The Gorilla/Websocket library is employed for managing the websocket endpoint, ensuring efficient and reliable communication.

- **Gorm for Postgres Integration:** Gorm, a powerful Object-Relational Mapping (ORM) library, is used for seamless integration with Postgres, simplifying database interactions in Go.

## Functionality

Go-Dev-Chat offers a range of functionalities for an immersive chat experience:

- **User Authentication:** Users can sign up by creating a unique username and password, providing a secure and personalized chat environment.

- **Login Credentials:** The application ensures secure user authentication, allowing users to log in with their credentials to access their personalized chat sessions.

- **Initiate Chats:** Users can start a chat with another user, requiring only the knowledge of the recipient's username.

- **Real-Time Messaging:** Enjoy real-time messaging capabilities, allowing users to send and receive messages instantaneously.

## Usage

To get started with Go-Dev-Chat:

1. **Start Backend:** Run `bash scripts/backend-start.sh` in the project root to launch the backend using Docker.

2. **Launch CLI App:** Execute `go run *.go` within the `cli` directory to initiate the CLI Chat App.

Ensure Docker is installed for a seamless experience.

---

## Preview
![image](https://github.com/devAdhiraj/go-dev-chat/assets/75645547/3e793e91-d78a-443d-99fd-89be80eb2e74)

---

For additional information and documentation, refer to the respective libraries and frameworks used:

- [tview](https://pkg.go.dev/github.com/rivo/tview#section-readme)
- [Gin](https://gin-gonic.com/)
- [gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket#section-readme)
- [Gorm](https://gorm.io/)

