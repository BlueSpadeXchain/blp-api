# blp-api

# API Repository

This repository contains the backend API for managing user accounts and orders. It is built with Go and integrates with Supabase for database management. The API handles order creation, user management, and data validation.

## Table of Contents

- [blp-api](#blp-api)
- [API Repository](#api-repository)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
  - [API Endpoints](#api-endpoints)
  - [Database Schema](#database-schema)
    - [Users Table](#users-table)
    - [Orders Table](#orders-table)
  - [Development](#development)
  - [Deployment](#deployment)

## Getting Started

To get started with the API, follow these steps:

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/yourusername/your-repo-name.git
    cd your-repo-name
    ```

2. **Set Up Environment Variables:**

    Create a `.env` file in the root directory with the following variables:

    ```env
    SUPABASE_URL=<your_supabase_url>
    SUPABASE_ANON_KEY=<your_supabase_anon_key>
    ```

3. **Install Dependencies:**

    ```bash
    go mod tidy
    ```

4. **Run the Application:**

    ```bash
    go run main.go
    ```

    The server will start on `http://localhost:8081`.

## API Endpoints

- **`GET /api/info`**: Retrieves information about the API.

- **`POST /api/order`**: Creates a new order. Requires a valid JSON payload with the following structure:

    ```json
    {
      "signer": "string",
      "createdOn": "string",
      "chainId": "string",
      "order": {
        "orderId": "string",
        "netValue": "string",
        "amount": "string",
        "collateral": "string",
        "markPrice": "string",
        "entryPrice": "string",
        "liquidationPrice": "string",
        "takeProfit": "string",
        "stopLoss": "string",
        "type": "string" // "long" or "short"
      },
      "messageId": "string",
      "signature": "string"
    }
    ```

## Database Schema

### Users Table

- **id**: Serial primary key
- **signer**: User's Ethereum address (VARCHAR)
- **date_created**: Timestamp of account creation
- **balances**: JSONB object containing currency balances

### Orders Table

- **key**: Serial primary key
- **signer**: User's Ethereum address (VARCHAR)
- **createdOn**: Timestamp of order creation
- **chainId**: Chain identifier (VARCHAR)
- **orderId**: Unique order identifier (VARCHAR)
- **pair**: Trading pair (VARCHAR)
- **netValue**: Net value of the order (VARCHAR)
- **amount**: Amount of the order (VARCHAR)
- **collateral**: Collateral for the order (VARCHAR)
- **markPrice**: Market price at order time (VARCHAR)
- **entryPrice**: Entry price of the order (VARCHAR)
- **liquidationPrice**: Liquidation price (VARCHAR)
- **takeProfit**: Take profit level (VARCHAR)
- **stopLoss**: Stop loss level (VARCHAR)
- **type**: Order type ("long" or "short")

## Development

- **To add a new currency:** Use SQL commands to update the `balances` JSONB field to include new currencies.

- **To handle database migrations:** Use Supabase tools or manual SQL scripts to adjust schema as needed.

## Deployment

- **Deploy to Vercel:** Currently, Vercel does not support WebSocket, so you may need to use a separate service like [Heroku](https://www.heroku.com/) or [Railway](https://railway.app/) for WebSocket support.

- **Configure Supabase:** Ensure that your Supabase instance is correctly set up with the `users` and `orders` tables.

