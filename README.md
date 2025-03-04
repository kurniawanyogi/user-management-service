# User Management Service

### Description
User Management Service is a RESTful API built using the Go to manage user resources.

### Features
- **Authentication**: Bearer Token
- **Database**: MySQL for storing user data
- **CRUD Operations** on the user resource


### Prerequisites

Before you begin, ensure that you have the following installed on your machine:

- [Go (Golang)](https://golang.org/dl/) version go 1.22.2 or later
- [MySQL](https://dev.mysql.com/downloads/) latest version
- [Postman](https://www.postman.com/) or a similar API testing tool

---

### Installation

1. **Clone the Repository**

   Clone the repository to your local machine:

   ```bash
   git clone https://github.com/kurniawanyogi/user-management-service.git
    ```

2. **Install Dependencies**

   After cloning the project, install the required dependencies by running:

   ```bash
   go mod tidy
    ```

3. **Configure MySQL Database**

    - Create a new MySQL database named user_management

    ```bash
   CREATE DATABASE user_management;
    ```

   - Run the migration to create the necessary tables in MySQL:
     - install goose 
     ```bash
        go install github.com/pressly/goose/v3/cmd/goose@latest` 
       ```
     - `cd migrations`
     - `goose mysql "user:password@/databaseName?parseTime=true" up` update user and password with your mysql credentials
   
4. **Configure Environment Variable**
   
   Update [env config](config.cold.json) based on your local machine

5. **Run Application**

   Once Everything is configured, run go application with `go run main.go` on your main directory

---

### Test

For run the unit test 
   ```bash
      go test ./service/... -v -cover
   ```

---

### Documentation

For API testing with Postman, you can import the [Postman Collection](user-management.postman_collection.json) provided in this repository.