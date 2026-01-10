ERP6-BE-GOLANG

ERP Backend (Capella Engine v6) written in Golang

🚀 Overview

ERP6-BE-GOLANG is the backend engine for Capella Engine v6 — a modular and scalable Enterprise Resource Planning (ERP) backend built with Go (Golang).
This project aims to provide a high-performance, extensible ERP backend that supports business processes, modular plugins, and clean architecture principles. It serves as the core API engine for ERP6, powering data management, workflows, and integrations.

📌 Key Features

✔️ Modular plugin system (extendable modules)
✔ RESTful API powered by Go
✔ Clean & scalable architecture
✔ Database agnostic (supports PostgreSQL / MySQL / SQLite)
✔ Built-in configuration & environment setup
✔ Swagger / OpenAPI documentation support

Add or update this list based on actual implemented features.

🧠 Architecture

The project follows a modular layered structure with separation of concerns. Core concepts include:

Core Modules — business logic & domain

Plugins — extend functionality without modifying core

Models — data models & ORM definitions

Response Helpers — standardized API responses

The design enables flexible expansion and service integration while keeping code testable and maintainable.

📦 Requirements

Before running the project, ensure you have:

Go 1.22+ installed

Supported database (e.g., PostgreSQL / MySQL / SQLite / SQL Server / Oracle)

Environment variables configured (see .env.example)

🔧 Installation

Clone the repository

git clone https://github.com/romyandre79/erp6-be-golang.git
cd erp6-be-golang


Set up environment

Copy the sample environment file:

cp .env.example .env


Modify .env values for your database and configuration.

Build and Run

go build -o erp6-be
./erp6-be


Or run directly:

go run .

🛠 API Documentation

Swagger / OpenAPI documentation is available at:

/swagger/index.html


(Include correct link path if hosted locally or deployed)

🛡 License

This project is licensed under GPL-3.0 License. See the LICENSE
 file for details.

🗣 Contributing

Contributions are welcome!
To contribute:

Fork the repository

Create your feature branch

Commit your changes

Open a Pull Request

Please ensure your code follows project standards and is tested.

📞 Support

If you have questions or need help:

Open an issue

Reach out via GitHub discussions

✨ Thank you for using / contributing to ERP6 Backend with Golang!