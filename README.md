# 📦 AVWarehouse – IoT-Based Asset Management System

AVWarehouse is a smart warehouse asset management system developed using Golang. It enables seamless tracking and control of goods inventory using CRUD operations, integrating both software and IoT automation capabilities. The backend is fully functional and tailored for warehouse environments requiring robust asset monitoring and management.

---

## 🔧 Features

- 📋 Full CRUD operations for asset data (Create, Read, Update, Delete)
- 🔍 Real-time item tracking
- 📦 Categorization of assets by name, type, and quantity
- 🌐 RESTful API integration
- 🏗️ Scalable architecture built using Golang
- 📁 Asset information stored in a structured format

---

## ⚙️ Tech Stack

- **Backend**: Go (Golang)
- **Framework**: net/http (standard library)
- **Database**: SQLite / GORM (if applicable)
- **Tools**: Postman (for API testing), cURL (optional), Git

---

## 🚀 Getting Started

To run this project locally, follow these steps:

### 1. Clone the repository

```bash
git clone https://github.com/farhanjaa/AVWarehouse.git
cd AVWarehouse

```
2. Run the application
```bash
go run main.go
```
3. Access the API
Once the server is running, you can access the endpoints at:
```bash
http://localhost:8080/
```
📡 API Endpoints
| Method | Endpoint      | Description             |
| ------ | ------------- | ----------------------- |
| GET    | `/assets`     | Retrieve list of assets |
| GET    | `/assets/:id` | Get asset by ID         |
| POST   | `/assets`     | Create a new asset      |
| PUT    | `/assets/:id` | Update asset by ID      |
| DELETE | `/assets/:id` | Delete asset by ID      |


