# Go Application

Ini adalah aplikasi Go untuk mengelola artikel. Aplikasi ini menyediakan titik akhir untuk membuat, membaca, memperbarui, dan menghapus artikel.

## Table of Contents

-   [Installation](#installation)
-   [Configuration](#configuration)
-   [Usage](#usage)

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/your-repo-name.git
    cd your-repo-name
    ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

## Configuration

1.  **Database Setup:**

    ```sql
    CREATE DATABASE article;
    ```

2.  **Tabel Setup:**

    ```sql
    CREATE TABLE posts (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Title VARCHAR(200),
    Content TEXT,
    Category VARCHAR(100),
    Created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    Status VARCHAR(100) CHECK (Status IN ('Publish', 'Draft', 'Thrash'))
    );
    ```

3.  **Insert Data:**

    Insert Data dengan menjakan query dari file 
    
    ```
    posts_202405231502.sql
    ```

## Usage

Run aplikasi:
    
    go run main.go


